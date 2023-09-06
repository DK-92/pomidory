package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type NumericalEntry struct {
	widget.Entry
}

func NewNumericalEntry() *NumericalEntry {
	entry := &NumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (n *NumericalEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' && len(n.Text) < 3 {
		n.Entry.TypedRune(r)
	}
}

func (n *NumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		n.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		n.Entry.TypedShortcut(shortcut)
	}
}

func (n *NumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
