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
	ID    string  `json:"id"`
	Msg   string  `json:"msg"`
	Lat   float64 `json:"lat"`
	NS    string  `json:"ns"`
	Lon   float64 `json:"lon"`
	EW    string  `json:"ew"`
	Alt   float64 `json:"alt"`
	Hdg   float64 `json:"hdg"`
	Spd   float64 `json:"spd"`
	Sats  int     `json:"sats"`
	Vbat  float64 `json:"vbat"`
	Baro  float64 `json:"baro"`
	Tin   float64 `json:"tin"`
	Tout  float64 `json:"tout"`
	Arate float64 `json:"arate"`
	Date  string  `json:"date"`
	Time  string  `json:"time"`
	Sep   string  `json:"sep"`
	Hpwr  bool    `json:"hpwr"`
}

func (t *Telemetry) ParsePacket(packet string, separator string) error {
	fields := strings.Split(packet, separator)
	if len(fields) >= TELEMETRY_FIELDS {
		var err error
		t.ID = strings.Split(fields[FIELD_LAT], "!")[0][2:]
		t.Date = fields[FIELD_DAT]
		t.Time = fields[FIELD_TIM]
		hdg, err := strconv.ParseFloat(strings.Split(fields[FIELD_LON], "O")[1], 64)
		t.Hdg = hdg
		spd, err := strconv.ParseFloat(fields[FIELD_SPD], 64)
		t.Spd = spd
		lat := strings.Split(strings.Split(fields[FIELD_GPS], "=")[1], ",")[0]
		latNum, err := strconv.ParseFloat(lat[:len(lat)-1], 64)
		t.Lat = latNum
		t.NS = string(lat[len(lat)-1])
		lon := strings.Split(strings.Split(fields[FIELD_GPS], "=")[1], ",")[1]
		lonNum, err := strconv.ParseFloat(lon[:len(lon)-1], 64)
		t.Lon = lonNum
		t.EW = string(lon[len(lon)-1])
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
