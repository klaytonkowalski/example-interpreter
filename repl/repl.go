package repl

import (
	"bufio"
	"example-go-interpreter/lexer"
	"example-go-interpreter/token"
	"fmt"
	"io"
)

const Prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, Prompt)
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
