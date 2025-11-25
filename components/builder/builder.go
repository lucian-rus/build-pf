package builder

import (
	"fmt"
	"gobi/components/env"
	"log"
	"os"
	"os/exec"
)

func PrepareBuildCommands(libraryProperties env.LibraryProperties) error {
	// order is important - for now
	// @todo: change order once everything is done
	setCompiler(libraryProperties)

	setFlags(libraryProperties)
	setDefines(libraryProperties)
	setIncludes(libraryProperties)

	setOutputProperties(libraryProperties)
	setInputProperties(libraryProperties)

	return nil
}

func RunBuilder() error {
	var commandToRun string
	commandToRun += compiler
	for _, arg := range argumentList {
		commandToRun += " " + arg
	}
	log.Println(commandToRun)
	fmt.Println("running: ", commandToRun)

	cmd := exec.Command(compiler, argumentList...)

	var so saveOutput
	cmd.Stdout = &so
	cmd.Stderr = os.Stderr

	// clear the slice regardless of output
	argumentList = nil

	// should capture output of gcc command
	if err := cmd.Run(); err != nil {
		fmt.Println("encountered an error: ", err)
		fmt.Print(so.savedOutput)
		return err
	}

	return nil
}
