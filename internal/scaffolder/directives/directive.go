package directives

type Directive interface {
	Execute(d string, n string) error
}
