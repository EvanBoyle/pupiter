package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/evanboyle/pupiter/eval"
	"github.com/evanboyle/pupiter/session"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, sessionName string) {
	session, err := session.NewSession(sessionName)
	if err != nil {
		panic(err)
	}
	evaluator := eval.NewEvaluator(session)

	fmt.Printf("new session started: %q\n", sessionName)
	fmt.Printf("resume session later with: 'pupiter %s'\n", sessionName)

	scanner := bufio.NewScanner(in)

	var input bytes.Buffer
	showPrompt := true

	for {
		if showPrompt {
			fmt.Printf(PROMPT)
		}
		showPrompt = false
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		input.Write(scanner.Bytes())
		input.WriteString("\n")

		if strings.HasSuffix(strings.TrimSpace(scanner.Text()), ";") {
			resp, err := evaluator.Eval(input.String())
			if err != nil {
				fmt.Println("error occured, try again:")
				fmt.Println("")
				fmt.Println("")
				fmt.Println(err)
				fmt.Println()
			} else {
				fmt.Println(resp)
			}
			input.Reset()
			showPrompt = true
		}

	}
}
