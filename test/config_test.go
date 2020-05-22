package test

import (
	"nodeAgent/fun"
	"testing"
)

func TestConfigSave(t *testing.T) {
	fun.Save(" {\"content\":\"{\\\"a\\\":1}\",\"path\":\"server/a.conf\"}")
}
