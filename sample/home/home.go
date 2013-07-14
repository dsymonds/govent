package home

import (
	"encoding/json"
	"net/http"

	"github.com/dsymonds/govent/govent"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := govent.NewState(&game{})
		s.ServeHTTP(w, r)
	})
}

type game struct {
	Count int

	Location int // 0 == inside house, 1 == outside
}

func (g *game) Title() string               { return "Home Quest" }
func (g *game) Marshal() ([]byte, error)    { return json.Marshal(g) }
func (g *game) Unmarshal(data []byte) error { return json.Unmarshal(data, g) }

func (g *game) Execute(cmd string) string {
	g.Count++

	const (
		insideText  = "You are standing in a house consisting of a single room."
		outsideText = "You are standing outside. It is a beautiful sunny day."
	)

	if cmd == "" || cmd == "LOOK" {
		switch g.Location {
		case 0:
			return insideText
		case 1:
			return outsideText
		}
	}

	switch g.Location {
	case 0:
		switch cmd {
		case "GO OUTSIDE":
			g.Location = 1
			return outsideText
		case "GO INSIDE":
			return "You're already inside, silly!"
		}
	case 1:
		switch cmd {
		case "GO INSIDE":
			g.Location = 0
			return insideText
		case "GO OUTSIDE":
			return "Outside? That's where you are now!"
		}
	}

	return "I don't understand."
}
