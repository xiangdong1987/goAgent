package fun

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"nodeAgent/tools"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

/**
获取提交
*/
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
	fmt.Println(commit)
	return commit
}

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

func CompareCode(path string, commitA string, commitB string) []string {
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
	//fmt.Println(res)
	files := ParseChange(path, res)
	return files
}

func ParseChange(path string, result string) []string {
	var res []string
	var fileName string
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
		res = append(res, fileName)
		fmt.Println(fileName)
	}
	return res
}

func PackageFiles(files []string) {

}

//压缩 使用gzip压缩成tar.gz
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := compress(file, "", tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压 tar.gz
func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}
func GetFilesByName(path string, files []string) []*os.File {
	var filesOs []*os.File
	for _, filename := range files {
		fmt.Println(path + filename)
		file, err := createFile(path + filename)
		if err == nil {
			filesOs = append(filesOs, file)
		} else {
			log.Println(err)
		}
	}
	return filesOs
}
func createFile(name string) (*os.File, error) {
	fmt.Println(name)
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "\\")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
