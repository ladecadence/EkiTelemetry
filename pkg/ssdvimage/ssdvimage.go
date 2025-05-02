package ssdvimage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SSDVImage struct {
	Content *fyne.Container
	Border  *fyne.Container
	Image   *canvas.Image
	Label   *widget.Label
}

func CreateSSDVImage() *SSDVImage {
	ssdvImage := SSDVImage{}
	ssdvImage.Image = canvas.NewImageFromFile("")
	ssdvImage.Label = widget.NewLabel("Waiting for image...")
	ssdvImage.Border = container.NewBorder(nil, ssdvImage.Label, nil, nil, ssdvImage.Image)
	ssdvImage.Content = container.NewPadded(ssdvImage.Border)
	return &ssdvImage
}

func (s *SSDVImage) UpdateImage(path string, info string) {
	s.Image = canvas.NewImageFromFile(path)
	s.Image.FillMode = canvas.ImageFillContain
	s.Border.Objects[0] = s.Image
	s.Label.SetText(info)
}
