package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	ignoreCase      bool
	inplace         bool
	ignoreNewLines  bool
	skipFirstChars  uint
	skipLastChars   uint
	skipFirstFields uint
	skipLastFields  uint
	oneField        uint
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Unordered uniq(1).\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %v [OPTIONS] [FILE]...\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nAll indexes begin at 1.")
	}
	flag.BoolVar(&ignoreCase, "i", false, "ignore case during comparison")
	flag.BoolVar(&ignoreNewLines, "n", false, "ignore duplicated newlines")
	flag.BoolVar(&inplace, "w", false, "write to files instead of stdout")
	flag.UintVar(&oneField, "f", 0, "use only field `N` for comparison")
	flag.UintVar(&skipFirstChars, "cf", 0, "skip comparing first `N` characters")
	flag.UintVar(&skipLastChars, "cl", 0, "skip comparing last `N` characters")
	flag.UintVar(&skipFirstFields, "ff", 0, "skip comparing first `N` fields")
	flag.UintVar(&skipLastFields, "fl", 0, "skip comparing last `N` fields")
	flag.Parse()

	var operations []func(string) string
	if ignoreCase {
		operations = append(operations, ignoreCaseOp)
	}
	if oneField != 0 {
		operations = append(operations, oneFieldOp)
	}
	if skipFirstChars != 0 {
		operations = append(operations, skipFirstCharsOp)
	}
	if skipLastChars != 0 {
		operations = append(operations, skipLastCharsOp)
	}
	if skipFirstFields != 0 {
		operations = append(operations, skipFirstFieldsOp)
	}
	if skipLastFields != 0 {
		operations = append(operations, skipLastFieldsOp)
	}
	if len(operations) == 0 {
		operations = append(operations, func(s string) string {
			return s
		})
	}

	if flag.NArg() == 0 {
		if err := uuniq(os.Stdin, os.Stdout, operations); err != nil {
			log.Print(err)
		}
	}

	var output io.StringWriter = os.Stdout
	for _, arg := range flag.Args() {
		input, err := os.Open(arg)
		if err != nil {
			log.Fatal(err)
		}
		if inplace {
			buf := make([]byte, 0, 1024)
			outputBuf := bytes.NewBuffer(buf)
			output = outputBuf
			defer overwriteFile(arg, outputBuf)
		}
		if err = uuniq(input, output, operations); err != nil {
			log.Print(err)
			return
		}
		if err = input.Close(); err != nil {
			log.Print(err)
			return
		}
	}
}

func overwriteFile(name string, data *bytes.Buffer) {
	realOutput, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = realOutput.Write(data.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err = realOutput.Close(); err != nil {
		log.Fatal(err)
	}
}

// "inspired" by https://github.com/ptrcnull/uuniq
func uuniq(input io.Reader, output io.StringWriter, operations []func(string) string) error {
	history := make(map[string]bool)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if ignoreNewLines && line == "" {
			_, err := output.WriteString("\n")
			if err != nil {
				return err
			}
			continue
		}
		processedLine := processLine(line, operations)
		if !history[processedLine] {
			_, err := output.WriteString(line + "\n")
			if err != nil {
				return err
			}
			history[processedLine] = true
		}
	}
	return nil
}

func processLine(line string, operations []func(string) string) string {
	for i := range operations {
		line = operations[i](line)
	}
	return line
}

func ignoreCaseOp(s string) string {
	return strings.ToLower(s)
}

func skipFirstCharsOp(s string) string {
	start := int(skipFirstChars)
	l := len(s)
	if start > l {
		return s
	}
	return s[start:]
}

func skipLastCharsOp(s string) string {
	end := int(skipLastChars)
	l := len(s)
	if end > l {
		return s
	}
	return s[:len(s)-end]
}

func skipFirstFieldsOp(s string) string {
	start := int(skipFirstFields)
	f := strings.Fields(s)
	l := len(f)
	if start > l {
		return s
	}
	return strings.Join(f[start:], " ")
}

func skipLastFieldsOp(s string) string {
	end := int(skipLastFields)
	f := strings.Fields(s)
	l := len(f)
	if end > l {
		return s
	}
	return strings.Join(f[:l-end], " ")
}

func oneFieldOp(s string) string {
	index := oneField
	f := strings.Fields(s)
	if oneField > uint(len(f)) {
		return s
	}
	return strings.Join(f[index-1:index], " ")
}
