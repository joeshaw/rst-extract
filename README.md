# rst-extract #

rst-extract is a tool for extracting specially tagged comments in Go
source code containing
[reStructured Text](http://docutils.sourceforge.net/rst.html).  Both
grouped line comments (`//`) and block comments (`/* */`) are
supported.

## How do I install this? ##

    go get github.com/joeshaw/rst-extract

## How do I use this? ##

Simply tag comments in your Go source code with `+rst` as the first
line of your comment.  For example,

```go
package main

// +rst
// API Documentation
// =================
//
// (FIXME: add documentation)
```

Then run the `rst-extract` command, providing it with a directory of
Go source files and an output directory:

    rst-extract ./example /tmp/example

rst-extract will output one `.rst` file per Go package.

Source files are processed in a predictable order:

1. A file name matching the package name (for instance, `main.go` in a
`main` package)
2. `doc.go`
3. lexicographic order

Comments within a file are processed in the order they appear.

This predictable ordering allows you to add, for instance, a header to
the output RST file by adding it to one of the special cases that
are processed first.

## Why would I want this? ##

I know what you're thinking.  You're thinking these two things:

1. [`godoc`](http://godoc.org) is awesome, what is the point of this?
2. RST sucks, Markdown is much better

And yes, you are right.  However, in past Python projects I've used
[Sphinx](http://sphinx-doc.org/) and its excellent
[httpdomain](https://pythonhosted.org/sphinxcontrib-httpdomain/)
extension for documenting HTTP APIs within my code.  Godoc is not
well-suited for this task, and while I prefer Markdown in general,
RST is what Sphinx deals in, and its additional structure works
well in this case.

In an ideal world I would have gone straight from Go comments to
HTML for my API docs, and at some future point I might do that.
In the meantime, however, this gets the job done for me, despite
the messy Python dependency for Sphinx.

## How does it work? ##

I hate writing parsers, and I am very thankful that I don't have to
write one for this.  This builds upon the excellent code provided in
the Go standard library for parsing Go code, namely
[`go/parser`](http://golang.org/pkg/go/parser) and
[`go/ast`](http://golang.org/pkg/go/ast).
