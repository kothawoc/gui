package main

import (
	"embed"
	"flag"
	"log/slog"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"

	"github.com/kothawoc/kothawoc"
)

//go:embed locales
var localeFS embed.FS

var kc *kothawoc.Client

var myApp fyne.App
var mainWindow fyne.Window

var parent *container.Split
var leftbox *fyne.Container

var navList *fyne.Container
var mainFeed *widget.List
var addGroup *fyne.Container
var postForm *widget.Form

var content *fyne.Container

var path = flag.String("path", os.Getenv("PWD")+"/data", "Path to data directires")
var port = flag.Int("port", 1119, "Default NNTP port to listen on")

func setContent(content fyne.CanvasObject) {
	updateGroupsList() // lol
	parent.Trailing = content
	parent.Trailing.Refresh()
	parent.Refresh()
}

//var navContainer, contentContainer *fyne.CanvasObject

var ScrollReset func()

var ADB *ArticleDB

func main() {
	flag.Parse()

	slog.Info("start new client")

	lang.AddTranslationsFS(localeFS, "locales")

	myApp = app.NewWithID("io.github.kothawoc.gui")

	configPath := "/../kothawoc/data"

	configPath = myApp.Storage().RootURI().Path()

	slog.Info("Config", "path", configPath)
	if path != nil {
		configPath = *path
	}
	slog.Info("Config", "path", configPath)

	// panic("test config")
	myWindow := myApp.NewWindow(lang.L("App: NNTP/TOR"))
	mainWindow = myWindow

	//kc = kothawoc.NewClient(os.Getenv("PWD") + "/../kothawoc/data")
	//slog.Info("Done new client start dial", "error", err)
	KC, _ := kothawoc.NewClient(configPath, *port)
	kc = KC

	ADB, _ = NewArticleDB(configPath)

	kc.Dial()

	slog.Info("Done done dial start app")
	content = container.NewVBox()

	//displayHome(content)
	displayWelcome()
	navList = newNavlist()
	mainFeed = newMainFeed()
	postForm = newPost()
	addGroup = newAddGroup()

	navContainer := container.NewScroll(navList)
	contentContainer := container.NewScroll(content)
	ScrollReset = func() {
		contentContainer.Offset.X = 0
		contentContainer.Offset.Y = 0
	}

	split := container.NewHSplit(navContainer, contentContainer)
	parent = split
	contentContainer.Offset.X = 0
	contentContainer.Offset.Y = 0

	split.Offset = 0.2 // Adjust the initial size of the left pane

	// Set up the main content with the split container
	myWindow.SetContent(split)

	myWindow.Resize(fyne.NewSize(1024, 768))

	myWindow.Show()

	//	})
	myApp.Run()
	//myWindow.ShowAndRun()
}
