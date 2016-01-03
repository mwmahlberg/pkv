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
	//	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

const (
	keycheck = `package verify

func Key(key string,bl []uint64) error {
	return KeyPart(key,{{ .Idx }},{{ index .Iv 0 }},{{ index .Iv 1 }},{{ index .Iv 2 }},bl)
}
`
	keyCheckFileName = "./pkv/verify/k.go"
	)

var (
	k      int
	stdout bool
	tmpl   *template.Template
)

// genkeyCmd respresents the genkey command
var genCodeCmd = &cobra.Command{
	Use:   "gencode",
	Short: "generates the code needed to check the selected part of the key",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a Cli library for Go that empowers applications. This
application is a tool to generate the needed files to quickly create a Cobra
application.`,
	Run: genCode,
}

func init() {
	RootCmd.AddCommand(genCodeCmd)

	// Here you will define your flags and configuration settings

	// Cobra supports Persistent Flags which will work for this command and all subcommands
	// genkeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command is called directly
	// genkeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle" )
	//	genkeyCmd.Flags().StringP("seed","s","","Seed value")
	genCodeCmd.Flags().IntVarP(&k, "key", "k", 0, "key to check in generated code [1-4]")
	genCodeCmd.Flags().BoolVarP(&stdout, "print", "p", false, "print generated code instead of writing it to files")
	tmpl = template.Must(template.New("Key").Parse(keycheck))

}

type Keychecker struct {
	Iv  [3]uint8
	Idx int
}

func genCode(cmd *cobra.Command, args []string) {

	if k < 1 || k > 4 {
		cmd.Usage()
		os.Exit(1)
	}

	if stdout {
		data, _ := Asset("verify/key.go")
		fmt.Println(string(data))
	} else {
		RestoreAsset("pkv", "verify/key.go")
	}

	key := readKeyFile()

	kc := Keychecker{Idx: k - 1, Iv: key.Matrix[k-1]}

	var out *os.File
	var err error
	
	if stdout {
		out = os.Stdout
	} else {
		out, err = os.OpenFile(keyCheckFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC,0644)
		
		if err != nil {
			fmt.Printf("Can not write output file '%s': %v\n",keyCheckFileName,err)
			os.Exit(1)
		}
	}
	tmpl.ExecuteTemplate(out, "Key", kc)

}
