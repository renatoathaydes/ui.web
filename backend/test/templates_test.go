package test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	uiweb "ui.web/server/src"
)

const expectedOutput = `import * as hello from './hello.js';
import * as foo from './foo.js';

async function evalWith(cmd, me, value) {
    // the eval'd code can see the variables:
	//   * 'me' - the component which invoked the cmd
	//   * 'value' - the result of a previous command, usually run in the backend
	//   * all frontend JS modules by name.
    return await eval(cmd);
}

window.evalWith = evalWith;
`

func TestEvalJsTemplate(t *testing.T) {
	mods := []string{"hello.mjs", "foo.mts"}
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	err := uiweb.WriteEvalJsWith(writer, mods, true)

	require.Nil(t, err)
	err = writer.Flush()
	require.Nil(t, err)
	out := buffer.String()

	require.Equal(t, expectedOutput, out)
}
