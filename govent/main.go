/*
Package govent ...
*/
package govent

import (
	"encoding/json"
	"fmt"
)

type Interface interface {
}

type State struct {
	Count int // number of commands
}

func (s *State) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *State) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *State) Execute(cmd string) string {
	s.Count++
	return fmt.Sprintf("You said %q. That was your %dth command.", cmd, s.Count)
}
