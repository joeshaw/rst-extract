// The rst-extract utility extracts reStructured Text (RST) from Go
// source comments tagged with the special token "+rst" in the first
// line.
//
// Usage: rst-extract <source dir> <output dir>
// where "source dir" is a directory containing Go source files ending
// in .go, and "output dir" is a directory to output .rst files.
//
// rst-extract will create one .rst file per package.  The source
// files are processed in a predictable order: (1) a file name
// matching the package name (for instance, "main.go" in a "main"
// package); (2) "doc.go"; and (3) lexicographic order.  Comments
// within a file are processed in the order they appear.  This
// predictable ordering allows you to add, for instance, a header to
// the output RST file by adding it to one of the special cases that
// are processed first.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"sort"
	"strings"
)

func stripExt(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[:i]
		}
	}
	return filename
}

// fileSorter sorts files passed in according to this heuristic:
// 1. file with name == package name
// 2. file with name "doc.go"
// 3. lexicographic order
type fileSorter struct {
	pkg   *ast.Package
	files []*ast.File
	names []string
}

func newFileSorter(pkg *ast.Package) *fileSorter {
	fs := &fileSorter{}
	fs.pkg = pkg

	for name, f := range pkg.Files {
		fs.files = append(fs.files, f)
		fs.names = append(fs.names, path.Base(stripExt(name)))
	}

	sort.Sort(fs)
	return fs
}

func (fs fileSorter) Len() int {
	return len(fs.files)
}

func (fs fileSorter) Swap(i, j int) {
	fs.files[i], fs.files[j] = fs.files[j], fs.files[i]
	fs.names[i], fs.names[j] = fs.names[j], fs.names[i]
}

func (fs fileSorter) Less(i, j int) bool {
	ni, nj := fs.names[i], fs.names[j]

	switch ni {
	case fs.pkg.Name:
		switch nj {
		case fs.pkg.Name:
			return false

		default:
			return true
		}

	case "doc":
		switch nj {
		case fs.pkg.Name, "doc":
			return false

		default:
			return true
		}

	default:
		switch nj {
		case fs.pkg.Name, "doc":
			return false

		default:
			return ni < nj
		}
	}
}

var srcDir, outDir string

func rstComment(cgrp *ast.CommentGroup) (string, bool) {
	s := cgrp.Text()
	parts := strings.SplitN(s, "\n", 2)
	if strings.TrimSpace(parts[0]) == "+rst" {
		return parts[1], true
	}
	return "", false
}

func parsePackage(pkg *ast.Package) []string {
	sorted := newFileSorter(pkg)

	var comments []string
	for _, f := range sorted.files {
		for _, c := range f.Comments {
			s, ok := rstComment(c)
			if ok {
				comments = append(comments, s)
			}
		}
	}

	return comments
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <source dir> <output dir>\n", os.Args[0])
		return
	}

	srcDir = os.Args[1]
	outDir = os.Args[2]

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, srcDir, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing files in %s: %s\n", srcDir, err)
		os.Exit(1)
	}

	if err := os.MkdirAll(outDir, os.FileMode(0755)); err != nil {
		fmt.Printf("Error creating %s: %s\n", outDir, err)
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		var outf *os.File
		filename := path.Join(outDir, pkg.Name+".rst")

		comments := parsePackage(pkg)
		if len(comments) > 0 {
			outf, err = os.Create(filename)
			if err != nil {
				fmt.Printf("Error creating %s: %s\n", filename, err)
				os.Exit(1)
			}

			for _, c := range comments {
				outf.WriteString(c)
				outf.WriteString("\n")
			}

			outf.Close()
			fmt.Printf("Wrote %s\n", filename)
		}
	}
}
