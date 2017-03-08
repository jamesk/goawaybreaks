package main

import (
	"go/token"
	"go/ast"

	"regexp"
	"os"
	"io/ioutil"
	"flag"
)
/*"io/ioutil"
"os"*/
var test = `import (
 "remote.com/other"
 "local"


 "remote.com/package"
 "local2"
)`

var whitespace = regexp.MustCompile(`\s*\n`)

func main() {
	write := flag.Bool("w", false, "write result to (source) file instead of stdout")
	args := flag.Args()

	var src []byte
	var err error


	if len(args) > 0 {
		src, err = ioutil.ReadFile(args[0])
		if err != nil {
			panic(err)
		}
	} else {
		src, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	}

	//src := []byte(test)

	out, err := joinImportGroups(src)

	if len(args) > 0 && *write {
		ioutil.WriteFile(args[0], out, 0)
	}

	_, err = os.Stdout.Write(out)
	if err != nil {
		panic(err)
	}
}

func joinImportGroups(src []byte) ([]byte, error) {
	opt := &Options{Fragment:true, Comments: true, TabIndent: true, TabWidth: 8}
	fileSet := token.NewFileSet()
	file, _, err := parse(fileSet, "<stdin>", src, opt)
	if err != nil {
		panic(err)
	}

	ends := findEndOfGroups(fileSet, file)
	srcString := string(src)
	newLinePositions := regexp.MustCompile(`\n`).FindAllIndex(src, -1)

	deleted := 0
	for _, end := range ends {
		endAdjusted := end - deleted

		for {
			index := newLinePositions[endAdjusted-1][1]
			nextIndex := newLinePositions[endAdjusted][1]

			nextLine := string(srcString[index:nextIndex])
			matchIndexes := whitespace.FindAllStringIndex(nextLine, -1)
			if len(matchIndexes) != 1 {
				break
			}
			if matchIndexes[0][0] != 0 || matchIndexes[0][1] != len(nextLine) {
				break
			}

			deleted++
			srcString = srcString[:index] + srcString[nextIndex:]
		}
	}

	return []byte(srcString), nil
}


func findEndOfGroups(fset *token.FileSet, f *ast.File) []int {
	lastOfGroup := make([]int, 0)

	for _, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Imports are always first.
			break
		}

		if !d.Lparen.IsValid() {
			// Not a block no empty lines to remove
			continue
		}

		// Identify import group breaks
		i := 0
		for j, s := range d.Specs {
			if j > i && fset.Position(s.Pos()).Line > 1+fset.Position(d.Specs[j-1].End()).Line {
				// j begins a new run. Add last line of group
				lastOfGroup = append(lastOfGroup, fset.Position(d.Specs[j-1].End()).Line,)
				i = j
			}
		}
	}

	return lastOfGroup
}