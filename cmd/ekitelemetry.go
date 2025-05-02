package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ladecadence/EkiTelemetry/pkg/config"
	"github.com/ladecadence/EkiTelemetry/pkg/console"
	"github.com/ladecadence/EkiTelemetry/pkg/datalog"
	"github.com/ladecadence/EkiTelemetry/pkg/maindata"
	"github.com/ladecadence/EkiTelemetry/pkg/serialport"
	"github.com/ladecadence/EkiTelemetry/pkg/ssdv"
	"github.com/ladecadence/EkiTelemetry/pkg/ssdvimage"
	"github.com/ladecadence/EkiTelemetry/pkg/telemetry"
)

func main() {
	dataChan := make(chan string)
	ssdvChan := make(chan []byte)

	// GUI
	telemApp := app.NewWithID("net.ladecadence.ekitelemetry")
	mainWindow := telemApp.NewWindow("EKI Telemetry")
	r, _ := fyne.LoadResourceFromPath("assets/eki-2.png")
	mainWindow.SetIcon(r)
	mainWindow.Resize(fyne.NewSize(800, 600))

	main := maindata.CreateMain()
	console := console.CreateConsole()
	config := config.CreateConfig(&mainWindow)
	log, err := datalog.CreateLog(config.Data.LogFile)
	if err != nil {
		fmt.Errorf("Can't create log file %v", err)
	}
	image := ssdvimage.CreateSSDVImage()

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

	labelStatus := widget.NewLabel("Waiting for data...")
	gui := container.NewBorder(nil, labelStatus, nil, nil, tabs)

	mainWindow.SetContent(gui)

	// run receiver
	err = serial.ListenAndDecode(dataChan, ssdvChan)
	if err != nil {
		fmt.Printf("Problem with the serial port: %v", err)
	}
	// telemetry
	go func() {
		for {
			msg := <-dataChan
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
				})
			}
		}
	}()

	// ssdv
	go func() {
		for {
			data := <-ssdvChan
			info := ssdv.SSDVPacketInfo(data)
			fyne.Do(func() {
				console.Append(time.Now().Format(time.TimeOnly) + fmt.Sprintf(" :: Received SSDV packet %d of image %d", info.Packet, info.Image))
				labelStatus.SetText(fmt.Sprintf("Received SSDV packet at %s", time.Now().Format(time.TimeOnly)))
			})
			imagePath, err := ssdv.SSDVDecodePacket(data, config.Data.ImageFolder)
			if err != nil {
				console.Append(fmt.Sprintf("Error decoding SSDV: %v", err))
			} else {
				fyne.Do(func() {
					if !info.LastPacket {
						image.UpdateImage(imagePath, fmt.Sprintf("Receiving image %s...", imagePath))
					} else {
						image.UpdateImage(imagePath, fmt.Sprintf("Received image %s.", imagePath))
					}
				})
			}
		}
	}()

	mainWindow.ShowAndRun()
}
