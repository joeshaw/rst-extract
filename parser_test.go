package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestStripExt(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"foo.bar", "foo"},
		{"foo.bar.baz", "foo.bar"},
		{"foo", "foo"},
		{".asdf", ""},
	}

	for i, test := range tests {
		s := stripExt(test.input)
		if s != test.expected {
			t.Fatalf("Iteration %d: Expected %s, got %s", i, test.expected, s)
		}
	}
}

func TestSorter(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "test", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("Expected len(pkgs) == 1, got %d", len(pkgs))
	}

	expected := []string{"test", "doc", "a", "z"}

	sorter := newFileSorter(pkgs["test"])

	if len(sorter.names) != len(expected) {
		t.Fatalf("Expected %#v, got %#v", expected, sorter.names)
	}

	for i := 0; i < len(sorter.names); i++ {
		if sorter.names[i] != expected[i] {
			t.Fatalf("Expected %s, got %s", expected, sorter.names)
		}
	}
}

func TestRSTComment(t *testing.T) {
	tests := []struct {
		comment  string
		expected bool
	}{
		{"//+rst", true},
		{"// +rst", true},
		{"//           +rst\n", true},
		{"/*+rst*/", true},
		{"/* +rst */", true},
		{"/*           +rst */", true},
		{"/*\n+rst */", true},
		{"/* +rst\nfoo */", true},
		{"/* foo\n+rst */", false},
		{"// rst", false},
	}

	for i, test := range tests {
		var c ast.Comment
		c.Text = test.comment

		cgrp := &ast.CommentGroup{List: []*ast.Comment{&c}}
		_, ok := rstComment(cgrp)
		if ok != test.expected {
			t.Fatalf("Iteration %d: Expected %t, got %t: %q", i, test.expected, ok, test.comment)
		}
	}
}

func TestParsePackage(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "test", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("Expected len(pkgs) == 1, got %d", len(pkgs))
	}

	comments := parsePackage(pkgs["test"])

	expected := []string{
		"A comment inside test.go\n",
		"More restructured text, in doc.go\n",
		"Here's a comment in a.go\n",
		"An interesting\nmulti-line\ncomment inside\nz.go\n",
	}

	if len(comments) != len(expected) {
		t.Fatalf("Expected %#v, got %#v", expected, comments)
	}

	for i := 0; i < len(comments); i++ {
		if comments[i] != expected[i] {
			t.Fatalf("Expected %s, got %s", expected, comments)
		}
	}
}
