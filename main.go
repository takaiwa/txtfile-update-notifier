package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	filename := "./sample.txt"
	var inTE, outTE *walk.TextEdit

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

				outTE.SetText(string(b))
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}

	var mw *walk.MainWindow

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "txtfile-update-notifier",
		Size:     Size{300, 200},
		MaxSize:  Size{400, 300},
		MinSize:  Size{100, 100},
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{AssignTo: &outTE, ReadOnly: true},
			TextEdit{AssignTo: &inTE},
			PushButton{
				Text: "ADD",
				OnClicked: func() {
					file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					fmt.Fprintln(file, inTE.Text()+"\n") //ファイルに書き込み
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
	<-done
}
