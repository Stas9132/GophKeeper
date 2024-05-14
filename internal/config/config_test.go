package config

import (
	"os"
	"testing"
)

func TestConfigInit(t *testing.T) {
	if err := os.Setenv("LISTEN_ADDRESS", "check"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("LISTEN_ADDRESS_REST", "check"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("S3_ENDPOINT", "check"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("S3_ACCESS_KEY", "check"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("S3_SECRET_KEY", "check"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("S3_USE_SSL", "NO"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("S3_LOCATION", "check"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("AES_KEY", "check"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("AES_NOONCE", "0123456789abcdef01234567"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("DATABASE_DSN", "check"); err != nil {
		t.Fatal(err)
	}

	Init()

	if ListenAddress != "check" {
		t.Fatal("ListenAddress", ListenAddress)
	}
	if ListenAddressR != "check" {
		t.Fatal("ListenAddressR", ListenAddressR)
	}
	if S3Endpoint != "check" {
		t.Fatal("S3Endpoint", S3Endpoint)
	}
	if S3AccessKeyID != "check" {
		t.Fatal("S3AccessKeyID", S3AccessKeyID)
	}
	if S3SecretAccessKey != "check" {
		t.Fatal("S3SecretAccessKey", S3SecretAccessKey)
	}
	if S3UseSSL {
		t.Fatal("S3UseSSL", S3UseSSL)
	}
	if string(AESKey) != "check" {
		t.Fatal("AESKey", AESKey)
	}
}
