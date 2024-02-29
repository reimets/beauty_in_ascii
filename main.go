package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Defines command-line flags with default values and descriptions.
	helpFlag := flag.Bool("h", false, "Enable help mode")
	encodeFlag := flag.Bool("e", false, "Enable encode mode")
	multiLineFlag := flag.Bool("m", false, "Enable multi-line mode")
	flag.Parse() // Processes the command-line flags and sets the corresponding variables.

	args := flag.Args() // Returns the command-line arguments that are not flags.

	// If the help flag is used, displays usage instructions and exits the program.
	if *helpFlag {
		displayTheUsage()
		return
	}

	var input string

	// Checks if an argument with .encoded.txt suffix is provided and reads its content.
	if len(args) == 1 && strings.HasSuffix(args[0], ".encoded.txt") {
		filePath := args[0]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\033[41mError:\033[0m\n Error reading file! \"%s\"\n", err)
			return
		}
		input = string(fileContent)

		decode(input, *multiLineFlag)
		return
	}

	// Checks if an argument with .art.txt suffix is provided and reads its content.
	if len(args) == 1 && strings.HasSuffix(args[0], ".art.txt") {
		filePath := args[0]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\033[41mError:\033[0m\n Error reading file! \"%s\"\n", err)
			return
		}
		input = string(fileContent)

		encode(input, *multiLineFlag)
		return
	}

	// If the encode flag is used and there is one argument provided.
	if *encodeFlag {
		if len(args) == 1 { // checks if one more argument is given besides flags, and for that it is args[0]
			encode(args[0], *multiLineFlag)
			fmt.Println()
			return
		}
	}

	// Handles multi-line input according to the multi-line flag.
	if *multiLineFlag {
		var input string
		if *encodeFlag {
			// Encode the input
			fmt.Println("Enter multi-line input for encoding (Ctrl+D to finish):")
			input = handleMultiLineInput()
			encode(input, *encodeFlag)
		} else if len(args) == 0 {
			fmt.Println("Enter multi-line input for decoding (Ctrl+D to finish):")
			input = handleMultiLineInput()
			decode(input, *multiLineFlag)
		} else {
			fmt.Println("\n\033[41mError:\033[0m\n Invalid usage with -m or -e flag.")
			displayTheUsage()
			return
		}

		if len(args) == 1 { // checks if one more argument is given besides flags, and for that it is args[0]
			decode(args[0], *multiLineFlag)
		} else {
			return
		}

	} else {
		if len(args) != 1 {
			fmt.Println("\n\033[41mError:\033[0m\nNo correct input provided or too many arguments.\nOr maybe you didn't used \"\"-s?")
			displayTheUsage()
			return
		}
		// If the -m flag is not used
		input = args[0]
		if len(args) == 1 {
			decode(input, *multiLineFlag)

		} else {
			displayTheUsage()
			return
		}
	}

	// Checks for empty square brackets in arguments.
	if strings.Contains(args[0], "[]") {
		fmt.Println("\n\033[41mError:\033[0m\n There are no arguments between square brackets")
		fmt.Println()
		return
	}
}

func decode(input string, multiLine bool) {
	if multiLine {
		// Splits the input into lines for multi-line decoding.
		lines := strings.Split(input, "\n")
		for _, line := range lines {
			decodedLine, success := decodeString(line)
			if success {
				fmt.Println(decodedLine)

			} else {
				fmt.Println("\n\033[41mError:\033[0m\n Multiline decoding failed")
				displayTheUsage()

				return
			}
		}
	} else {
		// Decodes a single line of input.
		decodedString, success := decodeString(input)
		if success {
			fmt.Println(decodedString)
		} else {
			fmt.Println("\n\033[41mError:\033[0m\n Decoding failed - check arguments between brackets [ ]")
			fmt.Println()
			return
		}
	}
}

// Logic for decoding, using regular expressions and checking for balanced brackets.
// Additional checks and the implementation of the decoding process.
func decodeString(input string) (string, bool) {
	if !isBracketsBalanced(input) {
		displayTheUsage()
		fmt.Println("\n\033[41mError:\033[0m\n Square brackets are unbalanced")
		return "", false
	}

	if strings.Contains(input, "[]") {
		displayTheUsage()
		return "", false
	}

	// Implement decoding logic
	pattern := regexp.MustCompile(`\[([^\]]+)\]|([^[]+)`)
	matches := pattern.FindAllStringSubmatch(input, -1)

	var result string
	for _, match := range matches {
		if match[1] != "" {
			arguments := strings.SplitN(match[1], " ", 2)
			if len(arguments) != 2 {
				displayTheUsage()
				return "", false
			}

			number, err := strconv.Atoi(arguments[0])
			if err != nil || arguments[1] == "" {
				displayTheUsage()
				return "", false
			}

			result += strings.Repeat(arguments[1], number)

		} else if match[2] != "" {
			result += match[2]
		}
	}
	return result, true
}

// Similar logic to the decode function but for encoding.
// Uses strings.Builder to create the encoded string.
func encode(input string, multiLine bool) {
	if multiLine {
		lines := strings.Split(input, "\n")
		for _, line := range lines {
			encodedLine, success := encodeString(line)
			if success {
				fmt.Println(encodedLine)
			} else {
				fmt.Println("\n\033[41mError:\033[0m\n Multiline encoding failed")
				displayTheUsage()
				return
			}
		}
	} else {
		encodedString, success := encodeString(input)
		if success {
			fmt.Println(encodedString)
		} else {
			fmt.Println("\n\033[41mError:\033[0m\n Encoding failed! Maybe you didn't used \"-s? ")
			fmt.Println()
			return
		}
	}
}

// Logic for encoding, handling character repetition and encoding accordingly.
func encodeString(input string) (string, bool) {
	var encodedBuilder strings.Builder
	count := 1

	for i := 1; i <= len(input); i++ {
		if i == len(input) || input[i] != input[i-1] {
			if count > 1 {
				encodedBuilder.WriteString(fmt.Sprintf("[%d %c]", count, input[i-1]))
			} else {
				encodedBuilder.WriteRune(rune(input[i-1]))
			}
			count = 1
		} else {
			count++
		}
	}
	return encodedBuilder.String(), true
}

// Reads multi-line input from standard input using bufio.Scanner.
func handleMultiLineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// Checks if the square brackets in the input are balanced.
func isBracketsBalanced(input string) bool {
	var brackets int

	for _, character := range input {
		switch character {
		case '[':
			brackets++
		case ']':
			brackets--
			if brackets < 0 {
				return false
			}
		}
	}
	return brackets == 0
}

// Displays usage instructions, including examples of how to use the program.
func displayTheUsage() {
	fmt.Println("\n\033[41m Usage instructions: \033[0m")
	fmt.Println()

	fmt.Println("\033[45mFor decoding\033[0m")

	fmt.Println("\033[35mFor single line decoding:          Follow this patter => go run main.go \"[\033[34m[number]\033[35m[single space]\033[34m[character(s)]\033[35m][same logic as in previous brackets][etc.]]\" \033[0m")
	fmt.Println("\033[35m             for example:          go run main.go \"[5 #][5 -_]-[5 #]\" \033[0m")
	fmt.Println("\033[34mFor decoding from file:            use file with the end \033[35m\".encoded.txt\"\033[34m. Example: go run main.go cats.encoded.txt\033[0m")
	fmt.Println("\033[35mFor multiline decoding from input: type \"go run main.go -m\" \033[0m")
	fmt.Println("\033[35mand into the next lines insert coded lines you want to decode.\033[0m")
	fmt.Println("\033[35mfor example:                       \n[5 |\\---/|]\n[5 | o_o |]\n[5  \\_^_/ ]\033[0m")
	fmt.Println("\033[45m\033[1m NB! After completing the multi-line input in the terminal, please push \"enter\" and then the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. \033[0m\033[22m")

	fmt.Println()
	fmt.Println("\033[44mFor encoding\033[0m")
	fmt.Println("\033[34mFor single line encoding:          add \"-e\" after main.go (For example: go run main.go -e \"[pattern_you_wish_to_encode]\") \033[0m")
	fmt.Println("\033[34m             for example:          go run main.go -e \"#####-_-_-_-_-_-#####\" \033[0m")
	fmt.Println("\033[35mFor decoding from file:            use file with the end \033[34m\".art.txt\"\033[35m. For example: go run main.go cats.art.txt \033[0m")
	fmt.Println("\033[34mFor multiline encoding from input: add \"-m\" & \"-e\" (example: go run main.go -m -e)\033[0m")
	fmt.Println("\033[34mand next lines enter for example:  \n" +
		"          \n" +
		"   *   *  \n" +
		"  *** *** \n" +
		"  ******* \n" +
		"   *****  \n" +
		"    ***   \n" +
		"     *    \n\033[0m")
	fmt.Println("\033[44m\033[1m NB! After completing the multi-line input in the terminal, please push \"enter\" and then the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. \033[0m\033[22m")
}
