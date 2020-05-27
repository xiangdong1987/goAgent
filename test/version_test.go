package test

import (
	"nodeAgent/fun"
	"testing"
)

func TestGetVersion(t *testing.T) {
	fun.GetCommit("\\Users\\xiangdd\\go\\src\\goAgent\\")
}

func TestUpdateVersion(t *testing.T) {
	fun.UpdateCode("\\Users\\xiangdd\\go\\src\\goAgent\\")
}

func TestCompareVersion(t *testing.T) {
	fun.CompareCode("\\Users\\xiangdd\\go\\src\\goAgent\\", "05cc39097af99b00c771205077fe73514c32c322", "464d8e8e24e83397737f3cee4562212f9d13fbae")
}

func TestPackageFiles(t *testing.T) {
	files := fun.CompareCode("/Users/xiangdd/go/src/goAgent/", "05cc39097af99b00c771205077fe73514c32c322", "464d8e8e24e83397737f3cee4562212f9d13fbae")
	fun.MoveFiles("/Users/xiangdd/go/src/goAgent/", files, "code_version/tmp/")
	fun.PackageFiles("/Users/xiangdd/go/src/goAgent/code_version/", "tmp", "test.tar.gz")
	//filesOs := fun.GetFilesByName("/Users/xiangdd/go/src/goAgent/code_version/tmp/",files)
	//err:=fun.Tar2("/Users/xiangdd/go/src/goAgent/code_version/tmp", "/Users/xiangdd/go/src/goAgent/code_version/test.tar",false)
	//fun.CreateTar("../code_version/tmp", "/Users/xiangdd/go/src/goAgent/code_version/test.tar.gz")
}
