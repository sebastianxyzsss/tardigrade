package globals

import "github.com/sebastianxyzsss/tardigrade/action"

type FilterType int

const (
	FilterNone FilterType = iota
	FilterFiles
	FilterTags
	FilterAnything
)

var FilterAction FilterType = FilterNone

var FilterStrings []string = make([]string, 0)

var RunMode string = "print-command"

var RunAction action.Action = action.Printer{}

type Settings struct {
	Height           int    `json:"height"`
	HistorySize      int    `json:"historysize"`
	IndicatorStyle   string `json:"indicatorstyle"`
	FooterKeyMaxSize int    `json:"footerkeymaxsize"`
	LogLevel         string `json:"loglevel"`
}

var ChildKeyMaxSize int = 8

var FlatParse bool = false

var HistorySize int = 10
