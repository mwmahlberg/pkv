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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	pkv "gopkg.in/mwmahlberg/pkv.v1/generate"
)

var lseed *uint32

// genkeyCmd respresents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "generates a product key",
	Long: `Generates a product key from the specified matrix and seed
Note that the seed must be unique per customer, so that you are able to blacklist individual keys.`,
	Run: genKey}

func init() {
	RootCmd.AddCommand(genkeyCmd)

	// Here you will define your flags and configuration settings

	// Cobra supports Persistent Flags which will work for this command and all subcommands
	// genkeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly
	// genkeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle" )

	// BUG(mwmahlberg): We have to convert the seed and it's usages to uint32, not just the input
	lseed = genkeyCmd.Flags().Uint32P("seed", "s", 0, "seed for key. Must be > 0")

	genkeyCmd.Flags().StringVarP(&keyfile, "file", "f", defaultKeyfilePath, "path to keyfile")
	//	t := *lseed
}

func genKey(cmd *cobra.Command, args []string) {
	seed = uint64(*lseed)

	if seed < 1 {
		fmt.Println("Seed can not be smaller than 1!")
		cmd.Usage()
		os.Exit(0)
	}

	k := readKeyFile(keyfile)

	pk := pkv.KeyMatrix{Matrix: k.Matrix}

	s := pk.GetKey(uint64(*lseed))
	if err := pkv.CheckCompleteKey(s, k.Matrix); err != nil {
		panic(err)
	}
	fmt.Println(s)
}
