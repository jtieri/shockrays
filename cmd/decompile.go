package cmd

import (
	"fmt"
	"github.com/jtieri/shockrays/shockrays"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const (
	// PROJECTOR_RAYS is the name of an environment variable that can be set to specify the directory where the ProjectorRays binary should be located.
	PROJECTOR_RAYS = "PROJECTOR_RAYS"
)

func decompileCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decompile",
		Aliases: []string{"dcmp", "dec"},
		Short:   "Decompiles a directory of Shockwave/Director files",
		Long:    "Decompiles Shockwave/Director files within an input directory and dumps the Lingo scripts into a target directory",
		Args:    withUsage(cobra.RangeArgs(0, 1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s decompile
$ %s dcmp C:\Users\Anon\Files
$ %s dec C:\Users\Anon\Files --projector-rays C:\Path\To\ProjectorRays\Bin`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}

			prDir, err := projectorRaysDir(cmd, a)

			// Check that the directory where the ProjectorRays binary is located actually exists.
			if _, err = os.Stat(prDir); err != nil {
				return err
			}

			// Check that the ProjectorRays binary actually exists.
			projectorRays := filepath.Join(prDir, a.config.PRBinName)

			file, err := os.Stat(projectorRays)
			if err != nil {
				return err
			}

			// Check that we actually have a reference to a file now and not a directory.
			if file.IsDir() {
				return fmt.Errorf("expected reference to ProjectorRays binary, got a directory: %s", projectorRays)
			}

			// If an arg is passed use that for the input path, otherwise use the specified home path.
			var inputPath string
			if len(args) == 1 {
				inputPath = args[0]
			} else {
				inputPath = home
			}

			// Get the names of all valid Shockwave/Director files in the input directory.
			validFiles, err := validFileNames(a, inputPath)
			if err != nil {
				return err
			}

			// Run ProjectorRays to decompile all of the files and dump the Lingo scripts.
			if err = decompileFiles(a, inputPath, projectorRays, validFiles); err != nil {
				return err
			}

			return nil
		},
	}

	return projectorRaysFlag(a.viper, cmd)
}

// projectorRaysDir returns the directory where the ProjectorRays binary should be located.
// If determines the directory by checking the config file, for an env variable, & a flag passed at the CLI.
// The order of precedence is as follows:
// 1. the value passed in with the flag
// 2. the value set in the environment variable
// 3. the value from the config file
func projectorRaysDir(cmd *cobra.Command, a *appState) (string, error) {
	prDir := a.config.PRDirectory

	tmpPrDir := os.Getenv(PROJECTOR_RAYS)
	if tmpPrDir != "" {
		prDir = tmpPrDir
	}

	tmpPrDir, err := cmd.Flags().GetString(flagPRDir)
	if err != nil {
		return "", err
	}

	if tmpPrDir != defaultPRDir {
		prDir = tmpPrDir
	}

	return prDir, nil
}

// validFileNames returns a slice containing the names of all the valid Shockwave/Director files found
// in the specified input directory.
func validFileNames(a *appState, inputPath string) ([]string, error) {
	dir, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read all files from the input directory.
	fileNames, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	if len(fileNames) == 0 {
		return nil, fmt.Errorf("the specified input directory contains no files to decompile: %s", inputPath)
	}

	// Filter out all the invalid files.
	var validFiles []string
	for _, name := range fileNames {
		if !strings.Contains(name, ".") {
			continue
		}

		// Split file name at . and check for valid file extension.
		// e.g. cct, cxt, dcr, dxr, cst, dir
		parts := strings.Split(name, ".")
		extension := parts[1]

		if shockrays.ValidFileExtension(extension) {
			a.log.Debug("Found valid Director file", "file_name", name)

			validFiles = append(validFiles, name)
		}
	}

	return validFiles, nil
}

// decompileFiles will attempt to run ProjectorRays to decompile each of the specified files in a concurrent manner.
// It will run the decompilation command from ProjectorRays for each file and then report any errors at the end.
func decompileFiles(a *appState, inputPath, projectorRays string, fileNames []string) error {
	var wg sync.WaitGroup
	errChan := make(chan ErrorReport, len(fileNames))

	for _, name := range fileNames {
		name := name

		wg.Add(1)

		go func(a *appState, inputPath, projectorRays, file string) {
			defer wg.Done()

			if errRep := decompileFile(a, inputPath, projectorRays, file); errRep.Err != nil {
				errChan <- errRep
			}
		}(a, inputPath, projectorRays, name)

	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []ErrorReport
	for err := range errChan {
		errors = append(errors, err)
	}

	a.log.Info(fmt.Sprintf("%d errors were encountered while decompiling %d files", len(errors), len(fileNames)))

	if len(errors) > 0 {
		for _, errRep := range errors {
			a.log.Error("Error returned when attempting to decompile file",
				"file_name", errRep.FileName,
				"error", errRep.Err,
				"cmd", errRep.Cmd,
				"output", errRep.Output,
			)
		}
	}

	return nil
}

// decompileFile creates a new subdirectory in the output directory which is named after the specified file.
// It then runs the decompilation command from ProjectorRays and dumps the Lingo scripts in the newly created directory.
func decompileFile(a *appState, inputPath, projectorRays, file string) ErrorReport {
	parts := strings.Split(file, ".")
	fileName := parts[0]

	outputDir := filepath.Join(a.config.ScriptOutputDir, fileName)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, os.ModePerm); err != nil {
			return ErrorReport{
				FileName: file,
				Err:      err,
			}
		}

		a.log.Debug("Created output directory", "dir_name", outputDir)
	}

	prArgs := []string{
		"decompile", filepath.Join(inputPath, file), "--output", outputDir, "--dump-scripts",
	}

	cmd := exec.Command(projectorRays, prArgs...)
	cmd.Dir = outputDir

	output, err := cmd.Output()
	if err != nil {
		e := ErrorReport{
			FileName: file,
			Err:      err,
			Cmd:      fmt.Sprintf("%s %s", projectorRays, prArgs),
			Output:   string(output),
		}

		return e
	}

	a.log.Info("Decompiled file", "output", output)

	return ErrorReport{}
}

// ErrorReport is used to return information when errors are encountered in the decompilation process.
type ErrorReport struct {
	FileName string `yaml:"file-name" json:"file-name"`
	Err      error  `yaml:"error" json:"error"`
	Cmd      string `yaml:"cmd" json:"cmd"`
	Output   string `yaml:"output" json:"output"`
}
