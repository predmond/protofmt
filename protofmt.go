package main

import (
	"fmt"
	"os"
	"strings"
	"text/scanner"
)

func main() {
	var s scanner.Scanner
	s.Init(os.Stdin)
	s.Filename = "stdin"
	s.Mode = scanner.ScanIdents |
		scanner.ScanFloats |
		scanner.ScanChars |
		scanner.ScanStrings |
		scanner.ScanRawStrings |
		scanner.ScanComments

	indentLevel := 0
	indentString := "\t"
	indent := func() {
		fmt.Print(strings.Repeat(indentString, indentLevel))
	}

	newLines := 0
	prevTok := 0
	prevLine := "syntax"

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tt := s.TokenText()

		if newLines > 0 {
			switch tt {
			case "syntax", "package", "message", "import":
				if prevLine != tt {
					newLines++
				}
			}
			prevLine = tt

			fmt.Print(strings.Repeat("\n", newLines))
			indent()
			newLines = 0
		} else {
			space := false
			switch prevTok {
			case scanner.Ident:
				switch tok {
				case '.', ')', ']', ',', ';':
				default:
					space = true
				}
			case scanner.Int, '=', ')', ',':
				space = true
			}
			if space {
				fmt.Print(" ")
			}
		}

		switch tt {
		case ";":
			fmt.Print(tt)
			newLines = 1

		case "{":
			fmt.Print(tt)
			indentLevel++
			newLines = 1

		case "}":
			indentLevel--
			fmt.Print(tt)
			newLines = 1

		default:
			fmt.Print(tt)
			//fmt.Printf("[%d]", tok)
		}

		if tok == scanner.Comment {
			newLines = 1
		}

		prevTok = int(tok)
	}

}
