module github.com/piratey7007/rediss/client

go 1.21.1

require github.com/piratey7007/rediss/lib v0.1.0

require github.com/spf13/cobra v1.7.0

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/piratey7007/rediss/lib => ../lib
