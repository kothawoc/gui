package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/mail"
	"net/textproto"
	"strings"

	"fyne.io/fyne/storage"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"

	"github.com/kothawoc/go-nntp"
	"github.com/kothawoc/kothawoc/pkg/messages"
)

func displayHome(content *fyne.Container) {

	slog.Info("test one")
	u := storage.NewURI("testfile.md")
	// a, err := storage.OpenFileFromURI(u)
	a := "aa"
	err := errors.New("poop")
	slog.Info("test moar")
	//	pathstore := storage. .Storage().RootURI().Path() + "file.txt"
	slog.Info("test storage stuff", "u", u, "a", a, "err", err)

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

func newReply(group, subject, references string) *widget.Form {
	editor := widget.NewMultiLineEntry()
	editor.SetMinRowsVisible(8)
	if strings.ToLower(subject[0:3]) != "re:" {
		subject = "Re: " + subject
	}

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
			//	fyne.CurrentApp().SendNotification(&fyne.Notification{
			//		Title:   "Form for: " + subject,
			//		Content: editor.Text,
			//	})
			mail := messages.NewMessageTool()
			mail.Article.Header.Add("Newsgroups", group)
			mail.Article.Header.Add("Subject", subject)
			mail.Article.Header.Add("References", references)
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
	form.Append("Message", editor)
	//form.Hidden = false
	//	form.Refresh()

	return form

}

var updateGroupsList func() = func() {}

func newPost() *widget.Form {

	groups, err := kc.NNTPclient.List("")
	if err != nil {
		return nil
	}
	//	mds := ""
	groupList := []string{}
	groupsEntry := widget.NewSelectEntry(groupList)
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
	//form.SetOnValidationChanged()

	updateGroupsList = func() {
		groupList := []string{}
		for _, g := range groups {
			groupList = append(groupList, g.Name)
		}
		groupsEntry.SetOptions(groupList)
		slog.Info("updating groups list")
	}
	return form

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
