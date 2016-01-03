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
	"os"

	pkv "github.com/mwmahlberg/pkv/generate"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyfile = "./pkvkey.json"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "partkey",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your application. For example:

Cobra is a Cli library for Go that empowers applications. This
application is a tool to generate the needed files to quickly create a Cobra
application.`,
	// Uncomment the following line if your bare application has an action associated with it
	//	Run: func(cmd *cobra.Command, args []string) { },
}

//Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings
	// Cobra supports Persistent Flags which if defined here will be global for your application

//	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.partkey.yaml)")

	// Cobra also supports local flags which will only run when this action is called directly

}

// Read in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".partkey") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func readKeyFile() (key *pkv.PartialKey) {
	b, err := ioutil.ReadFile(keyfile)

	if err != nil {
		fmt.Printf("Could not read keyfile '%s'\n", keyfile)
		fmt.Printf("Have you called 'init'?\n\n")
		RootCmd.Usage()
		os.Exit(1)
	}

	key = &pkv.PartialKey{}

	if err := json.Unmarshal(b, key); err != nil {
		fmt.Printf("Error while reading key data: %v\n", err)
		os.Exit(1)
	}
	
	return
}
