package maindata

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ladecadence/EkiTelemetry/pkg/telemetry"
)

type Main struct {
	Content   *fyne.Container
	DataGPS   *widget.Hyperlink
	DataAlt   *widget.Label
	DataSats  *widget.Label
	DataBaro  *widget.Label
	DataTin   *widget.Label
	DataTout  *widget.Label
	DataBatt  *widget.Label
	DataHdg   *widget.Label
	DataSpd   *widget.Label
	DataArate *widget.Label
}

func CreateMain() *Main {
	main := Main{}

	// create GUI
	labelGPS := widget.NewLabel("GPS:")
	dataGPS := widget.NewHyperlink("", nil)
	labelAlt := widget.NewLabel("Altitude:")
	dataAlt := widget.NewLabel("")
	labelSats := widget.NewLabel("Satellites:")
	dataSats := widget.NewLabel("")
	labelBaro := widget.NewLabel("Barometer:")
	dataBaro := widget.NewLabel("")
	labelTin := widget.NewLabel("Internal temperature:")
	dataTin := widget.NewLabel("")
	labelTout := widget.NewLabel("External temperature:")
	dataTout := widget.NewLabel("")
	labelBatt := widget.NewLabel("Battery voltage:")
	dataBatt := widget.NewLabel("")
	labelHdg := widget.NewLabel("Heading:")
	dataHdg := widget.NewLabel("")
	labelSpd := widget.NewLabel("Speed:")
	dataSpd := widget.NewLabel("")
	labelArate := widget.NewLabel("Ascension rate:")
	dataArate := widget.NewLabel("")

	logo := canvas.NewImageFromFile("assets/eki-2.png")
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 200))

	form := container.New(layout.NewFormLayout(),
		labelGPS,
		dataGPS,
		labelAlt,
		dataAlt,
		labelSats,
		dataSats,
		widget.NewSeparator(),
		widget.NewSeparator(),
		labelHdg,
		dataHdg,
		labelSpd,
		dataSpd,
		labelArate,
		dataArate,
		widget.NewSeparator(),
		widget.NewSeparator(),
		labelBaro,
		dataBaro,
		labelTin,
		dataTin,
		labelTout,
		dataTout,
		labelBatt,
		dataBatt,
	)

	padded := container.NewBorder(nil, nil, nil, logo, form)
	content := container.NewPadded(padded)

	main.DataGPS = dataGPS
	main.DataAlt = dataAlt
	main.DataBaro = dataBaro
	main.DataSats = dataSats
	main.DataTin = dataTin
	main.DataTout = dataTout
	main.DataBatt = dataBatt
	main.DataArate = dataArate
	main.DataHdg = dataHdg
	main.DataSpd = dataSpd
	main.Content = content

	return &main
}

func (m *Main) Update(t telemetry.Telemetry) {
	m.DataGPS.SetText(t.Lat + t.NS + ", " + t.Lon + t.EW)
	url, _ := url.Parse("http://maps.google.com/maps?z=12&t=m&q=loc:" + t.Lat + "+" + t.Lon)
	m.DataGPS.SetURL(url)
	m.DataAlt.SetText(fmt.Sprintf("%.2f m", t.Alt))
	m.DataBaro.SetText(fmt.Sprintf("%.2f mBar", t.Baro))
	m.DataSats.SetText(fmt.Sprintf("%d", t.Sats))
	m.DataHdg.SetText(fmt.Sprintf("%.2f º", t.Hdg))
	m.DataSpd.SetText(fmt.Sprintf("%.2f kn", t.Spd))
	m.DataArate.SetText(fmt.Sprintf("%.2f m/s", t.Arate))
	m.DataTin.SetText(fmt.Sprintf("%.2f ºC", t.Tin))
	m.DataTout.SetText(fmt.Sprintf("%.2f ºC", t.Tout))
	m.DataBatt.SetText(fmt.Sprintf("%.2f V", t.Vbat))
}
