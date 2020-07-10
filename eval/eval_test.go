package eval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	execInput := `
	var x = aws.s3.Bucket("abc");
	`
	execStatement := parse(execInput)
	assert.Equal(t, `var x = aws.s3.Bucket("abc");`, execStatement.Text)
	assert.Equal(t, Exec, execStatement.Type)
	assert.Equal(t, "x", execStatement.VarName)

	refInput := `x;`
	refStatement := parse(refInput)
	assert.Equal(t, `x;`, refStatement.Text)
	assert.Equal(t, Ref, refStatement.Type)
	assert.Equal(t, "x", refStatement.VarName)
}
