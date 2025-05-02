package datalog

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	numColumns = 6
)

type Row struct {
	ID       int
	Time     string
	Date     string
	Lat, Lon string
	Alt      float64
	Sats     int
	Baro     string
	VBatt    string
	Tin      string
	Tout     string
	Hdg      string
	Spd      string
	Arate    string
}

type DataLog struct {
	FilePath string
	Content  *fyne.Container
	Table    *widget.Table
	Rows     []Row
}

func CreateLog(file string) (*DataLog, error) {
	log := DataLog{}

	// check file, if exists use it
	if _, err := os.Stat(file); err == nil {
		log.FilePath = file
	} else {
		// create it
		f, err := os.Create(file)
		if err != nil {
			return nil, err
		}
		f.Close()
	}

	// GUI
	table := widget.NewTableWithHeaders(func() (int, int) {
		return len(log.Rows), numColumns
	},
		func() fyne.CanvasObject {
			l := widget.NewLabel("Data")
			return l
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			l.Truncation = fyne.TextTruncateEllipsis
			switch id.Col {
			case 0:
				l.Truncation = fyne.TextTruncateOff
				l.SetText(strconv.Itoa(log.Rows[id.Row].ID))
			case 1:
				l.SetText(log.Rows[id.Row].Time)
			case 2:
				l.SetText(log.Rows[id.Row].Lat)
			case 3:
				l.SetText(log.Rows[id.Row].Lon)
			case 4:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Alt, 'f', 2, 64))
			case 5:
				l.SetText(strconv.Itoa(log.Rows[id.Row].Sats))
			}

		})

	table.SetColumnWidth(0, 40)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 100)
	table.SetColumnWidth(5, 30)

	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewButton("000", func() {})
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			switch id.Col {
			case 0:
				b.SetText("ID")
			case 1:
				b.SetText("Time")
			case 2:
				b.SetText("Latitude")
			case 3:
				b.SetText("Longitude")
			case 4:
				b.SetText("Altitude")
			case 5:
				b.SetText("Sats")
			}
			b.Importance = widget.MediumImportance
			b.Enable()
			b.Refresh()
		}
	}

	log.Table = table

	padded := container.NewBorder(nil, nil, nil, nil, log.Table)
	content := container.NewPadded(padded)
	log.Content = content
	return &log, nil
}

func (l *DataLog) Append(row Row) error {
	l.Rows = append(l.Rows, row)
	file, err := os.OpenFile(l.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("%s,%s,%s,%s,%.2f,%d,%s,%s,%s,%s,%s,%s\n",
		row.Date, row.Time, row.Lat, row.Lon, row.Alt, row.Sats, row.Hdg,
		row.Spd, row.Arate, row.Tin, row.Tout, row.VBatt,
	))
	err = file.Close()
	return err
}
