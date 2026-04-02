package features

//go:generate mockery
type Logger interface {
	Debugf(format string, v ...any)
}
