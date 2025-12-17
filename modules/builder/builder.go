package builder

import (
	"fmt"
	"os"
	"os/exec"
)

// @todo this is taken over from stack, should be properly handled
type saveOutput struct {
	savedOutput []byte
}

var (
	// internal
	commandSequence [][]string
)

func AddBuildSequence(sequence []string) {
	commandSequence = append(commandSequence, sequence)
}

func Build() error {
	fmt.Println("---------------- running build ---------------------")

	for idx, cmdList := range commandSequence {
		commandToRun := cmdList[0]
		for _, arg := range cmdList[1:] {
			commandToRun += " " + arg
		}
		fmt.Printf("running step %d/%d: %s\n", idx+1, len(commandSequence), commandToRun)

		cmd := exec.Command(cmdList[0], cmdList[1:]...)

		var so saveOutput
		cmd.Stdout = &so
		cmd.Stderr = os.Stderr

		// clear the slice regardless of output
		cmdList = nil

		// should capture output of gcc command
		if err := cmd.Run(); err != nil {
			fmt.Println("encountered an error: ", err)
			fmt.Print(so.savedOutput)
			return err
		}
	}

	return nil
}

// consider doing this our own
func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}
