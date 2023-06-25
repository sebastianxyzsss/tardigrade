package globals

type FilterType int

const (
	FilterNone FilterType = iota
	FilterFiles
	FilterTags
	FilterAnything
)

var FilterAction FilterType = FilterNone

var FilterStrings []string = make([]string, 0)

type Settings struct {
	Height           int    `json:"height"`
	HistorySize      int    `json:"historysize"`
	IndicatorStyle   string `json:"indicatorstyle"`
	FooterKeyMaxSize int    `json:"footerkeymaxsize"`
}

var ChildKeyMaxSize int = 8

var FlatParse bool = false

var HistorySize int = 10
