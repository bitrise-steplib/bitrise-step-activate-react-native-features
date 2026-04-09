package step

//go:generate mockery
type Logger interface {
	EnableDebugLog(enable bool)
	Println()
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	Debugf(format string, v ...any)
	FormattedErrorf(err error)
}
