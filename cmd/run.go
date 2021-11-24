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
	controller "github.com/hchenc/devops-operator/pkg/controllers"
	"github.com/hchenc/devops-operator/pkg/models"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/spf13/cobra"
)

var (
	inputCfgFile         string
	kubeconfig           string
	homePath             string
	err                  error
	pipelineConfig       = &models.Config{}
	scheme               = runtime.NewScheme()
	metricsAddr          string
	enableLeaderElection bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run Devops Controller to watch kubernetes and kubesphere resource",
	Run: func(cmd *cobra.Command, args []string) {
		// use the current context in kubeconfig
		initConfig()

		var config *rest.Config

		config, err := rest.InClusterConfig()
		if err != nil {
			if err == rest.ErrNotInCluster {
				config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}

		cs := models.NewForConfigOrDie(config, pipelineConfig)
		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
			Scheme:             scheme,
			MetricsBindAddress: metricsAddr,
			Port:               9443,
			LeaderElection:     enableLeaderElection,
			LeaderElectionID:   "5e352c21.efunds.com",
		})
		c := controller.NewControllerOrDie(cs, mgr)
		c.Reconcile(ctrl.SetupSignalHandler())
	},
}

func init() {
	homePath, err = homedir.Dir()
	cobra.CheckErr(err)

	//cobra.OnInitialize(initConfig)

	runCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", filepath.Join(homePath, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	runCmd.Flags().StringVarP(&inputCfgFile, "config-path", "c", filepath.Join(homePath, ".devops-operator.yaml"), "(optional) config file path to load")

	rootCmd.AddCommand(runCmd)

}

func initConfig() {

	if kubeconfig != "" {
		viper.SetConfigFile(kubeconfig)
	}

	if inputCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(inputCfgFile)
	} else {
		viper.SetConfigType("yaml")
		viper.AddConfigPath(homePath)
		viper.SetConfigName(".devops-operator.yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// If a config file is found, read it in.
	if err := viper.Unmarshal(pipelineConfig); err != nil {
		panic(err)
	}

}
