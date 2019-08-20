package cfg

import (
	"dbb.com/lvfeng/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	yaml "gopkg.in/yaml.v2"
)

const (
	CFG_T_INT = iota
	CFG_T_FLOAT
	CFG_T_BOOL
	CFG_T_STR
	CFG_T_LIST 	// slice
	CFG_T_MAP	// dict, json
)

type SrvCfg struct {
	LocalPort int `yaml:"LocalPort"`
	RemoteHost string `yaml:"RemoteHost"`
	RemotePort int `yaml:"RemotePort"`
}

func (sc SrvCfg)Network() string{
	return "tcp"
}

func (sc SrvCfg)String() string{
	return sc.RemoteHost + ":" + strconv.Itoa(sc.RemotePort)
}


type ConCfg struct {
	MaxRate int  `yaml:"MaxRate"`// KB/s
	MaxConnCount int `yaml:"MaxConnCount"` // max connection count
}

type Cfg struct {
	SrvCfg SrvCfg
	ConCfg ConCfg
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
		yaml
	})
}