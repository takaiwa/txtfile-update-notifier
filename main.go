package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	filename := "./sample.txt"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	// ファイルパスからファイル情報取得
	fileinfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fsize := fileinfo.Size()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				file, err := os.Open(filename)
				if err != nil {
					fmt.Printf(err.Error())
				}
				defer file.Close()

				file.Seek(fsize, 0)
				b, err := ioutil.ReadAll(file)
				if err != nil {
					panic(err)
				}
				fsize = fsize + int64(len(b))
				fmt.Println("update:", string(b))
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
