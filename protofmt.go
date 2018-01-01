package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

type Formatter struct {
	scanner.Scanner

	Indent string

	indentLevel int
	indented    bool
	prevLine    string

	out io.Writer
}

func NewFormatter(in io.Reader, out io.Writer) *Formatter {
	f := &Formatter{
		Indent: "\t",
		out:    out,
	}
	f.Init(in)
	f.Mode = scanner.ScanIdents |
		scanner.ScanFloats |
		scanner.ScanChars |
		scanner.ScanStrings |
		scanner.ScanRawStrings |
		scanner.ScanComments
	return f
}

func (f *Formatter) print(a ...interface{}) {
	fmt.Fprint(f.out, a...)
}

func (f *Formatter) indent() {
	if !f.indented {
		f.print(strings.Repeat(f.Indent, f.indentLevel))
	}
	f.indented = true
}

func (f *Formatter) newLine(tt string, newLines int) {
	switch tt {
	case "syntax", "package", "message", "import":
		if f.prevLine != tt {
			newLines++
		}
	}
	f.print(strings.Repeat("\n", newLines))
	f.prevLine = tt
	f.indented = false
}

func (f *Formatter) space(tok, prevTok rune) {
	switch tok {

	case '{', '[', '=', scanner.String, scanner.Int:
		break

	case scanner.Ident:
		switch prevTok {
		case '(', '.':
			return
		default:
			break
		}

	default:
		return
	}
	f.print(" ")
}

func (f *Formatter) Format() {
	prevTok := rune(0)
	f.prevLine = "syntax"

	newLines := 0
	for tok := f.Scan(); tok != scanner.EOF; tok = f.Scan() {
		tt := f.TokenText()

		if newLines > 0 {
			f.newLine(tt, newLines)
			newLines = 0
		}

		if f.indented {
			f.space(tok, prevTok)
		}

		switch tt {
		case ";":
			f.print(tt)
			newLines = 1

		case "{":
			f.print(tt)
			f.indentLevel++
			newLines = 1

		case "}":
			f.indentLevel--
			f.indent()
			f.print(tt)
			newLines = 1

		default:
			f.indent()
			f.print(tt)
		}

		if tok == scanner.Comment {
			newLines = 1
		}

		prevTok = tok
	}

}

func formatString(in string) string {
	var out bytes.Buffer
	f := NewFormatter(strings.NewReader(in), &out)
	f.Format()
	return out.String()
}

func main() {
	f := NewFormatter(os.Stdin, os.Stdout)
	f.Format()
}
