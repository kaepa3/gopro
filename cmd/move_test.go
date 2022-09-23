package cmd

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// before
	conf.GoproRoot = "./test"
	conf.SavePath = "./save"
	pathAry := []string{conf.GoproRoot, conf.SavePath}
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

// TestSearch
func TestSearch(t *testing.T) {
	t.Log("zero")
	done := make(chan interface{})
	ch := make(chan FInfo)
	fName := "text.txt"
	srcPath := filepath.Join(conf.GoproRoot, fName)
	f, err := os.Create(srcPath)
	if err != nil {
		t.Errorf("file not created:%s\n", srcPath)
	}
	f.Close()
	go searchDir(conf.GoproRoot, ch, done)
	counter := 0
looper:
	for {
		select {
		case val := <-ch:
			counter += 1
			t.Logf("file input:%s\n", val.FileInfo.Name())
		case <-done:
			break looper
		case <-time.After(10 * time.Second):
			t.Error("time out!!!")
			break looper
		}
	}
	if counter <= 0 {
		t.Error("not counted")
	}
}
