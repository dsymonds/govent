/*
Package govent ...
*/
package govent

import (
	"encoding/json"
	"fmt"
	"strings"
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

	// Hack night!
	switch {
	case cmd == "LOOK":
		return "You are in a room full of gophers."
	case strings.HasPrefix(cmd, "LOOK AT "):
		return "He glares at you in return."
	case cmd == "GET LAMP":
		return "Okay, you have a lamp, smartass. Now what?"
	}

	return fmt.Sprintf("You said %q. That was your %dth command.", cmd, s.Count)
}
