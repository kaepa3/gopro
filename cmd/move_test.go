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
		case <-ch:
			counter += 1
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

// TestMove
func TestMove(t *testing.T) {
	fileName := "dummy.file"
	from := filepath.Join(conf.GoproRoot, fileName)
	to := filepath.Join(conf.SavePath, fileName)
	f, err := os.Create(from)
	defer f.Close()
	if err != nil {
		t.Error("create error")
	}
	err = moveProcess(from, to)
	if err != nil {
		t.Error(err.Error())
	}
}
