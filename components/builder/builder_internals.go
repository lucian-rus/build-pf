package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// @todo this is taken over from stack, should be properly handled
type saveOutput struct {
	savedOutput []byte
}

// internals
var (
	compiler     string
	argumentList []string
)

func setCompiler(libraryProprties LibraryProperties) {
	compiler = "gcc"
}

func setFlags(libraryProprties LibraryProperties) {
	argumentList = append(argumentList, libraryProprties.Flags...)
}

func setDefines(libraryProprties LibraryProperties) {

}

func setIncludes(libraryProprties LibraryProperties) {
	for _, item := range libraryProprties.Includes.Public {
		parsedArgument := "-I" + item
		argumentList = append(argumentList, parsedArgument)
	}

	for _, item := range libraryProprties.Includes.Private {
		parsedArgument := "-I" + item
		argumentList = append(argumentList, parsedArgument)
	}
}

func setOutputProperties(libraryProprties LibraryProperties) {
	// set output stuff
	argumentList = append(argumentList, "-o")
	argumentList = append(argumentList, libraryProprties.Name)
}

func setInputProperties(libraryProprties LibraryProperties) {
	// set target source file
	argumentList = append(argumentList, libraryProprties.Sources...)
	argumentList = append(argumentList, libraryProprties.Dependencies.Libraries...)
}

func runBuilder() error {
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

// consider doing this our own
func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}
