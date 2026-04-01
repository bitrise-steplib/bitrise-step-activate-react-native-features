package step

type Command interface {
	SetArgs(a []string)
	Execute() error
}
