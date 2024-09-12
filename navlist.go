package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

func newNavlist() *fyne.Container {

	// Create a list for the navigation pane
	navList := container.NewVBox()
	go func() {
		for _, i := range []func() fyne.CanvasObject{
			func() fyne.CanvasObject {
				vbox := container.NewVBox()

				button := widget.NewButton("Home", nil)
				button.Importance = widget.LowImportance
				button.Alignment = widget.ButtonAlignLeading

				button.OnTapped = func() {
					setContent(mainFeed)
				}

				vbox.Add(button)

				return vbox
			},
			func() fyne.CanvasObject {
				vbox := container.NewVBox()

				button := widget.NewButton(lang.L("List Newsgroups"), nil)
				button.Importance = widget.LowImportance
				button.Alignment = widget.ButtonAlignLeading

				button.OnTapped = func() {
					setContent(newNewsgroupList())
					return

				}

				vbox.Add(button)

				return vbox
			},
			func() fyne.CanvasObject {
				// peers entry
				vbox := container.NewVBox()
				leftbox = vbox
				createPeerList(vbox)
				return vbox
			},
			func() fyne.CanvasObject {
				//	label := widget.NewLabel("Copy")
				label := NewTappableLabel("New Post")
				vbox := container.NewVBox()
				vbox.Add(label)
				label.OnTapped = func(e *fyne.PointEvent) {
					//	updateGroupsList()
					setContent(postForm)
				}
				label.OnTappedSecondary = func(e *fyne.PointEvent) {
					log.Printf("This is the override of ontapped.")
					menuItem1 := fyne.NewMenuItem("A", nil)
					menuItem2 := fyne.NewMenuItem("B", nil)
					menuItem3 := fyne.NewMenuItem("C", nil)
					menu := fyne.NewMenu("File", menuItem1, menuItem2, menuItem3)

					popUpMenu := widget.NewPopUpMenu(menu, mainWindow.Canvas())

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
					setContent(addGroup)
				}

				return vbox
			},
		} {
			navList.Add(i())
		}
	}()

	return navList
}
