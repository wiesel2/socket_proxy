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

func LoadDefaultConfig()(error, Cfg){
	var defaultCfg Cfg  = Cfg{}
	defaultCfgPath := utils.DefaultCFGPath()
	yamlFile, err := ioutil.ReadFile(defaultCfgPath)
	if err != nil{
		return errors.New(fmt.Sprintf("Default cfg load failed, error: %s", err)), defaultCfg
	}
	err = yaml.UnmarshalStrict(yamlFile, &defaultCfg)
	if err != nil{
		return errors.New(fmt.Sprintf("Default cfg un marshal failed, error: %s", err)), defaultCfg
	}
	fmt.Printf("Default Config: %v", defaultCfg)
	return nil, defaultCfg
}


func init(){
	once.Do(func() {
		var err error
		err, DefaultCfg = LoadDefaultConfig()
		if err != nil{
			panic(err)
		}
	})
}