package main

import (
	"embed"
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

/*
	func main() {
		a := app.New()
		w := a.NewWindow("Hello")

		hello := widget.NewLabel("Hello Fyne!")
		button := widget.NewButton("Hi!", nil)
		newContainer := container.NewVBox(
			hello,
			button,
		)

		button.OnTapped = func() {
			hello.SetText("Welcome :)")
			if hello.Visible() {
				hello.Hide()
			} else {
				hello.Show()
			}

		}

		newBoarders := container.NewBorder(nil, nil, hello, nil, newContainer)

		w.SetContent(newBoarders)

		w.ShowAndRun()
	}
*/
var peerlist []string = []string{}

func createPeerList(vbox, content *fyne.Container, win fyne.Window) {
	//*
	group := kc.DeviceId() + ".peers"

	a, err := kc.NNTPclient.Group(group)
	peerlist := []string{}

	if err != nil {
		log.Printf("Error in listing newsgroup conftent: [%v]", err)
		return
	}
	res, err := kc.NNTPclient.Over(int(a.Low), int(a.High))
	for _, msg := range res {

		peer := strings.Split(msg.Subject, " ")[1]
		peerlist = append(peerlist, peer)
		//label := widget.NewLabel(line.Subject)

		//	item := NewTappableLabel(msg.Subject)
		//	item.OnTapped = func(e *fyne.PointEvent) {
		//		displayMessage(content, msg.Number)
		//	}

		//text := widget.NewRichTextWithText(fmt.Sprintf("%v", res))
		//	content.Add(item)
	}
	//*/
	vbox.RemoveAll()
	button := widget.NewButton(lang.L("Peers"), nil)
	vbox.Add(button)
	edit := widget.NewEntry()
	edit.OnSubmitted = func(text string) {
		edit.Hidden = true
		err := kc.AddPeer(text)
		if err != nil {
			log.Printf("GUI ERROR: failed to add peer.")
			button.Show()
			return
		}

		peerlist = append(peerlist, text)
		createPeerList(vbox, content, win)
		vbox.Refresh()
		//additem(text, vbox, content, edit, button)
		//vbox.Add(item)
	}

	edit.Hide()
	edit.SetText("")
	vbox.Add(edit)

	//			contacts := []string
	button.OnTapped = func() {
		button.Hide()
		edit.SetText("")
		edit.Show()
		win.Canvas().Focus(edit)
	}

	for _, text := range peerlist {

		//item := widget.NewButton(text, nil)

		//	label := widget.NewLabel("Copy")
		item := NewTappableLabel(text)
		//vbox := container.NewVBox()
		//vbox.Add(label)

		//	peerlist = append(peerlist, text)
		//	item.Importance = widget.LowImportance
		//	item.Alignment = widget.ButtonAlignLeading
		vbox.Add(item)

		item.OnTapped = func(e *fyne.PointEvent) {
			content.RemoveAll()
			label := widget.NewLabel(text)
			content.Add(label)

			button := widget.NewButton(lang.L("Delete"), nil)
			button.OnTapped = func() {

				// this is now broken, because it uses the nntp server for the peer list
				// so it needs to send a cancel message to delete it.
				for i := 0; i < len(peerlist); i++ {
					if peerlist[i] == text {
						peerlist = append(peerlist[:i], peerlist[i+1:]...)
					}
				}
				createPeerList(vbox, content, win)
			}

			content.Add(button)
			content.Refresh()
		}

		item.OnTappedSecondary = func(e *fyne.PointEvent) {

			menuItem1 := fyne.NewMenuItem(lang.L("Unpeer"), nil)
			menuItem2 := fyne.NewMenuItem(lang.L("Subcscribe to Groups"), nil)
			menu := fyne.NewMenu(lang.L("menu.Peer"), menuItem1, menuItem2)

			popUpMenu := widget.NewPopUpMenu(menu, win.Canvas())

			popUpMenu.ShowAtPosition(e.AbsolutePosition)
			popUpMenu.Show()
		}

	}

	vbox.Refresh()
}

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

	//	groupsLabel := widget.NewLabel(lang.L("News Groups"))
	groupsEntry := widget.NewEntry()
	//	groupsBox := container.NewPadded(groupsLabel, groupsEntry)
	//	groupsBox.Add(groupsLabel)
	//	groupsBox.Add(groupsEntry)

	//	groupsLabel.Resize(fyne.NewSize(100, 0))
	//	groupsEntry.Resize(fyne.NewSize(content.Size().Width-groupsLabel.Size().Width, 0))
	//groupsBox.Resize(fyne.NewSize(500, 0))

	//	subjectLabel := widget.NewLabel(lang.L("Subject"))
	subjectEntry := widget.NewEntry()
	//subjectEntry.
	//	subjectBox := container.NewHBox()
	//	subjectBox.Add(subjectLabel)
	//	subjectBox.Add(subjectEntry)

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

	//content.Set
	//	content.Add(groupsBox)
	//	content.Add(subjectBox)
	//	content.Add(editor)
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

func displayNewsgroupContent(content *fyne.Container, group string) {
	content.RemoveAll()
	label := widget.NewLabel("Select an article to read: " + group)
	content.Add(label)

	a, err := kc.NNTPclient.Group(group)

	if err != nil {
		log.Printf("Error in listing newsgroup conftent: [%v]", err)
		return
	}
	res, err := kc.NNTPclient.Over(int(a.Low), int(a.High))
	for _, line := range res {

		//label := widget.NewLabel(line.Subject)

		item := NewTappableLabel(line.Subject)
		item.OnTapped = func(e *fyne.PointEvent) {
			displayMessage(content, line.Number)
		}

		//text := widget.NewRichTextWithText(fmt.Sprintf("%v", res))
		content.Add(item)
	}

	content.Refresh()

}
func displayNewsgroupList(content *fyne.Container) {
	content.RemoveAll()
	label := widget.NewLabel("Select an item from the news groups list")
	content.Add(label)

	groups, err := kc.NNTPclient.List("")
	if err != nil {
		return
	}
	//	mds := ""
	for _, g := range groups {

		item := NewTappableLabel(g.Name)
		item.OnTapped = func(e *fyne.PointEvent) {
			log.Printf("Clicked on group: [%s]", g.Name)
			displayNewsgroupContent(content, g.Name)
		}
		content.Add(item)

		log.Printf("display groups[%#v]", g)
		//		mds += g.Name + "\n\n"
	}

	content.Show()
	content.Refresh()

	//log.Printf("display mds[%s]", mds)

	//	rt := widget.NewRichTextFromMarkdown(mds)

	//	content.Add(rt)

}

//var navContainer, contentContainer *fyne.CanvasObject

var ScrollReset func()

func main() {
	kc = kothawoc.NewClient(os.Getenv("PWD") + "/../kothawoc/data")

	kc.Dial()

	lang.AddTranslationsFS(localeFS, "locales")

	myApp := app.New()
	myWindow := myApp.NewWindow(lang.L("App: NNTP/TOR"))
	content := container.NewVBox()

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
	} {
		navList.Add(i())
	}

	// Create a content area
	//label := widget.NewLabel("Select an item from the navigation pane")
	//content.Add(label)

	displayHome(content)

	// Create the split container
	navContainer := container.NewScroll(navList)
	contentContainer := container.NewScroll(content)
	ScrollReset = func() {
		contentContainer.Offset.X = 0
		contentContainer.Offset.Y = 0
	}

	split := container.NewHSplit(navContainer, contentContainer)

	//contentContainer.Offset.X = 0
	//contentContainer.Offset.Y = 0

	split.Offset = 0.2 // Adjust the initial size of the left pane

	// Set up the main content with the split container
	myWindow.SetContent(split)

	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
