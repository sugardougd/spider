module spider_v1

go 1.20

require (
	github.com/desertbit/readline v1.5.1
	golang.org/x/crypto v0.26.0
	spider/grumble v1.0.0
)

require (
	github.com/desertbit/closer/v3 v3.1.3 // indirect
	github.com/desertbit/columnize v2.1.0+incompatible // indirect
	github.com/desertbit/go-shlex v0.1.1 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.23.0 // indirect
)

replace spider/grumble => ../grumble
