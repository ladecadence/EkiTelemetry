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
	Name     string
	Time     string
	Date     string
	Lat, Lon float64
	NS, EW   string
	Alt      float64
	Sats     int
	Baro     float64
	VBatt    float64
	Tin      float64
	Tout     float64
	Hdg      float64
	Spd      float64
	Arate    float64
}

type DataLog struct {
	FilePath string
	Content  *fyne.Container
	Table    *widget.Table
	Rows     []Row
}

func CreateLog(file string) (*DataLog, error) {
	log := DataLog{}

	// // check file, if exists use it
	// if _, err := os.Stat(file); err == nil {
	// 	log.FilePath = file
	// } else {
	// 	// create it
	// 	f, err := os.Create(file)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	f.Close()
	// }
	log.FilePath = file

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
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Lat, 'f', 6, 64) + log.Rows[id.Row].NS)
			case 3:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Lon, 'f', 6, 64) + log.Rows[id.Row].EW)
			case 4:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Alt, 'f', 2, 64))
			case 5:
				l.SetText(strconv.Itoa(log.Rows[id.Row].Sats))
			case 6:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Baro, 'f', 1, 64))
			case 7:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Hdg, 'f', 2, 64))
			case 8:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Spd, 'f', 2, 64))
			case 9:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Arate, 'f', 2, 64))
			case 10:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Tin, 'f', 2, 64))
			case 11:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].Tout, 'f', 2, 64))
			case 12:
				l.SetText(strconv.FormatFloat(log.Rows[id.Row].VBatt, 'f', 2, 64))

			}

		})

	table.SetColumnWidth(0, 40)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 100)
	table.SetColumnWidth(5, 30)
	table.SetColumnWidth(6, 50)
	table.SetColumnWidth(7, 30)
	table.SetColumnWidth(8, 30)
	table.SetColumnWidth(9, 30)
	table.SetColumnWidth(10, 30)
	table.SetColumnWidth(11, 30)
	table.SetColumnWidth(12, 30)

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
	_, err = file.WriteString(fmt.Sprintf("%s,%s,%s,%.6f%s,%.6f%s,%.6f,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f\n",
		row.Name, row.Date, row.Time, row.Lat, row.NS, row.Lon, row.EW, row.Alt, row.Sats, row.Hdg,
		row.Spd, row.Arate, row.Tin, row.Tout, row.VBatt,
	))
	err = file.Close()
	return err
}
