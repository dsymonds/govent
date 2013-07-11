package home

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"

	"github.com/dsymonds/govent/govent"
)

func init() {
	http.HandleFunc("/", handleFront)
	http.HandleFunc("/cmd", handleCmd)
}

var pages = template.Must(template.New("").ParseGlob("*.html"))

func handleFront(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var b bytes.Buffer
	if err := pages.ExecuteTemplate(&b, "front.html", nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	io.Copy(w, &b)
}

// State represents the marshaled state stored in datastore.
type State struct {
	X []byte
}

// CmdReply represents the JSON reply from /cmd.
type CmdReply struct {
	Reply string `json:"reply,omitempty"`
	Error error  `json:"error,omitempty"`
}

func handleCmd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "should be a POST", http.StatusMethodNotAllowed)
		return
	}

	c := appengine.NewContext(r)
	cmd := r.FormValue("cmd")
	c.Debugf("cmd %s", cmd)
	reply, err := doCmd(c, cmd)
	c.Debugf("-> %q, %v", reply, err)
	b, err := json.Marshal(&CmdReply{
		Reply: reply,
		Error: err,
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

func doCmd(c appengine.Context, cmd string) (string, error) {
	stateKey := datastore.NewKey(c, "State", "state:"+user.Current(c).Email, 0, nil)

	// load state
	state := new(govent.State)
	var encState State
	err := datastore.Get(c, stateKey, &encState)
	if err == nil {
		err = state.Unmarshal(encState.X)
	}
	if err != nil && err != datastore.ErrNoSuchEntity {
		return "", err
	}

	reply := state.Execute(cmd)

	// save state
	if encState.X, err = state.Marshal(); err != nil {
		return "", err
	}
	if _, err = datastore.Put(c, stateKey, &encState); err != nil {
		return "", err
	}

	return reply, nil
}
