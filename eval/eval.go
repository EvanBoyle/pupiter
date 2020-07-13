package eval

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bradfitz/slice"

	"github.com/evanboyle/pupiter/exec"
	"github.com/evanboyle/pupiter/session"
	"github.com/evanboyle/pupiter/write"
)

type Evaluator interface {
	isEvaluator()
	Eval(input string) (string, error)
	addVar(statement Statement) (string, error)
	getVar(statement Statement) (string, string, error)
	getVars() ([]string, error)
}

type evaluator struct {
	session session.Session
}

func (e *evaluator) isEvaluator() {}

func NewEvaluator(session session.Session) Evaluator {
	return &evaluator{
		session,
	}
}

// Eval takes the given input and executes it
// with some hand waving in between
func (e *evaluator) Eval(input string) (string, error) {
	statement := parse(input)
	switch statement.Type {
	case Exec:
		out, err := e.addVar(statement)
		if err != nil {
			return "", err
		}
		return out, nil
	case List:
		vars, err := e.getVars()
		if err != nil {
			return "", err
		}
		return strings.Join(vars, "\n"), nil
	case Eject:
		return e.eject()
	case Ref:
		fmt.Printf("retrieving var: %s\n", input)
		outs, secs, err := e.getVar(statement)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("outputs: %s\nsecret outputs: %s\n", outs, secs), nil
	default:
		return "", errors.New("unable to parse statement")
	}
}

func (e *evaluator) addVar(statement Statement) (string, error) {
	vars, err := e.getVars()
	if err != nil {
		return "", err
	}

	err = write.WriteFiles(statement.VarName, statement.Text, e.session, vars)
	if err != nil {
		return "", err
	}
	out, err := exec.Execute(statement.VarName, e.session)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (s *evaluator) getVar(statement Statement) (string, string, error) {
	return exec.Get(statement.VarName, s.session)
}

func (s *evaluator) getVars() ([]string, error) {
	return exec.List(s.session)
}

func (s *evaluator) eject() (string, error) {
	files, err := ioutil.ReadDir(s.session.Dir())
	if err != nil {
		return "", err
	}
	slice.Sort(files[:], func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	program := `
const pulumi = require("@pulumi/pulumi");
const aws = require("@pulumi/aws");

%s
`

	var codes []string
	for _, v := range files {
		v.Name()

		fName := filepath.Join(s.session.Dir(), v.Name(), "index.js")
		f, err := os.Open(fName)
		if err != nil {
			return "", nil
		}
		bytes := make([]byte, 10000000)
		f.Read(bytes)
		contents := string(bytes)
		code := strings.Split(contents, "// ---------------------------------------------")[1]
		codes = append(codes, code)
	}

	return fmt.Sprintf(program, strings.Join(codes, "\n")), nil
}

type StatementType = int

const (
	Exec StatementType = iota
	Ref
	List
	Eject
)

type Statement struct {
	VarName string
	Type    StatementType
	Text    string
}

func parse(input string) Statement {
	input = strings.TrimSpace(input)
	// TODO: not so robust...
	if input == "ls();" {
		return Statement{
			Type: List,
		}
	}
	// TODO: not so robust...
	if input == "eject();" {
		return Statement{
			Type: Eject,
		}
	}
	if !strings.Contains(input, "=") {
		return Statement{
			VarName: strings.Trim(input, ";"),
			Type:    Ref,
			Text:    input,
		}
	}
	// `var foo = ...;`
	parts := strings.Split(input, "=")
	varName := strings.Split(strings.TrimSpace(parts[0]), " ")[1]
	return Statement{
		VarName: varName,
		Type:    Exec,
		Text:    input,
	}
}
