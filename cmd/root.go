// Copyright ©2016 Markus W Mahlberg <markus@mahlberg.io>
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/*
Package cmd contains the sources of the "pkv" utility and should not be imported.
*/
package cmd

import (
	"encoding/gob"
	"fmt"
	"os"

	pkv "gopkg.in/mwmahlberg/pkv.v1/internal"

	"github.com/spf13/cobra"
)

const (
	defaultKeyfilePath = "./pkv.key"
)

var (
	cfgFile string
	keyfile string
	k       int
	seed    uint64
	stdout  bool
)

var ()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pkv",
	Short: "product key generation",
	Long: `
pkv generates product keys and the code you need to verify them with the go programming language.

For product code generation and verification, the Partial Key Verification scheme is used.

To start, first create a secret matrix file with the "init" command. This matrix is the basis for all keys
you generate. Never publish it, and keep backups of it.
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func readKeyFile(path string) (key *pkv.KeyMatrix) {

	var (
		file *os.File
		dec  *gob.Decoder
		err  error
	)

	if file, err = os.OpenFile(path, os.O_RDONLY, 0); err != nil {
		fmt.Printf("Error opening key file '%s': %s", path, err)
		os.Exit(1)
	}

	dec = gob.NewDecoder(file)
	key = &pkv.KeyMatrix{}

	if err = dec.Decode(&key); err != nil {
		fmt.Printf("Error reading key file '%s': %s", path, err)
		os.Exit(1)
	}

	return
}
