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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mwmahlberg/pkv/generate"
	"github.com/spf13/cobra"
)


// initCmd respresents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a Cli library for Go that empowers applications. This
application is a tool to generate the needed files to quickly create a Cobra
application.`,
	Run: initialize,
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings

	// Cobra supports Persistent Flags which will work for this command and all subcommands
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle" )

}

func initialize(cmd *cobra.Command, args []string) {
	
	if _, err := os.Stat("./pkvkey.json"); err == nil {
		fmt.Printf("File '%s' already exists!\n",keyfile)
		fmt.Println("Cowardly refusing to overwrite it…")
		os.Exit(1)
	}

	k := generate.NewKey()
	j, err := json.MarshalIndent(k, " ", "  ")

	if err != nil {
		fmt.Printf("Error while generating JSON for key")
		os.Exit(1)
	}
	ioutil.WriteFile("pkvkey.json", j, 0600)
}
