package component

type BuildOptions struct {
	Repository string
	Hash       string
}

type Component interface {
	Build(opts BuildOptions) (BuiltComponent, error)
}

type BuiltComponent interface {
	Run(opts interface{}) error
}
