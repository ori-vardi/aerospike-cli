package prompt

import (
	"aerospike-cli/src/aerospike/client"
	"aerospike-cli/src/logger"
	"aerospike-cli/src/property"
	"aerospike-cli/src/util"
	"errors"
	"fmt"
	as "github.com/aerospike/aerospike-client-go"
	"strings"
)

const (
	defaultSet           = "aerospike.default.set"
	defaultNamespace     = "aerospike.default.namespace"
	scanAllChunkSize     = "aerospike.scan.all.chunk.size"
	allSetsKeysChunkSize = "aerospike.all.set.keys.chunk.size"
)

type ASPrompt struct {
	aerospike           *client.Aerospike
	setName             string
	nameSpace           string
	ConsoleString       *ConsoleString
	scanAllChunkSize    int
	allSetKeysChunkSize int
}

func New(env string) *ASPrompt {
	res := &ASPrompt{aerospike: client.New(env),
		ConsoleString:       NewConsoleString(),
		scanAllChunkSize:    property.Props.GetInt(scanAllChunkSize, 5),
		allSetKeysChunkSize: property.Props.GetInt(allSetsKeysChunkSize, 50),
		setName:             property.Props.GetString(defaultSet, ""),
	}
	setDefaultNamespace(res)
	return res
}

func setDefaultNamespace(res *ASPrompt) {
	defaultNameSpace := property.Props.MustGetString(defaultNamespace)
	if len(defaultNameSpace) > 0 {
		res.nameSpace = defaultNameSpace
	}
}

func (asPrompt *ASPrompt) ConfigSetName() {
	setName := util.GetKeyboardInput(asPrompt.ConsoleString.KeyboardInputInsertSetName)
	if len(setName) == 0 || setName == "c" {
		util.PrintInfoToConsoleNewLine("operation canceled")
		return
	}
	asPrompt.setName = setName
}

func (asPrompt *ASPrompt) GetRecord() {
	setName, keyName, cancel, err := asPrompt.getSetAndKey()
	if err != nil {
		util.PrintErrorToConsole(err.Error())
		return
	} else if cancel {
		util.PrintInfoToConsoleNewLine("operation canceled")
		return
	}
	asPrompt.getData(setName, keyName)
}

func (asPrompt *ASPrompt) getData(setName string, keyName string) error {
	record, err := asPrompt.aerospike.Find(asPrompt.nameSpace, setName, keyName)
	if err != nil {
		util.PrintErrorToConsole(asPrompt.ConsoleString.ConsolePrintGeneralFailed, "failed to get data", setName, keyName, err.Error())
		return err
	}
	util.PrintInfoToConsoleNewLine("data:\n%s", util.GetBinMapPrettyJson(record.Bins))
	return nil
}

func (asPrompt *ASPrompt) UpdateField() {
	setName, keyName, cancel, err := asPrompt.getSetAndKey()
	if err != nil {
		util.PrintErrorToConsole(err.Error())
		return
	} else if cancel {
		util.PrintInfoToConsoleNewLine("operation canceled")
		return
	}
	asPrompt.getData(setName, keyName)
	bins := asPrompt.getUpdateBinMapFromKeyboard()
	if bins != nil && len(bins) > 0 {
		err = asPrompt.aerospike.Update(asPrompt.nameSpace, setName, keyName, bins...)
		if err != nil {
			util.PrintErrorToConsole(asPrompt.ConsoleString.ConsolePrintGeneralFailed, "update fields failed", setName, keyName, err.Error())
			return
		}
	}
}

func (asPrompt *ASPrompt) DeleteField() {
	setName, keyName, cancel, err := asPrompt.getSetAndKey()
	if err != nil {
		util.PrintErrorToConsole(err.Error())
		return
	} else if cancel {
		util.PrintInfoToConsoleNewLine("operation canceled")
		return
	}
	err = asPrompt.getData(setName, keyName)
	if err != nil {
		return
	}
	fieldToDelete := util.GetKeyboardInput("Going to delete filed, insert field name, c for cancel: ")
	if len(fieldToDelete) == 0 || fieldToDelete == "c" {
		util.PrintInfoToConsoleNewLine("operation canceled")
		return
	}
	bin := &as.Bin{Name: fieldToDelete, Value: as.NewNullValue()}
	err = asPrompt.aerospike.Update(asPrompt.nameSpace, setName, keyName, bin)
	if err != nil {
		util.PrintErrorToConsole(asPrompt.ConsoleString.ConsolePrintGeneralFailed, "delete field failed", setName, keyName, err.Error())
		return
	}

}

func (asPrompt *ASPrompt) getUpdateBinMapFromKeyboard() []*as.Bin {
	flag := true
	var bins []*as.Bin
	for flag {
		key, val, done, cancel, err := asPrompt.getMultiplePairFromKeyboard()
		if err != nil {
			util.PrintErrorToConsole(err.Error())
			continue
		}
		if cancel {
			util.PrintInfoToConsoleNewLine("operation canceled")
			return nil
		}
		if done {
			flag = false
			continue
		}
		bin := &as.Bin{Name: key, Value: as.NewValue(val)}
		bins = append(bins, bin)
	}
	return bins
}

func (asPrompt *ASPrompt) getSetAndKey() (string, string, bool, error) {
	asPromptDefaultSet := asPrompt.setName
	input := util.GetKeyboardInput(fmt.Sprintf(asPrompt.ConsoleString.KeyboardInputInsertSetAndKey, asPromptDefaultSet))
	if len(input) > 0 {
		splitArr := strings.Split(input, " ")
		if len(splitArr) == 1 {
			if splitArr[0] == "c" {
				return "", "", true, nil
			} else if len(asPromptDefaultSet) == 0 {
				return "", "", false, errors.New(asPrompt.ConsoleString.ConsolePrintMissingSet)
			} else {
				return asPromptDefaultSet, splitArr[0], false, nil
			}
		} else if len(splitArr) == 2 {
			first := splitArr[0]
			sec := splitArr[1]
			return first, sec, false, nil
		}
	}
	return "", "", false, errors.New(asPrompt.ConsoleString.ConsolePrintInvalidInputSetAndKey)
}

func (asPrompt *ASPrompt) getMultiplePairFromKeyboard() (string, string, bool, bool, error) { //key, val, done, cancel, err
	input := util.GetKeyboardInput(asPrompt.ConsoleString.KeyboardInputInsertKeyAndVal)
	if len(input) == 0 {
		return "", "", true, false, nil
	} else {
		splitArr := strings.Split(input, " ")
		if len(splitArr) == 1 {
			return "", "", false, splitArr[0] == "c", nil
		} else if len(splitArr) == 2 {
			first := splitArr[0]
			sec := splitArr[1]
			return first, sec, false, false, nil
		}
		return "", "", false, false, errors.New(asPrompt.ConsoleString.ConsolePrintInvalidInputKeyAndVal)
	}
}

func (asPrompt *ASPrompt) Close() {
	asPrompt.aerospike.Close()
}

func (asPrompt *ASPrompt) ScanAll() {
	setName := util.GetKeyboardInput(fmt.Sprintf(asPrompt.ConsoleString.KeyboardInputInsertSetNameWithDefault, asPrompt.setName))
	if setName == "c" {
		util.PrintInfoToConsoleNewLine("operation canceled")
		return
	}
	if len(setName) == 0 {
		if len(asPrompt.setName) > 0 {
			setName = asPrompt.setName
		} else {
			util.PrintErrorToConsole(asPrompt.ConsoleString.ConsolePrintMissingSet)
			return
		}
	}
	recordset, err := asPrompt.aerospike.ScanAll(asPrompt.nameSpace, setName)
	if err != nil {
		util.PrintErrorToConsole(asPrompt.ConsoleString.ConsolePrintGeneralFailedSet, "failed", setName, err.Error())
		return
	}
	defer func(recordset *as.Recordset) {
		err := recordset.Close()
		if err != nil {
			logger.Error.Printf("Error while trying to close recordset, setName: %s, %s", setName, err)
		}
	}(recordset)
	chunkCounter := 0
	for res := range recordset.Results() {
		if res.Err != nil {
			logger.Error.Printf("ScanAll, skipping cause the record contain an error, %s", res.Err.Error())
			continue
		}
		if chunkCounter == asPrompt.scanAllChunkSize {
			chunkCounter = 0
			operation := util.GetKeyboardInput("press enter to continue, c for cancel: ")
			if operation == "c" {
				util.PrintInfoToConsoleNewLine("operation canceled")
				return
			}
		}
		data := map[string]interface{}{
			"Key":        res.Record.Key.Value(),
			"Bins":       res.Record.Bins,
			"Generation": res.Record.Generation,
			"Expiration": res.Record.Expiration,
		}
		util.PrintInfoToConsoleNewLine("data:\n%s", util.GetPrettyJson(data))
		chunkCounter++
	}
}

func (asPrompt *ASPrompt) ShowSets() {
	sets, err := asPrompt.aerospike.GetSets(asPrompt.nameSpace)
	if err != nil {
		util.PrintErrorToConsole(fmt.Sprintf("get all %ss failed, %s", asPrompt.ConsoleString.ConsoleSetNameAlias, err.Error()))
		return
	}
	util.PrintInfoToConsoleNewLine("sets (total: %d):", len(sets))
	for i := 0; i < len(sets); i += asPrompt.allSetKeysChunkSize {
		end := i + asPrompt.allSetKeysChunkSize
		if end > len(sets) {
			end = len(sets)
		}
		util.PrintInfoToConsoleNewLine("data: %d/%d:\n%s", end, len(sets), util.GetPrettyJson(sets[i:end]))
		if end != len(sets) {
			operation := util.GetKeyboardInput("press enter to continue, c for cancel: ")
			if operation == "c" {
				util.PrintInfoToConsoleNewLine("operation canceled")
				return
			}
		}
	}
}

func (asPrompt *ASPrompt) ShowKeys() {
	setName := util.GetKeyboardInput(fmt.Sprintf(asPrompt.ConsoleString.KeyboardInputInsertSetNameWithDefault, asPrompt.setName))
	if setName == "c" {
		util.PrintInfoToConsoleNewLine("operation canceled")
		return
	}
	if len(setName) == 0 {
		if len(asPrompt.setName) > 0 {
			setName = asPrompt.setName
		} else {
			util.PrintErrorToConsole(asPrompt.ConsoleString.ConsolePrintMissingSet)
			return
		}
	}
	recordset, err := asPrompt.aerospike.GetAllKeysBySet(asPrompt.nameSpace, setName)
	if err != nil {
		util.PrintErrorToConsole(asPrompt.ConsoleString.ConsolePrintGeneralFailedSet, "failed", setName, err.Error())
		return
	}
	defer func(recordset *as.Recordset) {
		err := recordset.Close()
		if err != nil {
			logger.Error.Printf("Error while trying to close recordset, setName: %s, %s", setName, err)
		}
	}(recordset)
	chunkCounter := 0

	var keys []string
	for res := range recordset.Results() {
		if chunkCounter == asPrompt.allSetKeysChunkSize {
			chunkCounter = 0
			util.PrintInfoToConsoleNewLine("data:\n%s", util.GetPrettyJson(keys[:]))
			keys = nil
			operation := util.GetKeyboardInput("press enter to continue, c for cancel: ")
			if operation == "c" {
				util.PrintInfoToConsoleNewLine("operation canceled")
				return
			}
		}
		if res.Err != nil {
			logger.Error.Printf("ShowKeys, skipping cause the record contain an error, %s", res.Err.Error())
			continue
		}
		if res.Record.Key.Value() != nil {
			keys = append(keys, res.Record.Key.Value().String())
			chunkCounter++
		}
	}
	if len(keys) > 0 {
		util.PrintInfoToConsoleNewLine("data:\n%s", util.GetPrettyJson(keys[:]))
	}
}
