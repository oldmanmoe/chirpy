package main

import (
	"reflect"
	"testing"
)

// run test with: go test -coverprofile=c.out
func TestCleanBadWords(t *testing.T) {
	tests := map[string]struct {
		input	string
		want 	string
	}{
		"Regular Chirp":		{input: "This is just a regular string", want: "This is just a regular string"},
		"Kerffufle":			{input: "This is a kerfuffle opinion", want: "This is a **** opinion"},
		"Kurffufle Caps":		{input: "This is a KERFUFFLE but in caps", want: "This is a **** but in caps"},
		"Multiple Bad Words": 	{input: "This sharbert is a complete kerfuffle", want:  "This **** is a complete ****"},
		"Repeated Bad Words": 	{input: "sharbert sharbert sharbert", want:  "**** **** ****"},
		"Empty Input": 			{input: "", want:  ""},
		"Testing Sharbert":      {input: "Behold the power of sharbert", want: "Behold the power of ****"},
        "Testing Fornax Caps":  {input: "FORNAX is a constellation", want: "**** is a constellation"},
        "Punctuation Attached": {input: "Go to fornax!", want: "Go to fornax!"},
		"Random":				{input: "I think Sharbert is rude", want: "I think **** is rude"},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T){
			got := cleanBadWords(testCase.input)
			if !reflect.DeepEqual(testCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", name, testCase.want, got)
			}
		})
	}
}