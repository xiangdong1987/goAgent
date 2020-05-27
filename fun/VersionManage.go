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
	"path/filepath"
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

//压缩 使用gzip压缩成tar.gz
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := compress(file, "test", tw)
		if err != nil {
			return err
		}
	}
	return nil
}
func CreateTar(fileSource, fileTarget string) error {
	//创建目录
	tarTarget, err := os.Create(fileTarget)
	if err != nil {
		//如果文件存在删除文件
		if err == os.ErrExist {
			if err := os.Remove(fileTarget); err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			fmt.Println(err)
			return err
		}
	}
	defer tarTarget.Close()
	tarWriter := tar.NewWriter(tarTarget)
	sFileInfo, err := os.Stat(fileSource)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !sFileInfo.IsDir() {
		return tarFile(fileTarget, fileSource, sFileInfo, tarWriter)
	} else {
		return tarFolder(fileSource, tarWriter)
	}
	return nil
}

func tarFile(directory string, filesource string, sfileInfo os.FileInfo, tarwriter *tar.Writer) error {
	sfile, err := os.Open(filesource)
	if err != nil {
		panic(err)
		return err
	}
	defer sfile.Close()
	header, err := tar.FileInfoHeader(sfileInfo, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	header.Name = directory
	err = tarwriter.WriteHeader(header)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if _, err = io.Copy(tarwriter, sfile); err != nil {
		fmt.Println(err)
		panic(err)
		return err
	}
	return nil
}

func tarFolder(directory string, tarwriter *tar.Writer) error {
	var baseFolder string = filepath.Base(directory)
	//fmt.Println(baseFolder)
	return filepath.Walk(directory, func(targetpath string, file os.FileInfo, err error) error {
		//read the file failure
		if file == nil {
			panic(err)
			return err
		}
		if file.IsDir() {
			// information of file or folder
			header, err := tar.FileInfoHeader(file, "")
			if err != nil {
				return err
			}
			header.Name = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			fmt.Println(header.Name)
			if err = tarwriter.WriteHeader(header); err != nil {
				return err
			}
			os.Mkdir(strings.TrimPrefix(baseFolder, file.Name()), os.ModeDir)
			return nil
		} else {
			//baseFolder is the tar file path
			var fileFolder = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			return tarFile(fileFolder, targetpath, file, tarwriter)
		}
	})
}
func compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		fmt.Println(fileInfos)
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

func createFile(name string) (*os.File, error) {
	fmt.Println(name)
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
