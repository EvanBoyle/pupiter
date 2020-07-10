package eval

import (
	"errors"
	"fmt"
	"strings"

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

	// todo - addVar getVar
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

	return "", errors.New("unable to parse statement")
}

func (e *evaluator) addVar(statement Statement) (string, error) {
	err := write.WriteFiles(statement.VarName, statement.Text, e.session)
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

type StatementType = int

const (
	Exec StatementType = iota
	Ref
	List
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
