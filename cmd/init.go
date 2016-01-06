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
package cmd

import (
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	generate "gopkg.in/mwmahlberg/pkv.v1/internal"
)

// initCmd respresents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes a matrix used for generating product keys",
	Long: `In order to generate product keys, a secret matrix needs to be generated.
This command creates this matrix by using securely generated random numbers. 
	`,
	Run: initialize,
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&keyfile, "file", "f", defaultKeyfilePath, "path to your matrix file")

}

func initialize(cmd *cobra.Command, args []string) {
	var (
		file *os.File
		err  error
	)

	if _, err = os.Stat(keyfile); err == nil {
		fmt.Printf("File '%s' already exists!\n", keyfile)
		fmt.Println("Cowardly refusing to overwrite it…")
		os.Exit(1)
	}

	if file, err = os.OpenFile(keyfile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600); err != nil {
		fmt.Printf("Error while creating '%s': %s\n", keyfile, err)
		os.Exit(1)
	}
	defer file.Close()

	k := generate.NewKey(rand.Reader)
	enc := gob.NewEncoder(file)

	if err = enc.Encode(k); err != nil {
		fmt.Printf("Error writing key to '%s': %s\n", keyfile, err)
		os.Exit(1)

	}

	//	j, err := json.MarshalIndent(k, " ", "  ")

	//	if err != nil {
	//		fmt.Printf("Error while generating JSON for key")
	//		os.Exit(1)
	//	}
	//	ioutil.WriteFile("pkvkey.json", j, 0600)
}
