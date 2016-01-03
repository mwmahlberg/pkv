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

	pkv "github.com/mwmahlberg/pkv/generate"
	"github.com/spf13/cobra"
)

var (
	seed uint64
)

// genkeyCmd respresents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "generates a product key",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a Cli library for Go that empowers applications. This
application is a tool to generate the needed files to quickly create a Cobra
application.`,
	Run: genKey}

func init() {
	RootCmd.AddCommand(genkeyCmd)

	// Here you will define your flags and configuration settings

	// Cobra supports Persistent Flags which will work for this command and all subcommands
	// genkeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly
	// genkeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle" )

	genkeyCmd.Flags().Uint64VarP(&seed, "seed", "s", 0, "seed for key. Must be > 0")
}

func genKey(cmd *cobra.Command, args []string) {

	if seed < 1 {
		cmd.Usage()
		os.Exit(0)
	}

	k := readKeyFile()

	pk := pkv.PartialKey{Matrix: k.Matrix}
	
	s := pk.GetKey(seed)
	if err := pkv.CheckCompleteKey(s,k.Matrix); err != nil {
		panic(err)
	}
	fmt.Println(s)
}
