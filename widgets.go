package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type tappableLabel struct {
	widget.Label
	OnTapped          func(e *fyne.PointEvent)
	OnTappedSecondary func(e *fyne.PointEvent)
}

func NewTappableLabel(text string) *tappableLabel {
	label := &tappableLabel{}
	label.ExtendBaseWidget(label)
	label.SetText(text)
	label.OnTapped = func(e *fyne.PointEvent) {}
	label.OnTappedSecondary = func(e *fyne.PointEvent) {}

	return label
}

func (t *tappableLabel) Tapped(e *fyne.PointEvent) {
	t.OnTapped(e)
}

func (t *tappableLabel) TappedSecondary(e *fyne.PointEvent) {
	t.OnTappedSecondary(e)
}
