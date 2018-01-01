package main

import (
	"testing"
)

var formatTests = []struct {
	in  string
	out string
}{
	{
		"message{}",
		`message {
}`,
	},
	{
		"message{message{}}",
		`message {
	message {
	}
}`,
	},
	{
		`syntax="proto2";`,
		`syntax = "proto2";`,
	},
	{
		`repeated string foo=1;`,
		`repeated string foo = 1;`,
	},
	{
		`optional string foo=1[(gogoproto.nullable)=false];`,
		`optional string foo = 1 [(gogoproto.nullable) = false];`,
	},
}

func TestFormat(t *testing.T) {
	for _, tt := range formatTests {
		out := formatString(tt.in)
		if out != tt.out {
			t.Errorf("wrong format:\n%s\n%s", out, tt.out)
		}
	}
}
