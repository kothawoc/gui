package main

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"

	"github.com/emersion/go-vcard"
	"github.com/kothawoc/go-nntp"
	"github.com/kothawoc/kothawoc/pkg/messages"
)

func displayWelcome() {

	slog.Info("Done done dial start app", "uri", fyne.CurrentApp().Storage().RootURI())

	label := widget.NewLabel(lang.L("Welcome: Create a new account"))

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

	content.RemoveAll()
	//	label := widget.NewLabel("Select an item from the navigation pane")
	content.Add(label)
	content.Add(form)

	content.Refresh()

}
