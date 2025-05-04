package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BurntSushi/toml"
	"github.com/ladecadence/EkiTelemetry/pkg/serialport"
)

type ConfigData struct {
	PortName    string
	LogFile     string
	ImageFolder string
	Server      string
	User        string
	Password    string
}

type Config struct {
	Content *fyne.Container
	Data    ConfigData
}

func CreateConfig(w *fyne.Window) *Config {
	config := Config{}

	// create GUI
	ports := serialport.GetSerialPorts()
	labelSerial := widget.NewLabel("Serial port:")
	comboSerial := widget.NewSelect(ports, func(value string) {
		config.Data.PortName = value
	})
	labelLog := widget.NewLabel("Log file:")
	inputLog := widget.NewEntry()

	buttonLog := widget.NewButtonWithIcon("", theme.FileIcon(), func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer != nil {
				fmt.Println(writer.URI().Path())
				config.Data.LogFile = writer.URI().Path()
				inputLog.SetText(writer.URI().Path())
			}
		}, *w)

	})
	groupLog := container.NewBorder(nil, nil, nil, buttonLog, inputLog)

	labelImages := widget.NewLabel("Images folder:")
	inputImages := widget.NewEntry()

	buttonImages := widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				fmt.Println(uri.Path())
				config.Data.ImageFolder = uri.Path()
				inputImages.SetText(uri.Path())
			}
		}, *w)

	})
	groupImages := container.NewBorder(nil, nil, nil, buttonImages, inputImages)

	labelServer := widget.NewLabel("Server:")
	labelUser := widget.NewLabel("User:")
	labelPassword := widget.NewLabel("Password:")
	inputServer := widget.NewEntry()
	inputServer.OnChanged = func(s string) { config.Data.Server = s }
	inputUser := widget.NewEntry()
	inputUser.OnChanged = func(s string) { config.Data.User = s }
	inputPassword := widget.NewPasswordEntry()
	inputPassword.OnChanged = func(s string) { config.Data.Password = s }

	buttonSave := widget.NewButton("Save", func() {
		config.SaveConfig()
	})

	padded := container.New(layout.NewVBoxLayout(),
		labelSerial,
		comboSerial,
		labelLog,
		groupLog,
		labelImages,
		groupImages,
		labelServer,
		inputServer,
		labelUser,
		inputUser,
		labelPassword,
		inputPassword,
		buttonSave,
	)
	content := container.NewPadded(padded)

	config.Content = content

	// try to load config
	err := config.LoadConfig()
	if err != nil {
		// do nothing
	}
	// ok, fill values
	if slices.Contains(ports, config.Data.PortName) {
		comboSerial.SetSelected(config.Data.PortName)
	}
	inputLog.SetText(config.Data.LogFile)
	inputImages.SetText(config.Data.ImageFolder)
	inputServer.SetText(config.Data.Server)
	inputUser.SetText(config.Data.User)
	inputPassword.SetText(config.Data.Password)

	return &config
}

func (c *Config) SaveConfig() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(dir, "EkiTelemetry")
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err
	}

	data, err := toml.Marshal(c.Data)
	if err != nil {
		return err
	}
	confFile := filepath.Join(configPath, "config.toml")
	if err := os.WriteFile(confFile, data, 0666); err != nil {
		return err
	}
	return nil
}

func (c *Config) LoadConfig() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(dir, "EkiTelemetry", "config.toml")
	data := ConfigData{}
	_, err = toml.DecodeFile(configPath, &data)
	if err != nil {
		return err
	}
	c.Data = data

	return nil
}
