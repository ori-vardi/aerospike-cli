package prompt

import (
	"aerospike-cli/src/util"
	"fmt"
	"sort"
)

type Prompt struct {
	optionsMap     map[string]*Options
	OptionsSortKey []string
	runMainLoop    bool
}

type Options struct {
	key      string
	title    string
	function func()
}

func New() *Prompt {
	return &Prompt{optionsMap: map[string]*Options{}, runMainLoop: true}
}

func (prompt *Prompt) Add(key string, title string, fn func()) {
	option := &Options{
		key:      key,
		title:    title,
		function: fn,
	}

	prompt.optionsMap[key] = option
}

func (prompt *Prompt) printAllOptions() {
	for _, key := range prompt.OptionsSortKey {
		fmt.Printf("   (%s) %s\n", key, prompt.optionsMap[key].title)
	}
}

func (prompt *Prompt) Run() {
	prompt.sortPromptOptions()
	for prompt.runMainLoop {
		fmt.Println("\nSelect one of the options:\n-- MENU --")
		prompt.printAllOptions()
		input := util.GetKeyboardInput("[?] Enter menu selection: ")
		fmt.Printf("\n")
		if option, ok := prompt.optionsMap[input]; ok {
			option.function()
		}
	}
}

func (prompt *Prompt) sortPromptOptions() {
	keys := make([]string, 0)
	for k := range prompt.optionsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	prompt.OptionsSortKey = keys
}

func (prompt *Prompt) Close() {
	prompt.runMainLoop = false
}
