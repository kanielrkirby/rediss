module github.com/piratey7007/rediss/server

require (
	github.com/piratey7007/rediss/lib v0.1.0
	github.com/spf13/cobra v1.7.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

go 1.21.1

replace github.com/piratey7007/rediss/lib => ../lib
