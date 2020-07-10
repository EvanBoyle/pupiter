package exec

import (
	"encoding/json"
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
