package buildinfo

var version string

// Version returns the VCS version recorded at build time.
func Version() string {
	return version
}
