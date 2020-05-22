package fun

import (
	"io/ioutil"
	"log"
	"nodeAgent/tools"
	"path/filepath"
	"runtime"
)

func Save(params string) {
	fixDir := ""
	switch runtime.GOOS {
	case "darwin":
	case "windows":
		fixDir = "/data/config/"
	case "linux":
	}
	json := tools.JsonDecode([]byte(params))
	allPath := json["path"].(string)
	fileName := filepath.Base(allPath)
	path := json["path"].(string)[0 : len(allPath)-len(fileName)]
	//判断路径是否存在
	tools.CreateFile(fixDir + path)
	err := ioutil.WriteFile(fixDir+allPath, []byte(json["content"].(string)), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
