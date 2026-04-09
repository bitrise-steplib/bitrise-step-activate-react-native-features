package step

//go:generate mockery
type Command interface {
	SetArgs(a []string)
	Execute() error
}
