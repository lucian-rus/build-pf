package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

// @todo add support for library-specific compiler. until then, it shall be defaulted.
func setCompiler(libProprties LibraryProperties) {
	compiler = ProjectConfiguration.Compiler
}

func setFlags(libProprties LibraryProperties) {
	argumentList = append(argumentList, libProprties.Flags...)
}

func setDefines(libProprties LibraryProperties) {
	for _, item := range libProprties.Defines {
		parsedArgument := "-D" + item
		argumentList = append(argumentList, parsedArgument)
	}
}

func setIncludes(libProprties LibraryProperties) {
	for _, item := range libProprties.Includes.Public {
		parsedArgument := "-I" + item
		argumentList = append(argumentList, parsedArgument)
	}

	for _, item := range libProprties.Includes.Private {
		parsedArgument := "-I" + item
		argumentList = append(argumentList, parsedArgument)
	}
}

func setOutputProperties(libProprties LibraryProperties) {
	// set output stuff
	argumentList = append(argumentList, "-o")

	// @todo need to resolve global dependency for output
	objOutputPath := filepath.Join(ProjectConfiguration.OutputPath, libProprties.Name)
	argumentList = append(argumentList, objOutputPath)
}

func setInputProperties(libProprties LibraryProperties) {
	// set target source file
	argumentList = append(argumentList, libProprties.Sources...)
	argumentList = append(argumentList, libProprties.LinkedObjects...)
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
