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
	tools.CreateFile(fixDir + path)
	err := ioutil.WriteFile(fixDir+allPath, []byte(json["content"].(string)), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
