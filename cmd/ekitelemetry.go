package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ladecadence/EkiTelemetry/pkg/api"
	"github.com/ladecadence/EkiTelemetry/pkg/config"
	"github.com/ladecadence/EkiTelemetry/pkg/console"
	"github.com/ladecadence/EkiTelemetry/pkg/datalog"
	"github.com/ladecadence/EkiTelemetry/pkg/maindata"
	"github.com/ladecadence/EkiTelemetry/pkg/serialport"
	"github.com/ladecadence/EkiTelemetry/pkg/ssdv"
	"github.com/ladecadence/EkiTelemetry/pkg/ssdvimage"
	"github.com/ladecadence/EkiTelemetry/pkg/telemetry"
)

func clock(labelClock *widget.Label) {
	ticker := "🟠"
	for range time.Tick(time.Second) {
		fyne.Do(func() {
			labelClock.SetText(time.Now().Format(time.TimeOnly) + " " + ticker)
			if ticker == "🟠" {
				ticker = "🟢"
			} else {
				ticker = "🟠"
			}
		})
	}
}

func receiveTelemetry(config *config.Config, data chan string, console *console.Console, main *maindata.Main, log *datalog.DataLog, serial *serialport.Serial, labelStatus *widget.Label) {
	for {
		msg := <-data
		fyne.Do(func() { // fyne.Do for safe GUI updating between threads
			console.Append(time.Now().Format(time.TimeOnly) + " :: Received packet")
		})
		telem := telemetry.Telemetry{}
		err := telem.ParsePacket(msg, "/")
		if err != nil {
			fyne.Do(func() {
				console.Append(time.Now().Format(time.TimeOnly) + " :: Error decoding telemetry")
				labelStatus.SetText(fmt.Sprintf("Error decoding telemetry at %s", time.Now().Format(time.TimeOnly)))
			})
		} else {
			fyne.Do(func() {
				main.Update(telem)
				err = log.Append(datalog.Row{
					ID: serial.Packets, Name: telem.ID, Date: telem.Date, Time: telem.Time,
					Lat: telem.Lat, NS: telem.NS, Lon: telem.Lon, EW: telem.EW, Alt: telem.Alt, Sats: telem.Sats,
					Hdg: telem.Hdg, Spd: telem.Spd, Arate: telem.Arate,
					Tin: telem.Tin, Tout: telem.Tout, VBatt: telem.Vbat,
				})
				if err != nil {
					console.Append(time.Now().Format(time.TimeOnly) + fmt.Sprintf(" :: Error adding data to the log: %v", err))
				}
				console.Append(time.Now().Format(time.TimeOnly) +
					fmt.Sprintf(" :: Decoded telemetry: %.6f%s, %.6f%s, %.2f m, %d sats",
						telem.Lat, telem.NS, telem.Lon, telem.EW, telem.Alt, telem.Sats))
				labelStatus.SetText(fmt.Sprintf("Decoded telemetry packet at %s", time.Now().Format(time.TimeOnly)))

				// upload
				if config.Data.Server != "" && config.Data.User != "" && config.Data.Password != "" {
					api := api.API{Server: config.Data.Server, User: config.Data.User, Password: config.Data.Password}
					err := api.DataUpload(telem)
					if err != nil {
						labelStatus.SetText("Problem uploading telemetry.")
						console.Append(fmt.Sprintf(" :: Error uploading telemetry: %v", err))
					}
				}
			})
		}
	}
}

func receiveSSDV(ssdvChan chan []byte, console *console.Console, image *ssdvimage.SSDVImage, config *config.Config, labelStatus *widget.Label) {
	for {
		data := <-ssdvChan
		info := ssdv.SSDVPacketInfo(data)
		fyne.Do(func() {
			console.Append(time.Now().Format(time.TimeOnly) + fmt.Sprintf(" :: Received SSDV packet %d of image %d", info.Packet, info.Image))
			labelStatus.SetText(fmt.Sprintf("Received SSDV packet at %s", time.Now().Format(time.TimeOnly)))
		})
		imagePath, missionName, err := ssdv.SSDVDecodePacket(data, config.Data.ImageFolder)
		if err != nil {
			console.Append(fmt.Sprintf("Error decoding SSDV: %v", err))
		} else {
			fyne.Do(func() {
				if !info.LastPacket {
					image.UpdateImage(info.Image, imagePath, fmt.Sprintf("Receiving image %s...", imagePath), missionName)
				} else {
					image.UpdateImage(info.Image, imagePath, fmt.Sprintf("Received image %s.", imagePath), missionName)
				}
			})
		}
	}
}

func main() {
	dataChan := make(chan string)
	ssdvChan := make(chan []byte)

	// GUI
	telemApp := app.NewWithID("net.ladecadence.ekitelemetry")
	mainWindow := telemApp.NewWindow("EKI Telemetry")
	r, _ := fyne.LoadResourceFromPath("assets/eki-2.png")
	mainWindow.SetIcon(r)
	mainWindow.Resize(fyne.NewSize(900, 700))

	main := maindata.CreateMain()
	console := console.CreateConsole()
	config := config.CreateConfig(&mainWindow)
	log, err := datalog.CreateLog(config.Data.LogFile)
	if err != nil {
		fmt.Errorf("Can't create log file %v", err)
	}
	image := ssdvimage.CreateSSDVImage(config)

	serial, err := serialport.NewSerial(config.Data.PortName, 115200)
	if err != nil {
		fmt.Printf("Can't open serial port: %v", err)
	}
	defer serial.Close()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), main.Content),
		container.NewTabItemWithIcon("", theme.ComputerIcon(), console.Content),
		container.NewTabItemWithIcon("", theme.ListIcon(), log.Content),
		container.NewTabItemWithIcon("", theme.MediaPhotoIcon(), image.Content),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), config.Content),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	//tabs.SetTabLocation(container.TabLocationTop)

	// status bar
	labelStatus := widget.NewLabel("Waiting for data...")
	labelTime := widget.NewLabel("")
	statusBar := container.NewHBox(labelStatus, layout.NewSpacer(), labelTime)
	gui := container.NewBorder(nil, statusBar, nil, nil, tabs)

	mainWindow.SetContent(gui)

	// run receiver
	err = serial.ListenAndDecode(dataChan, ssdvChan)
	if err != nil {
		fmt.Printf("Problem with the serial port: %v", err)
	}

	// threads
	// telemetry
	go func() {
		receiveTelemetry(config, dataChan, console, main, log, serial, labelStatus)
	}()

	// ssdv
	go func() {
		receiveSSDV(ssdvChan, console, image, config, labelStatus)
	}()

	// clock
	go clock(labelTime)

	mainWindow.ShowAndRun()
}
