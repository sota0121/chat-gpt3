package application

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/exp/slog"
)

// extractCode extracts code from file
// This uses AST if a function name is provided
// Otherwise, this reads the whole file
func extractCode(file *os.File, funcName string) (string, error) {
	// Read the whole file
	if funcName == "" {
		return getFileContent(file)
	}

	// Parse file and extract the target function
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file.Name(), nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == funcName {
				// Get where the function starts and ends
				start := fset.Position(fn.Pos()).Line
				end := fset.Position(fn.End()).Line

				// Get the content of the file
				content, err := getFileContent(file)
				if err != nil {
					slog.Error("Error getting file content", err)
					return "", err
				}

				// Extract the target function
				lines := strings.Split(content, "\n")
				return strings.Join(lines[start-1:end], "\n"), nil
			}
		}
	}
	err = errors.New("function not found")
	return "", err
}

// getFileContent returns the content of the file
func getFileContent(file *os.File) (string, error) {
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	data := make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
