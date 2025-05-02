package telemetry

import (
	"errors"
	"strconv"
	"strings"
)

const (
	TELEMETRY_FIELDS = 14

	FIELD_LAT = 0
	FIELD_LON = 1
	FIELD_SPD = 2
	FIELD_ALT = 3
	FIELD_VLT = 4
	FIELD_BAR = 5
	FIELD_TIN = 6
	FIELD_TOU = 7
	FIELD_DAT = 8
	FIELD_TIM = 9
	FIELD_GPS = 10
	FIELD_SAT = 11
	FIELD_ASR = 12
	FIELD_MSG = 13
)

type Telemetry struct {
	ID    string
	Msg   string
	Lat   string
	NS    string
	Lon   string
	EW    string
	Alt   float64
	Hdg   string
	Spd   string
	Sats  int
	Vbat  string
	Baro  string
	Tin   string
	Tout  string
	Arate string
	Date  string
	Time  string
	Sep   string
	Hpwr  bool
}

func (t *Telemetry) ParsePacket(packet string, separator string) error {
	fields := strings.Split(packet, separator)
	if len(fields) >= TELEMETRY_FIELDS {
		t.Date = fields[FIELD_DAT]
		t.Time = fields[FIELD_TIM]
		t.Hdg = strings.Split(fields[FIELD_LON], "O")[1]
		t.Spd = fields[FIELD_SPD]
		t.Lat = strings.Split(strings.Split(fields[FIELD_GPS], "=")[1], ",")[0]
		t.Lon = strings.Split(strings.Split(fields[FIELD_GPS], "=")[1], ",")[1]
		alt, _ := strconv.ParseFloat(strings.Split(fields[FIELD_ALT], "=")[1], 64)
		t.Alt = alt
		t.Vbat = strings.Split(fields[FIELD_VLT], "=")[1]
		t.Tin = strings.Split(fields[FIELD_TIN], "=")[1]
		t.Tout = strings.Split(fields[FIELD_TOU], "=")[1]
		t.Baro = strings.Split(fields[FIELD_BAR], "=")[1]
		sats, _ := strconv.Atoi(strings.Split(fields[FIELD_SAT], "=")[1])
		t.Sats = sats
		t.Arate = strings.Split(fields[FIELD_ASR], "=")[1]
	} else {
		return errors.New("Not enough fields")
	}

	return nil
}
