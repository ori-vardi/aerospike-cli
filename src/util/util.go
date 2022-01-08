package util

import (
	"aerospike-cli/src/logger"
	"bufio"
	"encoding/json"
	"fmt"
	as "github.com/aerospike/aerospike-client-go"
	jsoniter "github.com/json-iterator/go"
	"os"
	"strings"
)

func GetKeyboardInput(title string) string {
	reader := bufio.NewReader(os.Stdin)
	PrintInfoToConsole(title)
	inputStr, _ := reader.ReadString('\n')
	input := strings.TrimSuffix(inputStr, "\n")
	logger.Info.Println("input from stdin: ", input)
	return input
}

func GetBinMapPrettyJson(binMap as.BinMap) string {
	var jsoniterStandard = jsoniter.ConfigCompatibleWithStandardLibrary
	marshal, err := jsoniterStandard.MarshalIndent(binMap, "", "    ")
	if err != nil {
		logger.Error.Println("json marshal failed, dataMap: %s", marshal, err)
	}
	return string(marshal)
}

func GetPrettyJson(dataMap interface{}) string {
	indent, err := json.MarshalIndent(dataMap, "", "    ")
	if err != nil {
		logger.Error.Println("json marshal failed, dataMap: %s", dataMap, err)
		return "json marshal failed"
		//TODO:ORI:fix
	}
	return string(indent)
}

func PrintInfoToConsole(title string, arg ...interface{}) {
	var log string
	if len(arg) > 0 {
		log = fmt.Sprintf(title, arg...)
	} else {
		log = title
	}
	logger.Info.Print(log)
	fmt.Print(log)
}

func PrintInfoToConsoleNewLine(title string, arg ...interface{}) {
	var log string
	if len(arg) > 0 {
		log = fmt.Sprintf(title, arg...)
	} else {
		log = title
	}
	logger.Info.Println(log)
	fmt.Println(log)
}

func PrintErrorToConsole(title string, objArr ...interface{}) {
	log := fmt.Sprintf(title, objArr...)
	logger.Error.Println(log)
	fmt.Println(log)
}
