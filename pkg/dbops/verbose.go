//nolint:gochecknoglobals // it will investigate later
package dbops

var verbose bool

func SetVerbose(isVerbose bool) {
	verbose = isVerbose
}
