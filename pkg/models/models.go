package models

import (
	"dem3_demo_v2/pkg/logging"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")
var ErrDoubleRecord = errors.New("models: такая запись уже есть")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type ProfData struct {
	ID                    int
	D2Number              string
	D2Profstroi           string
	D2Object              string
	D2Manager             string
	D2Kontragent          string
	D2ID                  string
	D2Diler               string
	D2City                string
	D2Napr                string
	D2SumProjToSk         string
	D2SumSkidka           string
	D2SumProjWithSkidka   string
	D2SumConstrWithSkidka string
	D2SumRabWithSkidka    string
	D2Status              string
	NoteOrder             string
	Logger                logging.Logger

	Details
}

type Details struct {
	ID         int
	OrderID    string
	Size       string
	Name       string
	Count      string
	Allowances string
	Color      string
	Type       string
	Price      string
	Location   string
	Width      string
	Height     string
	GForm      string
	NSq        string
	NRama      string
	AUg01      string
	AUg02      string
	CLke       string
}
