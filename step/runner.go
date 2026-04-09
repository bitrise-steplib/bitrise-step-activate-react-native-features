package step

// Runner defines the standard Bitrise step lifecycle phases.
type Runner interface {
	ProcessConfig() error
	InstallDeps() error
	Run() error
	ExportOutputs() error
}
