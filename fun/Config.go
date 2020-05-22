package fun

import (
	"fmt"
	"io/ioutil"
	"log"
	"nodeAgent/tools"
	"runtime"
)

func Save(params string) {
	fmt.Println(params)
	ficDir := ""
	switch runtime.GOOS {
	case "darwin":
	case "windows":
		ficDir = "d:"
	case "linux":
	}
	json := tools.JsonDecode([]byte(params))
	err := ioutil.WriteFile(ficDir+json["path"].(string), []byte(json["content"].(string)), 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(json["path"])
	fmt.Println(json)
}
