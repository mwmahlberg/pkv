//go:generate go-bindata -o cmd/bindata.go -ignore .*_test.go -pkg cmd ./verify


package main

import "github.com/mwmahlberg/pkv/cmd"

/*
pkv is a command line tool to manage and use Partial Key Verification.
*/ 
func main() {
	cmd.Execute()
}