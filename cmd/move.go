/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kaepa3/go-mtpfs/mtp"
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

type MtpInfo struct {
	MtpObjInfo *mtp.ObjectInfo
	Handle     uint32
}

// move file cmd
func move(cmd *cobra.Command, args []string) {
	err, dev, delFlg := initCommand(cmd)
	if err != nil {
		fmt.Printf("error :%s\n", err.Error())
		return
	}
	defer dev.Close()
	fmt.Println(dev.ID())

	searchDir(dev, delFlg)
}

// initCommand
func initCommand(cmd *cobra.Command) (error, *mtp.Device, bool) {
	v, err := cmd.Flags().GetBool("del")
	if err != nil {
		return err, nil, v
	}
	if err := canProc(); err != nil {
		return err, nil, v
	}
	dev, err := GetGopro()
	if err != nil {
		return err, nil, v
	}
	return nil, dev, v
}

// createFolderIfNeed
func createFolderIfNeed(folderPath string) {
	info, err := os.Stat(folderPath)
	if err == nil {
		if !info.IsDir() {
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

// moveProcess
func writeFile(dev *mtp.Device, handle uint32, name string) error {
	fs, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fs.Close()
	writer := bufio.NewWriter(fs)
	err = dev.GetObject(handle, writer)
	if err != nil {
		return err
	}
	return nil
}

// searchDir
func searchDir(dev *mtp.Device, isDel bool) {
	sids := mtp.Uint32Array{}
	err := dev.GetStorageIDs(&sids)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, id := range sids.Values {
		err, handles := getHandles(dev, id)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			transferObject(dev, handles)
		}
	}
}

func transferObject(dev *mtp.Device, handles []uint32) {
	for _, handle := range handles {
		var oi mtp.ObjectInfo
		dev.GetObjectInfo(handle, &oi)
		if strings.Contains(oi.Filename, "MP4") {
			folderName := createFolderName(oi.ModificationDate)
			savePath := filepath.Join(conf.SavePath, folderName)
			createFolderIfNeed(savePath)
			to := filepath.Join("save", oi.Filename)
			fmt.Printf("%s -> %s\n", oi.Filename, to)
			writeFile(dev, handle, to)
		}
	}
}

// storagehandle
func getHandles(dev *mtp.Device, id uint32) (error, []uint32) {
	hs := mtp.Uint32Array{}
	err := dev.GetObjectHandles(id, 0x0, 0x0, &hs)
	if err != nil {
		return err, nil
	}
	return nil, hs.Values
}

// GetGopro
func GetGopro() (*mtp.Device, error) {
	dev, err := mtp.SelectDevice("GoPro")
	if err != nil {
		return nil, err
	}
	dev.Configure()
	return dev, nil
}

func canProc() error {
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
