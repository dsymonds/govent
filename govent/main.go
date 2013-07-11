/*
Package govent ...
*/
package govent

import (
	"fmt"
)

type Interface interface {
}

type State struct {
}

func (s *State) Marshal() ([]byte, error) {
	// TODO
	return []byte{}, nil
}

func (s *State) Unmarshal(data []byte) error {
	// TODO
	return nil
}

func (s *State) Execute(cmd string) string {
	return fmt.Sprintf("You said %q", cmd)
}
