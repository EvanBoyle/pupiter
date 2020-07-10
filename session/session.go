package session

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Session interface {
	isSession()
	Name() string
	Dir() string
}

func NewSession(name string) (Session, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get users homedir")
	}
	dir := filepath.Join(home, ".pulumi", "pupiter", name)
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create pupiter session dir")
		}
	}
	return &session{
		name: name,
		dir:  dir,
	}, nil
}

type session struct {
	name string
	dir  string
}

func (s *session) Name() string { return s.name }
func (s *session) Dir() string  { return s.dir }
func (s *session) isSession()   {}
