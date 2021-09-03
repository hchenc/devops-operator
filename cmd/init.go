/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/hchenc/devops-operator/config/pipeline"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	host     = os.Getenv("HOST")
	port     = os.Getenv("PORT")
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
)


func init() {
	outputConfigPath := filepath.Join(DevOpsOperatorDir, ConfigFileName)

	initCmd.Flags().StringVarP(&cfgFile, "output-config-path", "o",outputConfigPath,"config file path to place (default is $HOME/devops-operator.yaml)")

	rootCmd.AddCommand(initCmd)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("init called")
		if err := setUp(); err !=nil{
			fmt.Println("init error")
		}

	},
}

func setUp() error {
	config := &pipeline.Config{}
	config.Devops.Gitlab.Version = "ee"
	config.Devops.Gitlab.User = user
	config.Devops.Gitlab.Host = host
	config.Devops.Gitlab.Port = port
	config.Devops.Gitlab.Password = password



	config.Devops.Pipelines = []pipeline.Pipelines{
		{
			Pipeline: "java",
			Template: "Spring",
			Ci: "devops/devops/-/raw/master/java.yml",
		},
	}
	return pipeline.WriteConfigTo(config, cfgFile)
}

