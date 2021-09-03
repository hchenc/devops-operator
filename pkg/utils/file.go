package utils

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"reflect"
)

func WriteConfigTo(obj interface{}, fpath string) error {
	data, _ := yaml.Marshal(obj)
	err := ioutil.WriteFile(fpath, data, 0666)
	return err
}

func GetDataFrom(fpath string, obj interface{}) error {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	vObj := reflect.ValueOf(obj)

	if err := yaml.Unmarshal(data, vObj); err != nil {
		return err
	}

	return nil
}
