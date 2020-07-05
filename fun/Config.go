package fun

import (
	"io/ioutil"
	"log"
	"nodeAgent/tools"
	"os"
	"path/filepath"
	"runtime"
)

func Save(params string) {
	fixDir := ""
	switch runtime.GOOS {
	case "darwin":
	case "windows":
		//node 所在磁盘的目录
		fixDir = "/data/config/"
	case "linux":
		fixDir = "/data/config/"
	}
	json := tools.JsonDecode([]byte(params))
	allPath := json["path"].(string)
	fileName := filepath.Base(allPath)
	path := json["path"].(string)[0 : len(allPath)-len(fileName)]
	//判断路径是否存在并创建
	err := tools.CreateFile(fixDir + path)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(fixDir+allPath, []byte(json["content"].(string)), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateDir(path string) error {
	fileName := filepath.Base(path)
	filePath := path[0 : len(path)-len(fileName)]
	//判断路径是否存在并创建
	return tools.CreateFile(filePath)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
