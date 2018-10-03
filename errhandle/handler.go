package errhandle

import (
	"fmt"
	help "goCtrl/help"
)

type ErrorHandler struct {
	helperTextPrinter help.HelpTextPrinter
}

func (handler ErrorHandler) HandleError() {
	err := recover()
	if err == nil {
		return
	}

	switch err {
	case InvalidCommand:
		fmt.Println(err)
		handler.printHelpText()
	default:
		fmt.Println(err)
	}
}

func (handler ErrorHandler) printHelpText() {
	handler.helperTextPrinter.Print()
}
