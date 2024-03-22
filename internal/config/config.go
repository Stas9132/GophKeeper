package config

import (
	"encoding/hex"
	"flag"
)

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

var (
	S3Endpoint        = "127.0.0.1:9000"
	S3AccessKeyID     = "aHLytUVhTKOPMYD6nYA2"
	S3SecretAccessKey = "F2Avh18pul7X8IsGhCTeWPnaQNhlOuda3iAYSO30"
	S3UseSSL          = false
	S3Location        = "us-east-1"
	AESKey            = []byte("AES256Key-32Characters1234567890")
	AESnonce, _       = hex.DecodeString("bb8ef84243d2ee95a41c6c57")
)
