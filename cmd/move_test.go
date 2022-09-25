package cmd

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// before
	conf.SavePath = "./save"
	pathAry := []string{conf.SavePath}
	for _, path := range pathAry {
		err := os.Mkdir(path, 0777)
		if err != nil {
			log.Println(err.Error())
		}
	}
	code := m.Run()
	// after
	for _, path := range pathAry {
		os.RemoveAll(path)
	}
	os.Exit(code)
}

// TestCanProc
func TestCanProc(t *testing.T) {
	err := canProc()
	if err != nil {
		t.Errorf("can proc error:%v", conf)
	}
}
