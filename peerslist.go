package main

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	//vcard "github.com/emersion/go-vcard"
)

var peerlist []string = []string{}

func createPeerList(vbox *fyne.Container) {
	//*
	//group := kc.DeviceId() + ".peers"
	go func() {

		res, err := kc.NNTPclient.List("")

		//slog.Info("groups list", "res", a, "error", err)
		//panic("w00t")
		peerlist := []string{}

		if err != nil {
			log.Printf("Error in listing newsgroup conftent: [%v]", err)
			return
		}
		//	res, err := kc.NNTPclient.Over(int(a.Low), int(a.High))
		for _, item := range res {
			isplit := strings.Split(item.Name, ".")
			if len(isplit) == 3 &&
				isplit[0] == kc.DeviceId() &&
				isplit[1] == "peers" {
				peerlist = append(peerlist, isplit[2])
			}
		}
		//*/
		vbox.RemoveAll()
		button := widget.NewButton(lang.L("Peers"), nil)
		vbox.Add(button)
		edit := widget.NewEntry()
		edit.OnSubmitted = func(text string) {
			edit.Hidden = true
			err := kc.AddPeer(text, "TODO: replace myname")
			if err != nil {
				log.Printf("GUI ERROR: failed to add peer.")
				button.Show()
				return
			}

			peerlist = append(peerlist, text)
			createPeerList(vbox)
			vbox.Refresh()

		}

		edit.Hide()
		edit.SetText("")
		vbox.Add(edit)

		//			contacts := []string
		button.OnTapped = func() {
			button.Hide()
			edit.SetText("")
			edit.Show()
			mainWindow.Canvas().Focus(edit)
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

					createPeerList(vbox)

					content.Add(button)
					content.Refresh()
				}

				item.OnTappedSecondary = func(e *fyne.PointEvent) {

					menuItem1 := fyne.NewMenuItem(lang.L("Unpeer"), nil)
					menuItem2 := fyne.NewMenuItem(lang.L("Subcscribe to Groups"), nil)
					menu := fyne.NewMenu(lang.L("menu.Peer"), menuItem1, menuItem2)

					popUpMenu := widget.NewPopUpMenu(menu, mainWindow.Canvas())

					popUpMenu.ShowAtPosition(e.AbsolutePosition)
					popUpMenu.Show()
				}

			}

			vbox.Refresh()
		}
	}()

}
