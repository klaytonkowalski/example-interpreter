package repl

////////////////////////////////////////////////////////////////////////////////
// DEPENDENCIES
////////////////////////////////////////////////////////////////////////////////

import (
	"bufio"
	"fmt"
	"io"

	"github.com/klaytonkowalski/example-interpreter/lexer"
	"github.com/klaytonkowalski/example-interpreter/parser"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES
////////////////////////////////////////////////////////////////////////////////

const prompt = ">> "

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS
////////////////////////////////////////////////////////////////////////////////

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, prompt)
		scan := scanner.Scan()
		if !scan {
			return
		}
		line := scanner.Text()
		lxr := lexer.New(line)
		prs := parser.New(lxr)
		program := prs.ParseProgram()
		if len(prs.Errors) > 0 {
			printParserErrors(out, prs.Errors)
			continue
		}
		io.WriteString(out, program.GetDebugString())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
