package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory(dir string) string {
	return substr(dir, 0, strings.LastIndex(dir, string(os.PathSeparator)))
}

func GetCurrentDirectory() string {
	_, file,_, ok := runtime.Caller(1)
	if ! ok{
		panic(errors.New("Can not get current file info"))
	}
	dir, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func ProxyRoot() string{
	curDir := GetCurrentDirectory()
	srcRoot := GetParentDirectory(GetParentDirectory(GetParentDirectory((curDir))))
	return srcRoot
}

func DefaultCFGPath() (path string){
	path = fmt.Sprintf( "%s%s%s%s%s", ProxyRoot(), string(os.PathSeparator),
		"config", string(os.PathSeparator), "default.yaml")
	return
}
