package step

//go:generate mockery
type InputParser interface {
	Parse(conf any) error
}
