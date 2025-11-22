package builder

import (
	"bytes"
	"fmt"
	"os/exec"
)

var (
	compiler     string
	argumentList []string
)

func PrepareBuildCommands(sourceFilesList []string) error {
	compiler = "gcc"

	// set output stuff
	argumentList = append(argumentList, "-o")
	argumentList = append(argumentList, "output")

	// set input stuff
	for _, sourceFile := range sourceFilesList {
		argumentList = append(argumentList, sourceFile)
	}

	return nil
}

func RunBuilder() error {
	var commandToRun string
	commandToRun += compiler
	for _, arg := range argumentList {
		commandToRun += " " + arg
	}
	fmt.Println(commandToRun)

	cmd := exec.Command(compiler, argumentList...)
	var cmdOutput bytes.Buffer
	cmd.Stdout = &cmdOutput

	// should capture output of gcc command
	if err := cmd.Run(); err != nil {
		fmt.Println("encountered an error: ", err)
		return err
	}

	fmt.Print(cmdOutput.String())
	return nil
}
