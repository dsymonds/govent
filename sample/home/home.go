package home

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dsymonds/govent/govent"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := govent.NewState(&game{})
		s.ServeHTTP(w, r)
	})
}

type game struct{
	Count int
}

func (g *game) Marshal() ([]byte, error) {
	return json.Marshal(g)
}

func (g *game) Unmarshal(data []byte) error {
	return json.Unmarshal(data, g)
}

func (g *game) Execute(cmd string) string {
	g.Count++

	// Hack night!
	switch {
	case cmd == "LOOK":
		return "You are in a room full of gophers."
	case strings.HasPrefix(cmd, "LOOK AT "):
		return "He glares at you in return."
	case cmd == "GET LAMP":
		return "Okay, you have a lamp, smartass. Now what?"
	}

	return fmt.Sprintf("You said %q. That was your %dth command.", cmd, g.Count)
}
