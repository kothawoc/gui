package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/mail"
	"net/textproto"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/kothawoc/go-nntp"
	"github.com/kothawoc/kothawoc/pkg/messages"
)

type ArtItem struct {
	Date      int
	MessageId string
	Article   *nntp.Article
}

func createMessageCard(messageNumber string) *fyne.Container {

	if messageNumber == "" {
		return container.NewWithoutLayout()
	}
	mesgCard := container.NewVBox()
	//label := widget.NewLabel("The content of: " + messageNumber)
	//mesgCard.Add(label)

	//a, messageId, rawArticle, err := kc.NNTPclient.Article(messageNumber)
	_, _, rawArticle, err := kc.NNTPclient.Article(messageNumber)

	/*
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
	*/
	slog.Info("Crash City", "messageNum", messageNumber, "rawArt", rawArticle)
	mailMsg, err := mail.ReadMessage(rawArticle)
	if err != nil {
		log.Fatal(err)
	}

	article := &nntp.Article{
		Header: textproto.MIMEHeader(mailMsg.Header),
		Body:   mailMsg.Body,
		//	Bytes:  len([]byte(body)),
		//	Lines:  strings.Count(body, "\n"),
	}

	msg := messages.NewMessageToolFromArticle(article)

	// content := container.NewBorder(
	hbox := container.NewHBox()
	shortFrom := msg.Article.Header.Get("From")
	shortFrom = shortFrom[:3] + "..." + shortFrom[len(shortFrom)-4:]
	tl := widget.NewRichTextWithText(shortFrom)
	tc := widget.NewRichTextFromMarkdown("**" + msg.Article.Header.Get("Subject") + "**")
	tr := widget.NewRichTextFromMarkdown(msg.Article.Header.Get("Date"))
	hbox.Add(tl)
	hbox.Add(tc)
	hbox.Add(tr)
	bbox := container.NewBorder(nil, nil, tl, tr, tc)

	mesgCard.Add(bbox)

	//content.Add(widget.NewLabel("From: " + msg.Article.Header.Get("From") + "\n" +
	//	"Subject: " + msg.Article.Header.Get("Subject") + "\n" +
	//	"Date: " + msg.Article.Header.Get("Date") + "\n" +
	//	"\n"))
	//content.Add(widget.NewLabel("Subject: " + msg.Article.Header.Get("Subject")))
	//content.Add(widget.NewLabel("Date: " + msg.Article.Header.Get("Date")))
	//content.Add(widget.NewLabel(" "))
	preamble := widget.NewLabel(msg.Preamble)
	preamble.Wrapping = fyne.TextWrap(fyne.TextWrapWord)
	mimeParts := widget.NewLabel(fmt.Sprintf("%s", msg.Parts))
	mimeParts.Wrapping = fyne.TextWrap(fyne.TextWrapWord)
	mesgCard.Add(preamble)
	mesgCard.Add(mimeParts)

	log.Printf("Finished displaying article")
	//label2 := widget.NewLabel(fmt.Sprintf("[%v]\n[%v]\n[%v]\n[%v]\n]", a, messageId, body, err))

	//content.Add(label2)

	return mesgCard
}

func displayMainFeed() {
	ArtDB := map[string]ArtItem{}
	ArtDBindex := []string{}
	//label := widget.NewLabel(lang.L("Welcome: Create a new account"))

	// nasty memory hungry hack, we need to load them into an sqlite cache for manipulating them
	groups, err := kc.NNTPclient.List("")
	if err != nil {
		slog.Warn("Can't list groups")
	}
	for _, group := range groups {

		a, err := kc.NNTPclient.Group(group.Name)

		if err != nil {
			log.Printf("Error in listing newsgroup conftent: [%v]", err)
			return
		}
		//tmp, err := kc.NNTPclient.ListOverviewFmt()
		//log.Printf("listoverviewformat [%#v][%v]", tmp, err)
		res, err := kc.NNTPclient.Over(int(a.Low), int(a.High))
		for _, art := range res {
			//date := time.art.Date

			//s := "Tue Sep 16 21:58:58 +0000 2014"
			// Date: Tue, 10 Sep 2024 15:11:19 +0000
			const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700 "
			//  MessageId:"Sat, 31 Aug 2024 19:59:07 +0000", References:"<1jd6tgb-snjkpy2l
			slog.Info(fmt.Sprintf("[%#v][%v]\n", art))
			// actually the data TODO FIX THIS BUG
			t, err := time.Parse(rfc2822, art.Date)
			if err != nil {
				slog.Info("Error", "error", err)
				return
			}

			// actually the message id, TO BUG FIX THIS
			idx := fmt.Sprintf("%d-%s", t, art.MessageId)
			ArtDB[idx] = ArtItem{
				MessageId: art.MessageId,
				Date:      int(t.Unix()),
			}
			ArtDBindex = append(ArtDBindex, idx)

			//fmt.Println(u)
			//f := t.Format(time.UnixDate)
			//fmt.Println(f)
		}
	}

	sort.Strings(ArtDBindex)

	Vbox := container.NewVBox()

	for i := len(ArtDBindex) - 1; i >= 0; i-- {
		//ArtDBindex = append(ArtDBindex, i)
		Vbox.Add(widget.NewSeparator())
		Vbox.Add(createMessageCard(ArtDB[ArtDBindex[i]].MessageId))
	}

	/*
		rt := widget.NewRichTextFromMarkdown(lang.L("App: Welcome Message"))
		rt.Wrapping = fyne.TextWrapWord

		dId := kc.DeviceId()
		//dln := len(dId)
		//shortId := dId[0:3] + "..." + dId[dln-3:]
		groupName := widget.NewEntry()
		groupDescription := widget.NewEntry()
		idAlias := widget.NewEntry()
		language := widget.NewEntry()
		url := widget.NewEntry()

		//{"post", "read", "reply", "cancel", "supersede"}
		//labelSupersede := lang.L("Supersede")
		//labelCancel := lang.L("Cancel")
		labelRead := lang.L("Read")
		labelReply := lang.L("Reply")
		labelPost := lang.L("Post")

		//	labelUserId := lang.L("UserId")

		checkGroup := widget.NewCheckGroup([]string{}, func(s []string) { fmt.Println("selected", s) })
		checkGroup.Selected = []string{labelRead, labelReply}
		checkGroup.Horizontal = true

		sendFunc := func() {}

		form := &widget.Form{
			//	Items: []*widget.FormItem{
			//	},
			OnCancel: func() {
				fmt.Println("Cancelled")
			},
			OnSubmit: func() {
				sendFunc()

			},
		}
		pm := NewPemsMgr(form)
		//pm.Add()
		pm.Render()

		sendFunc = func() {

			card := vcard.Card{}
			card.SetValue(vcard.FieldNickname, idAlias.Text)
			card.SetValue(vcard.FieldLanguage, language.Text)
			card.SetValue(vcard.FieldURL, url.Text)
			gParms := vcard.Params{}
			for _, item := range checkGroup.Selected {
				if item == labelRead {
					gParms["read"] = []string{"true"}
				}
				if item == labelReply {
					gParms["reply"] = []string{"true"}
				}
				if item == labelPost {
					gParms["post"] = []string{"true"}
				}
			}
			card.Add("X-KW-PERMS", &vcard.Field{
				Value:  "group",
				Params: gParms,
			})
			for _, item := range pm.Items {
				parms := vcard.Params{}
				if item.Read.Checked {
					parms["read"] = []string{"true"}
				}
				if item.Reply.Checked {
					parms["reply"] = []string{"true"}
				}
				if item.Post.Checked {
					parms["post"] = []string{"true"}
				}
				if item.Cancel.Checked {
					parms["cancel"] = []string{"true"}
				}
				if item.Supersede.Checked {
					parms["supersede"] = []string{"true"}
				}
				card.Add("X-KW-PERMS", &vcard.Field{
					Value:  item.TorId.Text,
					Params: parms,
				})
			}
			vcard.ToV4(card)

			msg, err := messages.CreateNewsGroupMail(kc.DeviceKey(), kc.Server.IdGenerator, dId+"."+groupName.Text, groupDescription.Text, card, nntp.PostingPermitted) //(string, error)
			if err != nil {
				log.Printf("Failed at creating new groups mail.")
			}

			kc.NNTPclient.Post(strings.NewReader(msg))

		}
		// add group name
		// group descripton
		// group vcard
		//   CATEGORIES:public\,cabbage
		//   LANG:enCheckGroup Item 2
		//   NICKNAME:The magic bus
		//   URL:https://github.com/kothawoc

		form.Append("ID Alias", idAlias)
		form.Append("Description", groupDescription)
		form.Append("Language", language)
		form.Append("URL", url)
		form.Append("Group Permissions", checkGroup)
		form.Append("Extra Perms", pm.Vbox)
	*/
	content.RemoveAll()
	//	label := widget.NewLabel("Select an item from the navigation pane")
	//content.Add(label)
	//content.Add(form)
	content.Add(Vbox)

	content.Refresh()

}
