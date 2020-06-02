package fun

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"nodeAgent/tools"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//获取当前提交
func GetCommit(path string) string {
	//判断路径是否存在
	if !tools.IsExist(path) {
		log.Fatal("program is not Exist")
	}
	//打开目录
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	commit := out.String()
	commit = strings.Replace(commit, "\n", "", -1)
	return commit
}

//更新代码
func UpdateCode(path string) {
	//判断路径是否存在
	if !tools.IsExist(path) {
		log.Fatal("program is not Exist")
	}
	//打开目录
	cmd := exec.Command("git", "pull")
	cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	res := out.String()
	fmt.Println(res)
}

//比较代码
func CompareCode(path string, commitA string, commitB string) ([]string, []string) {
	//判断路径是否存在
	if !tools.IsExist(path) {
		log.Fatal("program is not Exist")
	}
	//打开目录
	cmd := exec.Command("git", "diff", "--stat", commitA, commitB)
	cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	res := out.String()
	files, deleteFiles := ParseChange(res)
	fmt.Println(files, deleteFiles)
	return files, deleteFiles
}

//解析变化
func ParseChange(result string) ([]string, []string) {
	var res []string
	var fileName string
	var deleteFile []string
	strArray := tools.ExplodeStr(result, "\n")
	for _, v := range strArray {
		pos := strings.Index(v, "|")
		if pos < 0 {
			continue
		}
		tmpArray := tools.ExplodeStr(v, "|")
		//匹配一个或多个空白符的正则表达式
		reg := regexp.MustCompile("\\s+")
		fileName = reg.ReplaceAllString(tmpArray[0], "")
		typeArray := tools.ExplodeStr(tmpArray[1], " ")
		//判断删除
		if typeArray[2] == "-" {
			deleteFile = append(deleteFile, fileName)
		} else {
			res = append(res, fileName)
		}
	}
	return res, deleteFile
}

//打包代码更新
func PackageFiles(path, source, desc string) error {
	//判断路径是否存在
	if !tools.IsExist(source) {
		log.Fatal("program is not Exist")
	}
	//打开目录
	cmd := exec.Command("tar", "-zcvf", desc, source)
	cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	res := out.String()
	fmt.Println(res)
	return nil
}

//移动文件
func MoveFiles(path string, files []string, des string) {
	for _, filename := range files {
		//fmt.Println(path + filename)
		//判断路径是否存在并创建
		err := os.MkdirAll(string([]rune(path + des + filename)[0:strings.LastIndex(path+des+filename, "/")]), 0755)
		if err != nil {
			log.Println(err)
		}
		_, err = copyFile(path+filename, path+des+filename)
		if err != nil {
			log.Println(err)
		}
	}
}

//文件拷贝
func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
