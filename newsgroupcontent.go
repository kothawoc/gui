package main

import (
	"fmt"
	"log"
	"net/textproto"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"github.com/emersion/go-vcard"

	//vcard "github.com/emersion/go-vcard"

	"github.com/kothawoc/go-nntp"
	"github.com/kothawoc/kothawoc/pkg/messages"
)

func newAddGroup() *fyne.Container {
	content := container.NewVBox()

	dId := kc.DeviceId()
	dln := len(dId)
	shortId := dId[0:3] + "..." + dId[dln-3:]
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

	checkGroup := widget.NewCheckGroup([]string{labelRead, labelReply, labelPost}, func(s []string) { fmt.Println("selected", s) })
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

		msg, err := messages.CreateNewsGroupMail(kc.DeviceKey(),
			kc.Server.IdGenerator, dId+"."+groupName.Text, groupDescription.Text, card, nntp.PostingPermitted)
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

	form.Append("Name "+shortId+".", groupName)
	form.Append("Description", groupDescription)
	form.Append("ID Alias", idAlias)
	form.Append("Language", language)
	form.Append("URL", url)
	form.Append("Group Permissions", checkGroup)
	form.Append("Extra Perms", pm.Vbox)

	label := widget.NewLabel("Select an item from the navigation pane")
	content.Add(label)
	content.Add(form)

	content.Refresh()

	return content
}

func newNewsgroupList() *container.Scroll {
	content := container.NewVBox()

	label := widget.NewLabel("Select an item from the news groups list")
	label.Wrapping = fyne.TextWrapBreak
	content.Add(label)

	groups, err := kc.NNTPclient.List("")
	if err != nil {
		return nil
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

	content.Resize(parent.Trailing.Size())
	content.Show()
	content.Refresh()

	setContent(content)

	scroll := container.NewScroll(content)

	//log.Printf("display mds[%s]", mds)

	//	rt := widget.NewRichTextFromMarkdown(mds)

	//	content.Add(rt)
	return scroll
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
	//tmp, err := kc.NNTPclient.ListOverviewFmt()
	//log.Printf("listoverviewformat [%#v][%v]", tmp, err)
	res, err := kc.NNTPclient.Over(int(a.Low), int(a.High))
	for _, line := range res {

		//label := widget.NewLabel(line.Subject)
		/*
			** TODO FIXME:
			** This is broken, //tmp, err := kc.NNTPclient.ListOverviewFmt() Hangs, and it gets this in the wrong order, check if it's the server or client (or both).
			2024/09/02 18:30:52 AND THE LINES WAS[nntpclient.OverItem{Number:"1725134320", Subject:"AddPeer 3rm3lavawfdngj6tspw2rrsfjcz4pxh3o7ltjxaugyhnauhir7ngvrad", Date:"3rm3lavawfdngj6tspw2rrsfjcz4pxh3o7ltjxaugyhnauhir7ngvrad", MessageId:"Sat, 31 Aug 2024 19:59:07 +0000", References:"<1jd6tgb-snjkpy2ltz5s3s2j3rxo6igmp2xxmvkj@3rm3lavawfdngj6tspw2rrsfjcz4pxh3o7ltjxaugyhnauhir7ngvrad>", bytesMetadata:"", linesMetadata:"368"}]
			2024/09/02 18:30:52 AND THE LINES WAS[nntpclient.OverItem{Number:"1725136214", Subject:"AddPeer 3buuqev6fbwybjo6qch2bescuelkcqm4sf7w73dm3vmf55qnudug2kyd", Date:"3rm3lavawfdngj6tspw2rrsfjcz4pxh3o7ltjxaugyhnauhir7ngvrad", MessageId:"Sat, 31 Aug 2024 20:30:17 +0000", References:"<1jd6vap-ihwa3bloqlwdr2d5iabtvhg5egabm7j5@3rm3lavawfdngj6tspw2rrsfjcz4pxh3o7ltjxaugyhnauhir7ngvrad>", bytesMetadata:"", linesMetadata:"368"}]
			From:=Date
			Date:=MessageId
			MessageId:=References
		*/
		MessageId := line.MessageId

		item := NewTappableLabel(line.Subject)
		item.OnTapped = func(e *fyne.PointEvent) {
			displayMessage(content, line.Number)
		}
		log.Printf("AND THE LINES WAS[%#v]", line)

		item.OnTappedSecondary = func(e *fyne.PointEvent) {

			menuItem1 := fyne.NewMenuItem(lang.L("Cancel Article"), nil)
			menu := fyne.NewMenu(lang.L("menu.Peer"), menuItem1)
			menuItem1.Action = func() {
				msg := (&messages.MessageTool{
					Article: &nntp.Article{

						Header: textproto.MIMEHeader{
							"Subject":                   {"cmsg cancel " + MessageId},
							"Control":                   {"cancel " + MessageId},
							"Newsgroups":                {group},
							"Content-Type":              {"multipart/mixed; boundary=\"nxtprt\""},
							"Content-Transfer-Encoding": {"8bit"},
						},
					},
					Preamble: "This is a MIME control message.",
					Parts: []messages.MimePart{
						{
							Header:  textproto.MIMEHeader{"Content-Type": []string{"application/news-groupinfo;charset=UTF-8"}},
							Content: []byte("Cancel " + MessageId),
						},
						{
							Header:  textproto.MIMEHeader{"Content-Type": []string{"text/plain;charset=UTF-8"}},
							Content: []byte("This is a system control message to delete the article " + MessageId),
						},
					},
				}).RawMail()

				kc.NNTPclient.Post(strings.NewReader(msg))
				log.Printf("AND THE LINES WAS[%#v]", line)

				createPeerList(leftbox)
			}

			//	m := messages.MessageTool{}
			//	m.
			popUpMenu := widget.NewPopUpMenu(menu, mainWindow.Canvas())

			popUpMenu.ShowAtPosition(e.AbsolutePosition)
			popUpMenu.Show()
		}

		//text := widget.NewRichTextWithText(fmt.Sprintf("%v", res))
		content.Add(item)
	}

	content.Refresh()

}
