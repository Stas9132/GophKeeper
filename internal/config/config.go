package config

import "flag"

var (
	Version   = "1.0.0"
	BuildDate string
)

var (
	ListenAddress  string = "localhost:2345"
	ListenAddressR string = "localhost:2346"
	PrintVersion   bool
)

func Init() {
	b := flag.Bool("version", false, "print version and build date")
	flag.Parse()

	PrintVersion = *b
}
