package action

import (
	"fmt"

	"github.com/atotto/clipboard"
)

type CopyPaster struct {
}

func (cp CopyPaster) Execute(message string) {

	clipboard.WriteAll(message)
	fmt.Println("echo ** command pasted to clipboard **")
}
