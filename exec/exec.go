package exec

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/evanboyle/pupiter/session"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

// TODO pulumi up via auto

// executor

func Execute(varName string, session session.Session) (string, error) {
	p := auto.Project{
		Name:       session.Name(),
		SourcePath: filepath.Join(session.Dir(), varName),
	}
	s := &auto.Stack{
		Name:    varName,
		Project: p,
	}
	res, err := s.Up()
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute statement: %s, %s", res.StdErr, res.StdOut)
	}
	b, err := json.MarshalIndent(res.Outputs[varName], "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Get(varName string, session session.Session) (string, string, error) {
	p := auto.Project{
		Name:       session.Name(),
		SourcePath: filepath.Join(session.Dir(), varName),
	}
	s := &auto.Stack{
		Name:    varName,
		Project: p,
	}
	outs, secs, err := s.GetOutputs()
	oStr, err := json.MarshalIndent(outs, "", "  ")
	if err != nil {
		return "", "", err
	}
	sStr, err := json.MarshalIndent(secs, "", "  ")
	if err != nil {
		return "", "", err
	}
	return string(oStr), string(sStr), nil
}

func List(session session.Session) ([]string, error) {
	files, err := ioutil.ReadDir(session.Dir())
	if err != nil {
		return nil, err
	}

	var res []string
	for _, f := range files {
		res = append(res, f.Name())
	}
	return res, nil
}
