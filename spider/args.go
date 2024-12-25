package spider

type Args struct {
}

type Arg struct {
	Name     string
	Help     string
	HelpArgs string
	Default  interface{}

	isList   bool
	optional bool
}
