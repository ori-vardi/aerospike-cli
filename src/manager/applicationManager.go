package manager

import (
	aerospikePrompt "aerospike-cli/src/aerospike/prompt"
	"aerospike-cli/src/prompt"
	"aerospike-cli/src/property"
	"aerospike-cli/src/util"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	preCmd    = "app.manager.pre.cmd"
	preCmdArg = "app.manager.pre.cmd.arg"
)

type ApplicationManager struct {
	quitChannel chan os.Signal
	asPrompt    *aerospikePrompt.ASPrompt
	prompt      *prompt.Prompt
	preRunCmd   *exec.Cmd
}

func New(env string) *ApplicationManager {
	preRunCmd := preRunCmd()
	return &ApplicationManager{
		quitChannel: registerShutdownSignals(),
		asPrompt:    aerospikePrompt.New(env),
		prompt:      prompt.New(),
		preRunCmd:   preRunCmd,
	}
}

func (am *ApplicationManager) Start() {
	go am.run()
}

func (am *ApplicationManager) shutdown() {
	util.PrintInfoToConsoleNewLine("\nShutting down the service...")
	am.prompt.Close()
	am.asPrompt.Close()
	if am.preRunCmd != nil {
		util.PrintInfoToConsoleNewLine("kill pre run cmd, ps: %d", am.preRunCmd.Process.Pid)
		err := am.preRunCmd.Process.Kill()
		if err != nil {
			util.PrintErrorToConsole("failed to kill process: %s", err)
		}
	}

	util.PrintInfoToConsoleNewLine("goodbye 😀")
	os.Exit(0)
}

func registerShutdownSignals() chan os.Signal {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return interruptChan
}

func (am *ApplicationManager) run() {
	go am.runPrompt()
	<-am.quitChannel
	am.shutdown()
}

func (am *ApplicationManager) runPrompt() {
	am.prompt.Add("1", am.asPrompt.ConsoleString.MenuGetRecord, am.asPrompt.GetRecord)
	am.prompt.Add("2", am.asPrompt.ConsoleString.MenuUpdateField, am.asPrompt.UpdateField)
	am.prompt.Add("3", am.asPrompt.ConsoleString.MenuDeleteField, am.asPrompt.DeleteField)
	am.prompt.Add("4", am.asPrompt.ConsoleString.MenuScanAll, am.asPrompt.ScanAll)
	am.prompt.Add("5", am.asPrompt.ConsoleString.MenuShowSets, am.asPrompt.ShowSets)
	am.prompt.Add("6", am.asPrompt.ConsoleString.MenuShowKeys, am.asPrompt.ShowKeys)
	am.prompt.Add("7", am.asPrompt.ConsoleString.MenuConfigDefaultSet, am.asPrompt.ConfigSetName)
	am.prompt.Add("0", "exit", am.shutdown)
	am.prompt.Run()
}

func preRunCmd() *exec.Cmd {
	preCmd, ok := property.Props.Get(preCmd)
	if ok {
		cmdArg := property.Props.GetString(preCmdArg, "")
		cmdArgArr := strings.Split(cmdArg, " ")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, preCmd, cmdArgArr...)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if ctx.Err() == context.DeadlineExceeded {
			util.PrintErrorToConsole(ctx.Err().Error())
			return nil
		}
		if err != nil {
			util.PrintErrorToConsole(err.Error())
		}
		return cmd
	}
	return nil
}
