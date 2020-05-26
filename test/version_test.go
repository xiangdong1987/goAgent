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
	filesOs := fun.GetFilesByName("/Users/xiangdd/go/src/goAgent/", files)
	fun.Compress(filesOs, "code/version/test.tar.gz")
}
