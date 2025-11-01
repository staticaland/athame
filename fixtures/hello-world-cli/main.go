package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Define command-line flags
	name := flag.String("name", "World", "Name to greet")
	uppercase := flag.Bool("uppercase", false, "Print greeting in uppercase")
	repeat := flag.Int("repeat", 1, "Number of times to repeat the greeting")

	flag.Parse()

	// Build the greeting message
	greeting := fmt.Sprintf("Hello, %s!", *name)

	// Apply uppercase if requested
	if *uppercase {
		greeting = strings.ToUpper(greeting)
	}

	// Print the greeting the requested number of times
	for i := 0; i < *repeat; i++ {
		fmt.Println(greeting)
	}

	// Exit successfully
	os.Exit(0)
}
