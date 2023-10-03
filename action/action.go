package action

import "fmt"

type Action interface {
	Execute(message string)
}

type Printer struct {
}

func (printer Printer) Execute(message string) {
	fmt.Println(message)
}
