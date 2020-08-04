package main

import (
	"bitbucket.com/ZackArts/time-checker/controllers"
	"github.com/gotk3/gotk3/gtk"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	WindowName = "window"
	ButtonName = "button"
	UIMain     = "main.glade"
)

// TODO: Add Export Functionality (1w, 2w, etc)
// TODO: Beautify UI
// TODO: Misc bug fixes, etc

func main() {
	DEBUG_VAR, _ := strconv.ParseBool(os.Getenv("ENV"))
	path := "test.db"
	if DEBUG_VAR != true {
		path = "/usr/share/time-checker/data.db"
	}
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	client := controllers.Client{DB: db}
	client.Initialize()

	gtk.Init(&os.Args)

	bldr, err := getBuilder(UIMain)
	if err != nil {
		panic(err)
	}

	window, err := getWindow(bldr)
	if err != nil {
		panic(err)
	}

	window.SetTitle("Time Checker")
	window.SetDefaultSize(400, 300)
	_, err = window.Connect("destroy", func() {
		gtk.MainQuit()
	})
	if err != nil {
		panic(err)
	}

	window.ShowAll()

	button, err := getButton(bldr, ButtonName)
	if err != nil {
		panic(err)
	}

	button2, err := getButton(bldr, "button2")
	if err != nil {
		panic(err)
	}
	button3, err := getButton(bldr, "button3")
	if err != nil {
		panic(err)
	}
	_, err = button.Connect("clicked", func() {
		log.Printf("started the time at %s", time.Now().String())
		client.Start()
	})
	if err != nil {
		panic(err)
	}

	_, err = button2.Connect("clicked", func() {
		category, err := getInputVal(bldr, "entry2")
		if err != nil {
			panic(err)
		}
		comment, err := getInputVal(bldr, "entry1")
		if err != nil {
			panic(err)
		}
		err = client.End(category, comment)
		if err != nil {
			log.Fatalf("an error occured with sqlite3: %v", err)
		}
	})
	if err != nil {
		panic(err)
	}

	_, err = button3.Connect("clicked", func() {
		category, err := getInputVal(bldr, "entry1")
		if err != nil {
			panic(err)
		}
		client.Export(category)
	})

	gtk.Main()
}

func getWindow(bldr *gtk.Builder) (*gtk.Window, error) {
	obj, err := bldr.GetObject(WindowName)
	if err != nil {
		return nil, err
	}

	window, ok := obj.(*gtk.Window)
	if !ok {
		return nil, err
	}

	return window, nil
}

func getButton(bldr *gtk.Builder, buttonName string) (*gtk.Button, error) {
	obj, err := bldr.GetObject(buttonName)
	if err != nil {
		return nil, err
	}

	button, ok := obj.(*gtk.Button)
	if !ok {
		return nil, err
	}

	return button, nil
}

func getInputVal(bldr *gtk.Builder, inputName string) (string, error) {
	obj, err := bldr.GetObject(inputName)
	if err != nil {
		panic(err)
	}

	input, ok := obj.(*gtk.Entry)
	if !ok {
		return "", err
	}

	return input.GetText()
}

func getBuilder(uiMain string) (*gtk.Builder, error) {
	b, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}

	if uiMain != "" {
		err := b.AddFromFile(uiMain)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}
