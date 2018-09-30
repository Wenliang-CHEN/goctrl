package errhandle

import (
	"fmt"
	help "goCtrl/help"
)

func HandleError() {
	err := recover()
	if err == nil {
		return
	}

	switch err {
	case InvalidCommand:
		fmt.Println(err)
		printHelpText()
	default:
		fmt.Println(err)
	}
}

// TODO: Add help text
func printHelpText() {
	help.PrintHelpText()
}
