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
	controller "github.com/hchenc/devops-operator/pkg/controllers"
	"github.com/hchenc/devops-operator/pkg/models"
	"github.com/hchenc/devops-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/spf13/cobra"
)

var (
	kubeconfig           string
	pipelineConfig       *models.Config
	scheme               = runtime.NewScheme()
	metricsAddr          string
	enableLeaderElection bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// use the current context in kubeconfig

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		//setup config
		err = utils.GetDataFrom(cfgFile, pipelineConfig)
		if err != nil {
			panic(err.Error())
		}

		dc := &controller.DevopsClientet{}
		dc.Complete(config)
		//setup controller through config instance
		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme:             scheme,
			MetricsBindAddress: metricsAddr,
			Port:               9443,
			LeaderElection:     enableLeaderElection,
			LeaderElectionID:   "5e352c21.efunds.com",
		})
		c, _ := controller.New(dc, mgr, pipelineConfig)
		//run controller
		fmt.Println("run called")
		c.Reconcile(ctrl.SetupSignalHandler())
	},
}

func init() {
	configPath := filepath.Join(DevOpsOperatorDir, ConfigFileName)

	if home := homeDir(); home != "" {
		runCmd.Flags().StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		runCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	runCmd.Flags().StringVarP(&cfgFile, "config-path", "c", configPath, "config file path to load")

	rootCmd.AddCommand(runCmd)

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
