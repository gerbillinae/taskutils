package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filename := os.Args[1]

	is_env_file := filepath.Ext(filename) == ".env"

	flag.Parse()

	log.Println(filename)
	// Specify the template file name

	// Open the template file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)

	total_sub_failures := 0
	for scanner.Scan() {
		line := scanner.Text()

		// Replace occurrences of environment variables
		// Environment variables are expected in the format ${VAR_NAME}
		line, sub_failures := replaceEnvVariables(line)
		total_sub_failures += sub_failures
		fmt.Println(line)

		if is_env_file {
			trimmed := strings.TrimSpace(line)
			if !(trimmed == "" || strings.HasPrefix(trimmed, "#")) {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					os.Setenv(key, value)
				}
			}
		}
	}

	os.Exit(total_sub_failures)
}

func replaceEnvVariables(text string) (string, int) {
	// Look for occurrences of ${VAR_NAME}
	start := "${"
	end := "}"

	sub_failures := 0
	for {
		startIdx := strings.Index(text, start)
		if startIdx == -1 {
			break
		}

		endIdx := strings.Index(text[startIdx:], end)
		if endIdx == -1 {
			break
		}

		// Extract the variable name
		varName := text[startIdx+len(start) : startIdx+endIdx]

		// Get the value of the environment variable
		varValue := os.Getenv(varName)

		if varValue != "" {
			// Replace the variable occurrence in the text
			text = strings.Replace(text, start+varName+end, varValue, 1)
		} else {
			sub_failures++
		}
	}

	return text, sub_failures
}
