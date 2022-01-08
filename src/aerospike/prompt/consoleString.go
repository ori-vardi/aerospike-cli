package prompt

import (
	"aerospike-cli/src/property"
	"fmt"
)

const (
	invalidInput                 = "invalid input, one or more of the following fields are missing: "
	cancel                       = "c for cancel: "
	aerospikeConsoleSetNameAlias = "aerospike.console.set.name.alias"
	aerospikeConsoleKeyNameAlias = "aerospike.console.key.name.alias"
)

type ConsoleString struct {
	ConsoleSetNameAlias string
	ConsoleKeyNameAlias string

	MenuGetRecord              string
	MenuConfigDefaultSet       string
	MenuGetDataByKeyDefaultSet string
	MenuUpdateField            string
	MenuDeleteField            string
	MenuScanAll                string
	MenuShowSets               string
	MenuShowKeys               string

	KeyboardInputInsertSetName            string
	KeyboardInputInsertSetNameWithDefault string
	KeyboardInputInsertKey                string
	KeyboardInputInsertSetAndKey          string
	KeyboardInputInsertKeyAndVal          string

	ConsolePrintMissingSet            string
	ConsolePrintGeneralFailed         string
	ConsolePrintGeneralFailedSet      string
	ConsolePrintInvalidInputSetAndKey string
	ConsolePrintInvalidInputKeyAndVal string
}

func NewConsoleString() *ConsoleString {
	cs := &ConsoleString{
		ConsoleSetNameAlias: property.Props.GetString(aerospikeConsoleSetNameAlias, "set"),
		ConsoleKeyNameAlias: property.Props.GetString(aerospikeConsoleKeyNameAlias, "key"),
	}
	configMenuTitle(cs)
	configKeyboardInput(cs)
	configConsolePrint(cs)
	return cs
}

func configConsolePrint(cs *ConsoleString) {
	cs.ConsolePrintMissingSet = fmt.Sprintf("%s needed, you can configure it using config option or in the prop file", cs.ConsoleSetNameAlias)
	cs.ConsolePrintGeneralFailed = fmt.Sprintf("%%s, %s: %%s, %s: %%s, error: %%s", cs.ConsoleSetNameAlias, cs.ConsoleKeyNameAlias)
	cs.ConsolePrintGeneralFailedSet = fmt.Sprintf("%%s, %s: %%s, error: %%s", cs.ConsoleSetNameAlias)
	cs.ConsolePrintInvalidInputSetAndKey = fmt.Sprintf("%s %s, %s", invalidInput, cs.ConsoleSetNameAlias, cs.ConsoleKeyNameAlias)
	cs.ConsolePrintInvalidInputKeyAndVal = fmt.Sprintf("%s %s, value", invalidInput, cs.ConsoleKeyNameAlias)
}

func configKeyboardInput(cs *ConsoleString) {
	cs.KeyboardInputInsertSetName = fmt.Sprintf("Insert %[1]v, %[2]v", cs.ConsoleSetNameAlias, cancel)
	cs.KeyboardInputInsertSetNameWithDefault = fmt.Sprintf("Press enter for the default %[1]v (%%s) or insert %[1]v, %[2]v", cs.ConsoleSetNameAlias, cancel)
	cs.KeyboardInputInsertKey = fmt.Sprintf("Insert %s: ", cs.ConsoleKeyNameAlias)
	cs.KeyboardInputInsertSetAndKey = fmt.Sprintf("Insert %[2]v (default %[1]v: %%s) or\ninsert %[1]v and %[2]v (split by space),\n%[3]v", cs.ConsoleSetNameAlias, cs.ConsoleKeyNameAlias, cancel)
	cs.KeyboardInputInsertKeyAndVal = fmt.Sprintf("Insert %s: and value (split by space), enter to finish, %s", cs.ConsoleKeyNameAlias, cancel)
}

func configMenuTitle(cs *ConsoleString) {
	cs.MenuGetRecord = fmt.Sprintf("Get record")
	cs.MenuConfigDefaultSet = fmt.Sprintf("Config default %s", cs.ConsoleSetNameAlias)
	cs.MenuUpdateField = fmt.Sprintf("Update record's field (bin)")
	cs.MenuDeleteField = fmt.Sprintf("Delete record's field (bin)")
	cs.MenuScanAll = fmt.Sprintf("Scan %s's records", cs.ConsoleSetNameAlias)
	cs.MenuShowKeys = fmt.Sprintf("Get %s's %ss (in case of sendKey (true) was seted at creation time)", cs.ConsoleSetNameAlias, cs.ConsoleKeyNameAlias)
	cs.MenuShowSets = fmt.Sprintf("Get %ss", cs.ConsoleSetNameAlias)
}
