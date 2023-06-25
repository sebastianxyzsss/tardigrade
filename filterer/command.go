package filterer

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/sahilm/fuzzy"
	"github.com/sebastianxyzsss/tardigrade/globals"
	"github.com/sebastianxyzsss/tardigrade/parser"
	"github.com/sebastianxyzsss/tardigrade/utils"

	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run(element *parser.Element) (*parser.Element, error) {

	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.PromptStyle = o.PromptStyle.ToLipgloss()
	i.Placeholder = o.Placeholder
	i.Width = o.Width

	v := viewport.New(o.Width, o.Height)

	choices := getChoices(element)

	if len(choices) == 0 {
		return nil, errors.New("no options provided, see `gum filter --help`")
	}

	options := []tea.ProgramOption{tea.WithOutput(os.Stderr)}
	if o.Height == 0 {
		options = append(options, tea.WithAltScreen())
	}

	var matches []fuzzy.Match
	if o.Value != "" {
		i.SetValue(o.Value)
	}
	switch {
	case o.Value != "" && o.Fuzzy:
		matches = fuzzy.Find(o.Value, choices)
	case o.Value != "" && !o.Fuzzy:
		matches = exactMatches(o.Value, choices)
	default:
		matches = matchAll(choices)
	}

	if o.NoLimit {
		o.Limit = len(choices)
	}

	p := tea.NewProgram(model{
		element:               element,
		choices:               choices,
		indicator:             o.Indicator,
		matches:               matches,
		header:                o.Header,
		textinput:             i,
		viewport:              &v,
		indicatorStyle:        o.IndicatorStyle.ToLipgloss(),
		selectedPrefixStyle:   o.SelectedPrefixStyle.ToLipgloss(),
		selectedPrefix:        o.SelectedPrefix,
		unselectedPrefixStyle: o.UnselectedPrefixStyle.ToLipgloss(),
		unselectedPrefix:      o.UnselectedPrefix,
		matchStyle:            o.MatchStyle.ToLipgloss(),
		headerStyle:           o.HeaderStyle.ToLipgloss(),
		textStyle:             o.TextStyle.ToLipgloss(),
		height:                o.Height,
		selected:              make(map[string]struct{}),
		limit:                 o.Limit,
		reverse:               o.Reverse,
		fuzzy:                 o.Fuzzy,
	}, options...)

	tm, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("unable to run filter: %w", err)
	}
	m := tm.(model)
	if m.aborted {
		return nil, ErrAborted
	}

	isTTY := isatty.IsTerminal(os.Stdout.Fd())

	var chosen *parser.Element

	// allSelections contains values only if limit is greater
	// than 1 or if flag --no-limit is passed, hence there is
	// no need to further checks
	if len(m.selected) > 0 {
		for k := range m.selected {
			if isTTY {
				fmt.Println(k)
			} else {
				fmt.Println(Strip(k))
			}
		}
	} else if len(m.matches) > m.cursor && m.cursor >= 0 {
		index := m.cursor
		if isTTY {
			content := m.matches[index].Str
			if m.left {
				if m.element.Parent != nil {
					chosen = m.element.Parent
				} else {
					chosen = m.element
				}
			} else {
				chosen = element.Children[content]
			}
			if chosen.IsCommand {
				if parser.DoesElementHaveCommandTag(chosen.Parent) {
					content = utils.ReplaceContentWithChoices(chosen.Parent.Content, content)
				}
				fmt.Println(content)
			}
		} else {
			content := Strip(m.matches[index].Str)
			if m.left {
				if m.element.Parent != nil {
					chosen = m.element.Parent
				} else {
					chosen = m.element
				}
			} else {
				chosen = element.Children[content]
			}
			if chosen.IsCommand {
				if parser.DoesElementHaveCommandTag(chosen.Parent) {
					content = utils.ReplaceContentWithChoices(chosen.Parent.Content, content)
				}
				fmt.Println(content)
			}
		}
	}

	if !o.Strict && len(m.textinput.Value()) != 0 && len(m.matches) == 0 {
		fmt.Println(m.textinput.Value())
	}
	return chosen, nil
}

func getChoices(element *parser.Element) []string {
	var choices []string
	var choicesCompressed []string
	var newChildrenSorted []*parser.Element
	for _, child := range element.ChildrenSorted {
		if elementPassesPostCriteria(child) {
			choices = append(choices, child.Content)
			choicesCompressed = append(choicesCompressed, parser.TruncateString(child.Content, globals.ChildKeyMaxSize))
			newChildrenSorted = append(newChildrenSorted, child)
		}
	}
	element.ChildKeys = &choicesCompressed
	element.ChildrenSorted = newChildrenSorted
	return choices
}

func elementPassesPostCriteria(element *parser.Element) bool {
	if element.IsCommand == false && len(element.Children) < 1 {
		// log.Println("does not pass post criteria:", element.String(), len(element.Children))
		return false
	}
	return true
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}
