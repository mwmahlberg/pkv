//go:generate go-bindata -o cmd/bindata.go -ignore .*_test.go -pkg cmd ./verify
 
package main

import "github.com/mwmahlberg/pkv/cmd"

func main() {
	cmd.Execute()
}