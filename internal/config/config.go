package config

import (
	"encoding/hex"
	"flag"
	"log"
	"os"
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

	if s, ok := os.LookupEnv("LISTEN_ADDRESS"); ok {
		ListenAddress = s
	}
	if s, ok := os.LookupEnv("LISTEN_ADDRESS_REST"); ok {
		ListenAddressR = s
	}
	if s, ok := os.LookupEnv("S3_ENDPOINT"); ok {
		S3Endpoint = s
	}
	if s, ok := os.LookupEnv("S3_ACCESS_KEY"); ok {
		S3AccessKeyID = s
	}
	if s, ok := os.LookupEnv("S3_SECRET_KEY"); ok {
		S3SecretAccessKey = s
	}
	if s, ok := os.LookupEnv("S3_USE_SSL"); ok {
		S3UseSSL = s == "YES"
	}
	if s, ok := os.LookupEnv("S3_LOCATION"); ok {
		S3Location = s
	}
	if s, ok := os.LookupEnv("AES_KEY"); ok {
		AESKey = []byte(s)
	}
	if s, ok := os.LookupEnv("AES_NOONCE"); ok {
		var err error
		AESnonce, err = hex.DecodeString(s)
		if err != nil {
			log.Fatal(err)
		}
	}
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
