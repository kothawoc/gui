package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type PermsT struct {
	Idx                                  int
	Vbox, hbox                           *fyne.Container
	TorId                                *widget.Entry
	Read, Reply, Post, Cancel, Supersede *widget.Check

	Remove *widget.Button
}
type PermsMgr struct {
	Form  *widget.Form
	Vbox  *fyne.Container
	Items []PermsT
}

func NewPemsMgr(f *widget.Form) PermsMgr {
	p := PermsMgr{
		Form:  f,
		Vbox:  container.NewVBox(),
		Items: make([]PermsT, 0),
	}
	return p
}

func (p *PermsMgr) Render() {
	p.Vbox.RemoveAll()
	for idx, item := range p.Items {
		item.Idx = idx
		item.Remove.OnTapped = func() {
			p.Items = append(p.Items[:idx], p.Items[idx+1:]...)
			p.Render()
		}
		p.Vbox.Add(widget.NewSeparator())
		p.Vbox.Add(item.Vbox)
	}

	p.Vbox.Add(widget.NewButton("+", func() {
		p.Add()
	}))

	p.Vbox.Refresh()
	p.Form.Refresh()
	if content != nil {
		content.Refresh()
	}
}

func NewCheck(s string, checked bool) *widget.Check {
	cb := widget.NewCheck(s, func(i bool) {})
	if checked {
		cb.Checked = true
	}
	return cb
}

func (p *PermsMgr) Add() {

	//{"post", "read", "reply", "cancel", "supersede"}
	labelSupersede := lang.L("Supersede")
	labelCancel := lang.L("Cancel")
	labelRead := lang.L("Read")
	labelReply := lang.L("Reply")
	labelPost := lang.L("Post")

	newItem := PermsT{
		Vbox:      container.NewVBox(),
		hbox:      container.NewHBox(),
		TorId:     widget.NewEntry(),
		Read:      NewCheck(labelRead, true),
		Reply:     NewCheck(labelReply, true),
		Post:      NewCheck(labelPost, false),
		Cancel:    NewCheck(labelCancel, false),
		Supersede: NewCheck(labelSupersede, false),
		Remove:    widget.NewButton("-", func() {}),
	}
	newItem.hbox.Add(newItem.Read)
	newItem.hbox.Add(newItem.Reply)
	newItem.hbox.Add(newItem.Post)
	newItem.hbox.Add(newItem.Cancel)
	newItem.hbox.Add(newItem.Supersede)
	newItem.hbox.Add(widget.NewLabel("    "))
	newItem.hbox.Add(newItem.Remove)

	newItem.Vbox.Add(newItem.TorId)
	newItem.Vbox.Add(newItem.hbox)

	p.Items = append(p.Items, newItem)
	p.Render()
}
