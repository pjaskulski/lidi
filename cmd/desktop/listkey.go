package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// list widget
type keyList struct {
	*widget.List
}

func (l *keyList) onUp() {
	fmt.Println("Up")
}

func (l *keyList) onDown() {
	fmt.Println("Down")
}

func newKeyList(data binding.ExternalStringList) *keyList {

	list := &keyList{widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(2,
				widget.NewLabel("template"),
			)
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			item := i.(binding.String)
			textLabel := o.(*fyne.Container).Objects[0].(*widget.Label)
			textLabel.Bind(item)
		})}

	list.ExtendBaseWidget(list)
	return list
}

func (l *keyList) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		l.onUp()
	case fyne.KeyReturn:
		l.onDown()
	}
}

func (l *keyList) SelectNew() {
	l.Unselect(0)
	l.Select(0)
}
