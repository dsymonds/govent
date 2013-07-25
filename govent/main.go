/*
Package govent ...
*/
package govent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

var frontPage = template.Must(template.New("").Parse(string(frontHTML))) // in front.html.go

// Interface defines the API for a govent game.
type Interface interface {
	// Title returns the name of the game.
	Title() string
	// Execute runs a command.
	Execute(cmd string) string
	// GameOver reports whether the game is over.
	GameOver() bool

	// state save and restore
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}

type State struct {
	iface Interface
}

func NewState(iface Interface) *State {
	return &State{iface: iface}
}

func (s *State) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		s.handleFront(w, r)
	case "/cmd":
		s.handleCmd(w, r)
	case "/dump-state":
		s.handleDumpState(w, r)
	case "/reset":
		s.handleReset(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *State) handleFront(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var b bytes.Buffer
	data := struct {
		Title string
		User  *user.User
	}{
		Title: s.iface.Title(),
		User:  user.Current(c),
	}
	if err := frontPage.Execute(&b, &data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	io.Copy(w, &b)
}

// cmdReply represents the JSON reply from /cmd.
type cmdReply struct {
	Reply    string `json:"reply,omitempty"`
	Error    error  `json:"error,omitempty"`
	GameOver bool   `json:"gameOver,omitempty"`
}

// encState represents the marshaled state stored in datastore.
type encState struct {
	X []byte
}

func (s *State) stateKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "State", "state:"+user.Current(c).Email, 0, nil)
}

func (s *State) handleCmd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "should be a POST", http.StatusMethodNotAllowed)
		return
	}

	c := appengine.NewContext(r)
	cmd := r.FormValue("cmd")
	c.Debugf("cmd %s", cmd)
	reply, err := s.doCmd(c, cmd)
	c.Debugf("-> %q, %v", reply, err)
	b, err := json.Marshal(&cmdReply{
		Reply:    reply,
		Error:    err,
		GameOver: s.iface.GameOver(),
	})
	if err != nil {
		// should not happen
		c.Errorf("json.Marshal: %v", err)
		http.Error(w, "oopsie, internal error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func (s *State) doCmd(c appengine.Context, cmd string) (string, error) {
	// load state
	var encState encState
	err := datastore.Get(c, s.stateKey(c), &encState)
	if err != nil && err != datastore.ErrNoSuchEntity {
		return "", err
	}
	if err == nil {
		if err := s.iface.Unmarshal(encState.X); err != nil {
			// press on anyway
			c.Errorf("Bad state, forgetting about it: %v", err)
		}
	}

	reply := s.iface.Execute(cmd)

	// save state
	if encState.X, err = s.iface.Marshal(); err != nil {
		return "", err
	}
	if _, err = datastore.Put(c, s.stateKey(c), &encState); err != nil {
		return "", err
	}

	return reply, nil
}

func (s *State) handleDumpState(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if !user.IsAdmin(c) {
		http.Error(w, "admins only!", http.StatusForbidden)
		return
	}

	var encState encState
	err := datastore.Get(c, s.stateKey(c), &encState)
	if err != nil && err != datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(encState.X)))
	w.Write(encState.X)
}

func (s *State) handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "should be a POST", http.StatusMethodNotAllowed)
		return
	}

	c := appengine.NewContext(r)
	stateKey := s.stateKey(c)
	if err := datastore.Delete(c, stateKey); err != nil && err != datastore.ErrNoSuchEntity {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "OK; reset state for %v\n", stateKey)
}
