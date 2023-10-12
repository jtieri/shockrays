package shockrays

// Valid Director/Shockwave file extensions.
const (
	cct = "cct"
	cxt = "cxt"
	dcr = "dcr"
	dxr = "dxr"
	dir = "dir"
	cst = "cst"
)

// ValidFileExtension will return true if the extension passed in is a valid Director/Shockwave file extension.
func ValidFileExtension(ext string) bool {
	switch ext {
	case cct, cxt, dcr, dxr, dir, cst:
		return true
	default:
		return false
	}
}
