package logger

import "fmt"

type uc struct {
}

func New() UC {
	return &uc{}
}

func (u *uc) Log(level int64, message string) {
	fmt.Printf("%d: %s\n", level, message)
}
