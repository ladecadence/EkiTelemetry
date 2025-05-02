package console

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Console struct {
	Scroll  *container.Scroll
	Content *fyne.Container
	Output  *widget.RichText
}

func CreateConsole() *Console {
	//output := widget.NewTextGrid()
	output := widget.NewRichText()
	outputScroll := container.NewScroll(output)
	padded := container.NewBorder(nil, nil, nil, nil, outputScroll)
	content := container.NewPadded(padded)
	return &Console{Content: content, Output: output, Scroll: outputScroll}
}

func (c *Console) Append(text string) {
	c.Output.AppendMarkdown(text)
	c.Scroll.ScrollToBottom()
}
