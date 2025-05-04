package ssdvimage

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ladecadence/EkiTelemetry/pkg/api"
	"github.com/ladecadence/EkiTelemetry/pkg/config"
)

type SSDVImage struct {
	Config       *config.Config
	Content      *fyne.Container
	Border       *fyne.Container
	Image        *canvas.Image
	ImageNum     uint8
	Mission      string
	Label        *widget.Label
	UploadButton *widget.Button
}

func CreateSSDVImage(conf *config.Config) *SSDVImage {
	ssdvImage := SSDVImage{Config: conf}
	ssdvImage.Image = canvas.NewImageFromFile("")
	ssdvImage.Label = widget.NewLabel("Waiting for image...")
	ssdvImage.UploadButton = widget.NewButton("Upload Image", ssdvImage.UploadClicked)
	ssdvImage.UploadButton.Disable()
	ssdvBottomBar := container.NewHBox(ssdvImage.Label, layout.NewSpacer(), ssdvImage.UploadButton)
	ssdvImage.Border = container.NewBorder(nil, ssdvBottomBar, nil, nil, ssdvImage.Image)
	ssdvImage.Content = container.NewPadded(ssdvImage.Border)
	return &ssdvImage
}

func (s *SSDVImage) UpdateImage(number uint8, path string, info string, mission string) {
	s.Mission = mission
	s.ImageNum = number
	s.Image = canvas.NewImageFromFile(path)
	if s.Image != nil {
		s.UploadButton.Enable()
	} else {
		s.UploadButton.Disable()
	}
	s.Image.FillMode = canvas.ImageFillContain
	s.Border.Objects[0] = s.Image
	s.Label.SetText(info)
}

func (s *SSDVImage) UploadClicked() {
	err := s.Upload()
	if err != nil {
		s.Label.SetText(fmt.Sprintf("Problem uploading image: %v", err))
	} else {
		s.Label.SetText("Image uploaded.")
	}
}

func (s *SSDVImage) Upload() error {
	// image
	imgFile := s.Config.Data.ImageFolder + fmt.Sprintf("/ssdv%d.jpg", s.ImageNum)

	// make upload
	api := api.API{Server: s.Config.Data.Server, User: s.Config.Data.User, Password: s.Config.Data.Password}
	err := api.ImageUpload(imgFile, s.Mission)

	return err

}
