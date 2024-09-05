package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"net/textproto"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"

	"github.com/kothawoc/go-nntp"
	"github.com/kothawoc/kothawoc"
	"github.com/kothawoc/kothawoc/pkg/messages"
)

//go:embed locales
var localeFS embed.FS

var kc *kothawoc.Client

var mainWindow fyne.Window
var leftbox *fyne.Container
var content *fyne.Container

var path = flag.String("path", os.Getenv("PWD")+"/data", "Path to data directires")
var port = flag.Int("port", 1119, "Default NNTP port to listen on")

func displayHome(content *fyne.Container) {
	content.RemoveAll()
	label := widget.NewLabel("Select an item from the navigation pane")
	rt := widget.NewRichTextFromMarkdown(lang.L("App: Welcome Message"))
	rt.Wrapping = fyne.TextWrapWord
	/*
			rt := widget.NewRichTextFromMarkdown(`
		# What a load of shit goes here
		This is a master title

		## these are sub titles
		## and another

		### even deeper
		#### deeper still

		* Will this work?
		  + how about sub lists
		  + this will be cool
		* I assume it will
		* this is stars


		+ this one is plus
		  - or lets try to
		  - go down the rabbit hole
		    * and make it a level
			* deeper
		+ signs, I wonder if it also
		+ accepts others


		- this is another, how
		  + how about sub lists
		  + this will be cool
		- about this, does it
		- work too?

		## lower

		And is this ok?
		1. [ ] What do you think?
		2. [/] This is another?
		3. [?] Another
		4. [x] and again...


		> I'm not so **sure** about it.
		> Seems *very* limited to me.
		> I'm not even sure this quote thigng is working, is it that I
		> chose to add some word formatting mid line?

		### higher too

		What a load of crap goes here.
			`

	*/

	content.Add(rt)

	content.Add(label)
}

func newPost(content *fyne.Container) {

	groupsEntry := widget.NewEntry()

	subjectEntry := widget.NewEntry()

	editor := widget.NewMultiLineEntry()
	editor.SetMinRowsVisible(8)

	form := &widget.Form{
		Items: []*widget.FormItem{
			//		{Text: "Name", Widget: name, HintText: "Your full name"},
			//	    {Text: "Email", Widget: email, HintText: "A valid email address"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Form for: " + groupsEntry.Text,
				Content: editor.Text,
			})
			mail := messages.NewMessageTool()
			mail.Article.Header.Add("Newsgroups", groupsEntry.Text)
			mail.Article.Header.Add("Subject", subjectEntry.Text)
			//	mail.Article.Header.Add("Content-Type", "multipart/mixed; boundary=\"nxtprt\"")
			mail.Article.Header.Add("Content-Type", "text/plain")
			mail.Article.Header.Add("Content-Transfer-Encoding", "8bit")
			mail.Preamble = editor.Text
			mail.Parts = []messages.MimePart{
				{Header: textproto.MIMEHeader{
					"Content-Type": []string{"text/plain"},
				}, Content: []byte(editor.Text)},
				//		{Header: textproto.MIMEHeader{
				//			"Content-Type": []string{"text/markdown"},
				//		}, Content: []byte(editor.Text)},
			}
			log.Printf("POSTING MAIL LIKE THIS: [%#v][%s]", mail, mail.RawMail())

			kc.NNTPclient.Post(strings.NewReader(mail.RawMail()))
			//kc.Post(mail)
		},
	}
	form.Append("News Groups", groupsEntry)
	form.Append("Subject", subjectEntry)
	form.Append("Message", editor)

	content.Add(form)

}

func displayMessage(content *fyne.Container, messageNumber string) {

	content.RemoveAll()

	label := widget.NewLabel("The content of: " + messageNumber)
	content.Add(label)

	//a, messageId, rawArticle, err := kc.NNTPclient.Article(messageNumber)
	_, _, rawArticle, err := kc.NNTPclient.Article(messageNumber)

	buf := make([]byte, 8129)
	n, err := rawArticle.Read(buf)
	rawMail := string(buf[:n])
	splitMail := (strings.SplitN(rawMail, "\r\n\r\n", 2))
	body := ""
	if len(splitMail) == 2 {
		body = splitMail[1]
	} else {
		splitMail := (strings.SplitN(rawMail, "\n\n", 2))

		if len(splitMail) == 2 {
			body = splitMail[1]
		} else {
			body = rawMail
		}
	}

	mailMsg, err := mail.ReadMessage(strings.NewReader(rawMail))
	if err != nil {
		log.Fatal(err)
	}

	article := &nntp.Article{
		Header: textproto.MIMEHeader(mailMsg.Header),
		Body:   mailMsg.Body,
		Bytes:  len([]byte(body)),
		Lines:  strings.Count(body, "\n"),
	}

	msg := messages.NewMessageToolFromArticle(article)
	content.Add(widget.NewLabel("From: " + msg.Article.Header.Get("From") + "\n" +
		"Subject: " + msg.Article.Header.Get("Subject") + "\n" +
		"Date: " + msg.Article.Header.Get("Date") + "\n" +
		"\n"))
	//content.Add(widget.NewLabel("Subject: " + msg.Article.Header.Get("Subject")))
	//content.Add(widget.NewLabel("Date: " + msg.Article.Header.Get("Date")))
	//content.Add(widget.NewLabel(" "))
	content.Add(widget.NewLabel(body))

	log.Printf("Finished displaying article")
	//label2 := widget.NewLabel(fmt.Sprintf("[%v]\n[%v]\n[%v]\n[%v]\n]", a, messageId, body, err))

	//content.Add(label2)

	ScrollReset()
	content.Refresh()
}

//var navContainer, contentContainer *fyne.CanvasObject

var ScrollReset func()

func main() {
	flag.Parse()

	kc = kothawoc.NewClient(*path, *port)
	//kc = kothawoc.NewClient(os.Getenv("PWD") + "/../kothawoc/data")

	kc.Dial()

	lang.AddTranslationsFS(localeFS, "locales")

	myApp := app.New()
	myWindow := myApp.NewWindow(lang.L("App: NNTP/TOR"))
	mainWindow = myWindow
	content = container.NewVBox()

	// Create a list for the navigation pane

	navList := container.NewVBox()
	for _, i := range []func() fyne.CanvasObject{
		func() fyne.CanvasObject {
			vbox := container.NewVBox()

			button := widget.NewButton("Home", nil)
			button.Importance = widget.LowImportance
			button.Alignment = widget.ButtonAlignLeading

			button.OnTapped = func() {
				displayHome(content)
			}

			vbox.Add(button)

			return vbox
		},
		func() fyne.CanvasObject {
			vbox := container.NewVBox()

			button := widget.NewButton(lang.L("List Newsgroups"), nil)
			button.Importance = widget.LowImportance
			button.Alignment = widget.ButtonAlignLeading
			n := 0
			button.OnTapped = func() {
				displayNewsgroupList(content)
				return
				n++
				button := widget.NewButton(fmt.Sprintf("Add [%d]", n), nil)
				button.Importance = widget.LowImportance
				button.Alignment = widget.ButtonAlignLeading

				content.Objects = append([]fyne.CanvasObject{button}, content.Objects...)

			}

			vbox.Add(button)

			return vbox
		},
		func() fyne.CanvasObject {
			// peers entry
			vbox := container.NewVBox()
			leftbox = vbox
			createPeerList(vbox, content, myWindow)
			return vbox
		},
		func() fyne.CanvasObject {
			//	label := widget.NewLabel("Copy")
			label := NewTappableLabel("New Post")
			vbox := container.NewVBox()
			vbox.Add(label)
			label.OnTapped = func(e *fyne.PointEvent) {
				newPost(content)

			}
			label.OnTappedSecondary = func(e *fyne.PointEvent) {
				log.Printf("This is the override of ontapped.")
				menuItem1 := fyne.NewMenuItem("A", nil)
				menuItem2 := fyne.NewMenuItem("B", nil)
				menuItem3 := fyne.NewMenuItem("C", nil)
				menu := fyne.NewMenu("File", menuItem1, menuItem2, menuItem3)

				popUpMenu := widget.NewPopUpMenu(menu, myWindow.Canvas())

				popUpMenu.ShowAtPosition(e.AbsolutePosition)
				popUpMenu.Show()
			}
			//	label.SetText("New Post")
			return vbox
		},
		func() fyne.CanvasObject {
			//	label := widget.NewLabel("Copy")
			label := NewTappableLabel("Add Group")
			vbox := container.NewVBox()
			vbox.Add(label)
			label.OnTapped = func(e *fyne.PointEvent) {
				//newPost(content)
				displayAddGroup(content)
			}

			return vbox
		},
	} {
		navList.Add(i())
	}

	//displayHome(content)
	displayWelcome()

	navContainer := container.NewScroll(navList)
	contentContainer := container.NewScroll(content)
	ScrollReset = func() {
		contentContainer.Offset.X = 0
		contentContainer.Offset.Y = 0
	}

	split := container.NewHSplit(navContainer, contentContainer)

	contentContainer.Offset.X = 0
	contentContainer.Offset.Y = 0

	split.Offset = 0.2 // Adjust the initial size of the left pane

	// Set up the main content with the split container
	myWindow.SetContent(split)

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
