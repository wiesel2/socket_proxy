package test

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"socket_proxy/src/dbb.com/lvfeng/utils"
	"testing"
)

func TestDefaultCfg(t *testing.T){
	var c = make(map[interface{}]interface{})
	defaultCfgPath := utils.DefaultCFGPath()
	yamlFile, err := ioutil.ReadFile(defaultCfgPath)
	if err != nil{
		panic(errors.New(fmt.Sprintf("Default cfg load failed, error: %s", err)))
	}
	err = yaml.UnmarshalStrict(yamlFile, &c)

	fmt.Printf("%v",c)
	fmt.Printf("%v", c )
}