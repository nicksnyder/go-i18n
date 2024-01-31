package main

import "testing"

func TestMarshal(t *testing.T) {
	actual, err := marshal(map[string]string{
		"&<key>": "&<val>",
	}, "json")

	if err != nil {
		t.Fatal(err)
	}

	expected := `{
  "&<key>": "&<val>"
}
`
	if a := string(actual); a != expected {
		t.Fatalf("\nexpected:\n%s\n\ngot\n%s", expected, a)
	}
}
