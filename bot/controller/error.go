package controller

import "fmt"

type ErrController struct {
	Err     error
	Message string
}

func (b *ErrController) Error() string {
	return fmt.Sprintf("%s: %v", b.Message, b.Err)
}

func NewErrController(Message string, Err error) *ErrController {
	return &ErrController{
		Message: Message,
		Err:     Err,
	}
}
