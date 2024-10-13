//go:build testing

package config

// init function will be built (and this called) only if tag testing (go test --tags=testing ./... in Makefile)
// this is needed because JR will call Init* directly from the CLI because of init ordering
func init() {
	// InitEnvironmentVariables()
	//
	// InitEmitters()
}
