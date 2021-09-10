module github.com/hchenc/devops-operator

go 1.15

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/hchenc/application v1.0.1
	github.com/hchenc/go-harbor v0.0.3
	github.com/hchenc/pager v0.0.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/xanzy/go-gitlab v0.50.3
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	sigs.k8s.io/controller-runtime v0.6.0
)
