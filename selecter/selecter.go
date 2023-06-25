package selecter

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/sebastianxyzsss/tardigrade/filterer"
	"github.com/sebastianxyzsss/tardigrade/globals"
	"github.com/sebastianxyzsss/tardigrade/logger"
	"github.com/sebastianxyzsss/tardigrade/parser"
	"github.com/sebastianxyzsss/tardigrade/reader"
	"github.com/sebastianxyzsss/tardigrade/utils"
)

var ll = logger.SetupLog()

func createOptions(settings *globals.Settings) *filterer.Options {
	ll.Debug().Msg("createOptions")

	filterOpts := filterer.Options{}

	filterOpts.Indicator = "•"
	filterOpts.IndicatorStyle = style.Styles{
		Foreground: settings.IndicatorStyle,
	}

	filterOpts.Header = "}}}}@ tardigrade"
	filterOpts.HeaderStyle = style.Styles{
		Foreground: "33",
	}

	filterOpts.Prompt = "> "
	filterOpts.Placeholder = "..."

	filterOpts.Limit = 1

	if globals.FlatParse {
		settings.Height += 4
	}

	filterOpts.Height = settings.Height

	filterOpts.MatchStyle = style.Styles{
		Foreground: "201",
	}

	filterOpts.PromptStyle = style.Styles{
		Foreground: "23",
	}

	return &filterOpts
}

func finalElementApply(rootElement *parser.Element, element *parser.Element) {
	ll.Debug().Msg("final element:" + element.String())

	if element.Parent != nil {
		if parser.DoesElementHaveCommandTag(element.Parent) {
			content := utils.ReplaceContentWithChoices(element.Parent.Content, element.Content)
			element.Content = content
		}
	}

	historyElement := parser.HistoryTemp

	childrenSorted := []*parser.Element{}
	childrenSorted = append(childrenSorted, element)

	for _, historyChild := range historyElement.ChildrenSorted {
		if historyChild.Content != element.Content {
			childrenSorted = append(childrenSorted, historyChild)
		}
		if len(childrenSorted) >= globals.HistorySize {
			break
		}
	}
	historyElement.ChildrenSorted = childrenSorted

	historyMap := parser.SimpleElementToMap(historyElement)

	yaml, err := reader.Marshall(&historyMap)
	if err != nil {
		log.Print(err)
		return
	}

	userTardiHistory, err := reader.GetHistoryPath()
	if err != nil {
		ll.Debug().Msg("unable to get history path: " + err.Error())
		return
	}

	reader.WriteToFile(*userTardiHistory, *yaml)
}

func Chooser(rootElement *parser.Element, settings *globals.Settings) {
	ll.Debug().Msg("about to do choosing ...")

	lipgloss.SetColorProfile(termenv.NewOutput(os.Stderr).Profile)

	ll.Debug().Msg("filter start")

	filterOpts := createOptions(settings)

	element := rootElement

	for {

		chosen, err := filterOpts.Run(element)
		if err != nil {
			ll.Debug().Msg("there was an interruption: " + err.Error())
			if strings.Contains(err.Error(), "no options provided") {
				ll.Debug().Msg("making an exception here")
				element = element.Parent
				continue
			} else {
				fmt.Println("pwd")
				break
			}
		}

		if chosen != nil {

			ll.Debug().Msg("chosen:" + chosen.String())

			if *&chosen.IsCommand {
				finalElementApply(rootElement, chosen)
				break
			}

			element = chosen

		} else {
			// fmt.Println("pwd")
			break
		}

	}

	ll.Debug().Msg("filter end")
}
