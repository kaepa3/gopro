/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: move,
}

func init() {
	moveCmd.Flags().BoolP("del", "d", false, "non delete")
	rootCmd.AddCommand(moveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type FInfo struct {
	FileInfo fs.FileInfo
	Root     string
}

// move file cmd
func move(cmd *cobra.Command, args []string) {
	if err := canProc(); err != nil {
		fmt.Printf("error :%s\n", err.Error())
		return
	}
	done := make(chan interface{})
	ch := make(chan FInfo)
	go searchDir(conf.GoproRoot, ch, done)
	procMessage(ch, done)
}

func procMessage(ch <-chan FInfo, loopDone <-chan interface{}) {
	wg := sync.WaitGroup{}

loopLabel:
	for {
		select {
		case f := <-ch:
			wg.Add(1)
			go moving(&wg, f)
		case <-loopDone:
			break loopLabel
		}
	}
	wg.Wait()
}

// moving
func moving(wg *sync.WaitGroup, f FInfo) {
	defer wg.Done()
	folderName := createFolderName(f.FileInfo.ModTime())
	savePath := filepath.Join(conf.SavePath, folderName)
	createFolderIfNeed(savePath)

	from := filepath.Join(f.Root, f.FileInfo.Name())
	to := filepath.Join(savePath, f.FileInfo.Name())
	fmt.Printf("%s -> %s\n", from, to)
	if err := moveProcess(from, to); err != nil {
		fmt.Println(err.Error())
	}
}

// createFolderIfNeed
func createFolderIfNeed(folderPath string) {
	info, err := os.Stat(folderPath)
	if err == nil {
		if info.IsDir() {
			fmt.Printf("same file exits:%s\n", folderPath)
		}
	} else {
		os.Mkdir(folderPath, 0777)
	}
}

// createFolderName
func createFolderName(t time.Time) string {
	return t.Format("2006-01-02")
}

// move`Process
func moveProcess(from string, to string) error {
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

// searchDir
func searchDir(path string, ch chan<- FInfo, done chan interface{}) {
	defer close(done)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		info := FInfo{FileInfo: f, Root: path}
		ch <- info
	}
}

func canProc() error {
	if !isExistDir(conf.GoproRoot) {
		return errors.New(fmt.Sprintf("gopro not exist:%s\n", conf.GoproRoot))

	}
	if !isExistDir(conf.SavePath) {
		return errors.New(fmt.Sprintf("save dir not exist:%s\n", conf.SavePath))
	}
	return nil
}

func isExistDir(path string) bool {
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		return false
	}
	return true
}
