# shockrays

## What?

A simple command line based application that is used to interact with [ProjectorRays](https://github.com/ProjectorRays/ProjectorRays). `ProjectorRays`
is a decompiler for Shockwave/Director related files e.g. DCR, DXR, CCT, CXT. It's capable of un-protecting these published
files with the Lingo bytecode intact, so that the Lingo source code can be edited again in Director, as well as dumping 
the reconstructed Lingo source code into separate files for each script member.

## Why?

`ProjectorRays` can handle decompiling a directory of Director related files,
but it outputs the Lingo scripts for all the decompiled files in the same directory. This becomes incredibly messy when 
you are decompiling multiple files at once. The primary goal of `shockrays` is to provide a way to programmatically dump the Lingo scripts into 
directories that have subdirectories named after the decompiled files which they came from.

## How?

### Initial Configuration

To get started with `shockrays` you will need to initialize the config file.

```
shockrays confit init
```

Now you should find a directory named `.shockrays` in your home directory. 

If you are on Windows this will be something like `C:\Users\Anon\`, if you are on Linux it will be something like `~/Anon`.
Inside the `.shockrays` directory you should find a directory named `projector-rays`, download or compile the latest 
version of ProjectorRays [here](https://github.com/ProjectorRays/ProjectorRays/releases) and place the binary inside of
this directory. 

By default, the applications expects the binary to be named `projectorrays`, but this can be changed via the config file.
You can also change the directory where the application looks for the binary via your config file, the `--projector-rays`
flag, or by creating an environment variable named `PROJECTOR_RAYS` with the full path to the binary.

### Decompiling Files

To decompile a directory full of Director files located at `C:\Users\Anon\Files` you could run the following command.
The path argument is optional and if it's not provided `shockrays` will look for the files in the `.shockrays` directory.

```
shockrays decompile C:\Users\Anon\Files 
```

This will create a subdirectory in the `ouptput` directory for each target file and decompile every Director related 
file in `C:\Users\Anon\Files`. The file's Lingo scripts are dumped into the newly created subdirectory to keep scripts
neatly organized.

### Help

For more information on the various commands and flags available run `shockrays` with the `--help` flag.

```
shockrays --help
```

Additionally, the `--debug` flag can be used with all the commands to have additional information logged
to the console which may be helpful when debugging issues.

## Where?

The `os` and `filepath` packages from Go's standard library are used to handle path names and
traversing the filesystem in an operating system agnostic manner, so theoretically `shockrays` _should_ work on any system that
`ProjectorRays` will run on.

## Credit Where Credit's Due

To all the wonderful people who made [ProjectorRays](https://github.com/ProjectorRays/ProjectorRays) possible, including but not limited to:
- [Debby Servilla](https://github.com/djsrv)
- [Anthony Kleine](https://github.com/tomysshadow)
- [Earthquake Project team](https://github.com/Earthquake-Project)
- [Just Solve the File Format Problem wiki](http://fileformats.archiveteam.org/wiki/Lingo_bytecode)
- [ScummVM Director engine team](https://www.scummvm.org/credits/#:~:text=Director:)