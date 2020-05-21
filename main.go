package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	flag.Parse()
	filename := "./sample.txt"
	if len(flag.Args()) != 0 {
		filename = flag.Arg(0)
	}

	var inTE, outTE *walk.TextEdit
	var fileL *walk.LinkLabel
	var timeL *walk.Label

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	defer close(done)

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

				file.Seek(fsize, 0)
				b, err := ioutil.ReadAll(file)
				if err != nil {
					panic(err)
				}
				fsize = fsize + int64(len(b))
				file.Close()

				outTE.SetText(string(b))
				inTE.SetText("") // クリア
				t := time.Now()
				timeL.SetText(t.Format("15:04"))
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
			Composite{
				Layout: HBox{},
				Children: []Widget{
					LinkLabel{
						AssignTo: &fileL,
						Text:     `<a id="filename" href="">` + filename + `</a>`,
						OnLinkActivated: func(link *walk.LinkLabelLink) {
							wincmd := "/C start " + filename
							err := exec.Command("cmd", wincmd).Run()
							if err != nil {
								panic(err)
							}
						},
					},
					Label{
						AssignTo:      &timeL,
						TextAlignment: AlignFar,
						Text:          "",
					},
				},
			},
			TextEdit{AssignTo: &outTE, ReadOnly: true},
			TextEdit{AssignTo: &inTE},
			PushButton{
				Text: "ADD",
				OnClicked: func() {
					file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Fprintln(file, inTE.Text()+"\n") //ファイルに書き込み
					file.Close()
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}
