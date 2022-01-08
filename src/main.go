package main

import (
	"aerospike-cli/src/logger"
	"aerospike-cli/src/manager"
	"aerospike-cli/src/property"
	"aerospike-cli/src/util"
	"fmt"
	"os"
)

const (
	defaultEnv = "aerospike.default.env"
)

func main() {
	quit := make(chan struct{})
	logger.Init(property.Props)
	env := getEnv()
	applicationManager := manager.New(env)
	applicationManager.Start()
	<-quit
}

func getEnv() string {
	defaultEnvProp := property.Props.GetString(fmt.Sprintf(defaultEnv), "")
	env := util.GetKeyboardInput(fmt.Sprintf("select env (%s) or press enter for default (%s): ", property.GetEnvOptions(), defaultEnvProp))
	if len(env) == 0 {
		if len(defaultEnvProp) > 0 {
			env = defaultEnvProp
		} else {
			util.PrintErrorToConsole("insert env, or config a default one in the property file")
			os.Exit(-1)
		}
	}
	return env
}
