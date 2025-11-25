package builder

import (
	"gobi/components/env"
	"os"
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

func setCompiler(libraryProprties env.LibraryProperties) {
	compiler = "gcc"
}

func setFlags(libraryProprties env.LibraryProperties) {
	for _, item := range libraryProprties.Flags {
		argumentList = append(argumentList, item)
	}
}

func setDefines(libraryProprties env.LibraryProperties) {

}

func setIncludes(libraryProprties env.LibraryProperties) {
	for _, item := range libraryProprties.Includes.Public {
		parsedArgument := "-I" + item
		argumentList = append(argumentList, parsedArgument)
	}

	for _, item := range libraryProprties.Includes.Private {
		parsedArgument := "-I" + item
		argumentList = append(argumentList, parsedArgument)
	}
}

func setOutputProperties(libraryProprties env.LibraryProperties) {
	// set output stuff
	argumentList = append(argumentList, "-o")
	argumentList = append(argumentList, libraryProprties.Name)
}

func setInputProperties(libraryProprties env.LibraryProperties) {
	// set target source file

	for _, sourceFile := range libraryProprties.Sources {
		argumentList = append(argumentList, sourceFile)
	}

	for _, library := range libraryProprties.Dependencies.Libraries {
		argumentList = append(argumentList, library)
	}
}

// consider doing this our own
func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}
