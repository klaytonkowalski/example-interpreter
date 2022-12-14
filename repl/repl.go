package repl

import (
	"bufio"
	"example-interpreter/lexer"
	"example-interpreter/token"
	"fmt"
	"io"
)

// Prefix to display before any user input.
const prefix = ">> "

// REPL entry point.
// Launches the REPL and continuously waits, parses, and outputs.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, prefix)
		scan := scanner.Scan()
		if !scan {
			return
		}
		line := scanner.Text()
		lexer_ := lexer.New(line)
		for token_ := lexer_.NextToken(); token_.Category != token.End; token_ = lexer_.NextToken() {
			fmt.Fprintf(out, "%+v\n", token_)
		}
	}
}
