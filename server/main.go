package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/evanboyle/pupiter/exec"
	"github.com/evanboyle/pupiter/session"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/pkg/errors"
)

func Serve() {
	// TODO need to serve this relative from a command

	//buildHandler := http.StripPrefix("/", http.FileServer(http.Dir(filepath.Join(".", "notebook", "build"))))
	buildHandler := http.StripPrefix("/", http.FileServer(pkger.Dir("/notebook/build")))
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(buildHandler)
	r.HandleFunc("/sessions", sessionsHandler)
	r.HandleFunc("/session/{session}", sessionHandler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":1337", nil))
}

// TODO write a handler to list all sessions
func sessionHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal("")
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to marshal stack outputs").Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(js))
}

func sessionsHandler(w http.ResponseWriter, r *http.Request) {
	sName := mux.Vars(r)["session"]
	s, err := session.NewSession(sName)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to initialize session").Error(), http.StatusInternalServerError)
		return
	}
	vars, err := exec.List(s)
	if err != nil {
		http.Error(w, errors.Wrap(err, "unable to list vars for session.").Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"vars": vars,
	}

	v, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to marshal session variables").Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(v))
}
