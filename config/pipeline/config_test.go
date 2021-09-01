package pipeline

import (
	"fmt"
	"testing"
)

func TestWriteTo(t *testing.T) {
	config,_ := GetConfigFrom("./local_config.yaml")
	WriteConfigTo(config, "./local.yaml")
	config,_ = GetConfigFrom("./local.yaml")
	fmt.Println(config.Devops.Gitlab)
}
