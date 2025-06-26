package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type InterpreterRequest struct {
	Code string `json:"code"`
}

type InterpreterResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func handleInterpret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	var req InterpreterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	output, errMsg := runLoxCode(req.Code)

	response := InterpreterResponse{
		Output: output,
		Error:  errMsg,
	}

	json.NewEncoder(w).Encode(response)
}

func runLoxCode(source string) (string, string) {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	var output bytes.Buffer
	var errorMsg string

	done := make(chan bool)
	go func() {
		io.Copy(&output, r)
		done <- true
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				errorMsg = fmt.Sprintf("Runtime error: %v", r)
			}
		}()

		scanner := NewScanner(source)
		tokens, scanErrors := scanner.ScanTokens()
		if len(scanErrors) > 0 {
			var errStrs []string
			for _, err := range scanErrors {
				errStrs = append(errStrs, err.Error())
			}
			errorMsg = strings.Join(errStrs, "\n")
			return
		}

		parser := NewParser(tokens)
		statements, err := parser.parse()
		if err != nil {
			errorMsg = err.Error()
			return
		}

		interpreter := NewInterpreter()
		resolver := NewResolver(interpreter)

		func() {
			defer func() {
				if r := recover(); r != nil {
					errorMsg = fmt.Sprintf("Resolution error: %v", r)
				}
			}()
			resolver.Resolve(statements)
		}()

		if errorMsg != "" {
			return
		}

		interpreter.Interpret(statements)
	}()

	w.Close()
	os.Stdout = oldStdout

	<-done

	return output.String(), errorMsg

}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "web/index.html")
	} else {
		http.FileServer(http.Dir("web/")).ServeHTTP(w, r)
	}
}
