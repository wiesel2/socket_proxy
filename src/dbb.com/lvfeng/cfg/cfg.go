package cfg

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"socket_proxy/src/dbb.com/lvfeng/utils"
	"sync"
)

//const (
//	CFG_T_INT = iota
//	CFG_T_FLOAT
//	CFG_T_BOOL
//	CFG_T_STR
//	CFG_T_LIST 	// slice
//	CFG_T_MAP	// dict, json
//)

type SrvCfg struct {
	LocalPort int `yaml:"LocalPort"`
	RemoteHost string `yaml:"RemoteHost"`
	RemotePort int `yaml:"RemotePort"`
}

type conCfg struct {
	MaxRate int  `yaml:"MaxRate"`// KB/s
	MaxConnCount int `yaml:"MaxConnCount"` // max connection count
}

type Cfg struct {
	ServerConfig SrvCfg `yaml:"ServerConfig"`
	ConCfg conCfg `yaml:"ConCfg"`
}

type BlackList struct {
	HOSTS []string
}

var DefaultCfg Cfg
var once sync.Once




func init(){
	once.Do(func() {
		DefaultCfg = Cfg{}
		defaultCfgPath := utils.DefaultCFGPath()
		yamlFile, err := ioutil.ReadFile(defaultCfgPath)
		if err != nil{
			panic(errors.New(fmt.Sprintf("Default cfg load failed, error: %s", err)))
		}
		err = yaml.UnmarshalStrict(yamlFile, DefaultCfg)
		if err != nil{
			panic(errors.New(fmt.Sprintf("Default cfg un marshal failed, error: %s", err)))
		}
		fmt.Printf("Default Config: %v", DefaultCfg)
	})
}