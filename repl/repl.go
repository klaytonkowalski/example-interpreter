package repl

import (
	"bufio"
	"example-interpreter/lexer"
	"example-interpreter/token"
	"fmt"
	"io"
)

const prefix = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, prefix)
		scan := scanner.Scan()
		if !scan {
			return
		}
		line := scanner.Text()
		lxr := lexer.New(line)
		for tok := lxr.GetNextToken(); tok.Category != token.End; tok = lxr.GetNextToken() {

		}
	}
}
