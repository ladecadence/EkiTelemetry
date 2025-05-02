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
	Hdg   float64
	Spd   float64
	Sats  int
	Vbat  float64
	Baro  float64
	Tin   float64
	Tout  float64
	Arate float64
	Date  string
	Time  string
	Sep   string
	Hpwr  bool
}

func (t *Telemetry) ParsePacket(packet string, separator string) error {
	fields := strings.Split(packet, separator)
	if len(fields) >= TELEMETRY_FIELDS {
		var err error
		t.Date = fields[FIELD_DAT]
		t.Time = fields[FIELD_TIM]
		hdg, err := strconv.ParseFloat(strings.Split(fields[FIELD_LON], "O")[1], 64)
		t.Hdg = hdg
		spd, err := strconv.ParseFloat(fields[FIELD_SPD], 64)
		t.Spd = spd
		t.Lat = strings.Split(strings.Split(fields[FIELD_GPS], "=")[1], ",")[0]
		t.Lon = strings.Split(strings.Split(fields[FIELD_GPS], "=")[1], ",")[1]
		alt, err := strconv.ParseFloat(strings.Split(fields[FIELD_ALT], "=")[1], 64)
		t.Alt = alt
		vbatt, err := strconv.ParseFloat(strings.Split(fields[FIELD_VLT], "=")[1], 64)
		t.Vbat = vbatt
		tin, err := strconv.ParseFloat(strings.Split(fields[FIELD_TIN], "=")[1], 64)
		t.Tin = tin
		tout, err := strconv.ParseFloat(strings.Split(fields[FIELD_TOU], "=")[1], 64)
		t.Tout = tout
		baro, err := strconv.ParseFloat(strings.Split(fields[FIELD_BAR], "=")[1], 64)
		t.Baro = baro
		sats, _ := strconv.Atoi(strings.Split(fields[FIELD_SAT], "=")[1])
		t.Sats = sats
		arate, err := strconv.ParseFloat(strings.Split(fields[FIELD_ASR], "=")[1], 64)
		t.Arate = arate
		return err
	} else {
		return errors.New("Not enough fields")
	}
}
