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
		lexer_ := lexer.New(line)
		for token_ := lexer_.GetNextToken(); token_.Category != token.End; token_ = lexer_.GetNextToken() {
			fmt.Fprintf(out, "%+v\n", token_)
		}
	}
}
