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
	//  Kui käsureal oli kasutusel lipp -m, siis multiLineFlag on true, muidu on see false.
	helpFlag := flag.Bool("h", false, "Enable help mode")
	encodeFlag := flag.Bool("e", false, "Enable encode mode")
	multiLineFlag := flag.Bool("m", false, "Enable multi-line mode")
	// webFlag := flag.Bool("w", false, "Enable web mode")
	flag.Parse()
	// flag.Parse() käivitamine töötleb käsureal sisestatud lüliteid ja määrab iga lülitiga seotud muutuja väärtuse vastavalt sellele,
	// kas lüliti on kasutusel või mitte.

	args := flag.Args()
	// args muutuja hoiab käsurea argumente, mis ei ole lipud. Pärast flag.Parse() väljakutsumist tagastab flag.Args() stringide jada, mis sisaldab
	// kõiki argumente, mis ei olnud lipud. Näiteks, kui programm käivitatakse go run main.go -m "input.txt", siis args sisaldab ["input.txt"]

	if *helpFlag {
		displayTheUsage()
		return
	}
	// var inputBytes []byte
	// var err error

	if *multiLineFlag {
		var input string

		if len(args) == 1 && strings.HasSuffix(args[0], ".encoded.txt") {
			filePath := args[0]
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("\033[41mError:\033[0m\n Error reading file: %s\n", err)
				return
			}
			input = string(fileContent)
			decode(input, *multiLineFlag)

		} else if len(args) == 0 {
			fmt.Println("Enter multi-line input (Ctrl+D to finish):")
			input = handleMultiLineInput()
			decode(input, *multiLineFlag)

			// } else if len(args) == 1 && strings.HasSuffix(args[0], ".encoded.txt") {
			// 	filePath := args[0]
			// 	fileContent, err := os.ReadFile(filePath)
			// 	if err != nil {
			// 		fmt.Println("\033[41mError:\033[0m\n Error reading file: %s\n", err)
			// 		return
			// 	}
			// 	input = string(fileContent)

		} else {
			fmt.Println("\033[41mError:\033[0m\n Invalid usage with -m flag.")
			displayTheUsage()
			return
		}

		if *encodeFlag {
			// Encode the input
			fmt.Println("Encoding not implemented yet")
		} else {
			// Decode the input
			if len(args) == 1 { // kontrollib, kas peale lippude (flags) on antud veel üks argument ja selleks on siis args[0]
				decode(args[0], *multiLineFlag)
			} else {
				displayTheUsage()
				return
			}
			// fmt.Println("That was decoding!")
			fmt.Println(input) // Selle asemel peaks olema decode funktsiooni väljakutse

		}
	} else {

		if len(args) != 1 {
			fmt.Println("\n\033[41mError:\033[0m\n No input provided or too many arguments")
			displayTheUsage()
			return
		}
		// Kui -m lippu ei kasutata
		input := args[0]
		if len(args) == 1 { // kontrollib, kas peale lippude (flags) on antud veel üks argument ja selleks on siis args[0]
			decode(input, *multiLineFlag)
			// fmt.Println("\033[32mNice and simple: That was just decoding without multi-line input!\n\033[0m")

		} else {
			displayTheUsage()
			return
		}
	}

}

//	func decode(input string, multiLine bool) {
//		if multiLine {
//			handleMultiLineInput()
func decode(input string, multiLine bool) {
	if multiLine {
		// Mitmerealise sisendi dekodeerimine
		lines := strings.Split(input, "\n")
		for _, line := range lines {
			decodedLine, success := decodeString(line)
			if success {
				fmt.Println(decodedLine)
			} else {
				// Kui dekodeerimine ebaõnnestub, võib väljastada veateate või käidelda vea
				fmt.Println("\n\033[41mError:\033[0m\n Multiline decoding failed")
				displayTheUsage()

				return
			}
		}
	} else {
		// Decode single line
		decodedString, success := decodeString(input)
		if success {
			fmt.Println(decodedString)
		} else {
			fmt.Println("\n\033[41mError:\033[0m\n Decoding failed")
			return

		}
	}
}

func encode(input string, multiLine bool) {
	if multiLine {
		// Handle multi-line encoding from STDIN or a file
		handleMultiLineInput()
	} else {
		// Encode single line
		fmt.Println(encodeString(input))
	}
}

func decodeString(input string) (string, bool) {

	if !isBracketsBalanced(input) {
		displayTheUsage()
		fmt.Println("\n\033[41mError:\033[0m\n Square brackets are unbalanced\n")
		return "", false
	}

	// var result string
	if strings.Contains(input, "[]") {
		displayTheUsage()
		return "\n\033[41mError:\033[0m\n There is no arguments between square brackets\n", false
	}

	// Implement decoding logic
	pattern := regexp.MustCompile(`\[([^\]]+)\]|([^[]+)`)
	// Esimene osa (\[([^\]]+)\]) otsib tekste, mis on kandiliste sulgude sees, välistades kandiliste sulgude endi esinemise sisus.
	// Teine osa (([^[]+)) otsib tekste, mis ei alga kandilise suluga, püüdes kinni kõik, mis ei ole osa kodeeritud järjestusest.
	matches := pattern.FindAllStringSubmatch(input, -1)

	var result string
	for _, match := range matches {
		if match[1] != "" {
			arguments := strings.SplitN(match[1], " ", 2)
			if len(arguments) != 2 {
				displayTheUsage()
				return "\n\033[41mError:\033[0m\n Incorrect input between brackets: One parameter is missing or no space between two parameters\n", false
			}

			number, err := strconv.Atoi(arguments[0])
			if err != nil || arguments[1] == "" {
				displayTheUsage()
				return "\n\033[41mError:\033[0m\n First argument between brackets is not a number or missing second parameter\n", false
			}

			result += strings.Repeat(arguments[1], number)

		} else if match[2] != "" {
			result += match[2]
		}
	}
	return result, true
}

func encodeString(input string) string {
	// Implement encoding logic
	return ""
}

// Funktsioon mitmerealise sisendi lugemiseks STDIN'ist
func handleMultiLineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func isBracketsBalanced(input string) bool {
	var brackets int

	for _, character := range input {
		switch character {
		case '[':
			brackets++ // Suurendame loendurit iga avava sulgu '[' puhul
		case ']':
			brackets-- // Vähendame loendurit iga sulguva sulgu ']' puhul
			if brackets < 0 {
				// Kui loendur muutub negatiivseks, tähendab see, et leidsime sulguva sulgu
				// enne vastava avava sulgu leidmist, mis on viga
				return false
			}
		}
	}
	// Kui loendur ei ole 0, tähendab see, et on olemas avavaid sulge,
	// millele ei vasta sulguv sulg, mis on samuti viga
	return brackets == 0
}

func displayTheUsage() {
	fmt.Println("\n")

	fmt.Println("\033[41m Usage instructions are coming here: \033[0m")
	fmt.Println("\n")
}

// 	fmt.Println("\033[45mFor decoding\033[0m")

// 	fmt.Println("\033[35mFor single line decoding:          Follow this patter => go run main.go \"[\033[34m[number]\033[35m[single space]\033[34m[character(s)]\033[35m][same logic as in previous brackets][etc.]]\" \033[0m")
// 	fmt.Println("\033[35m             for example:          go run main.go \"[5 #][5 -_]-[5 #]\" \033[0m")
// 	fmt.Println("\033[34mFor decoding from file:            add \"-m\" => Follow this pattern => go run main.go -m \"[file_with_code\033[35m.encoded.txt\033[34m]\" \033[0m")
// 	fmt.Println("\033[35mFor multiline decoding from input: add \"-m\" => go run main.go -m \033[0m")
// 	fmt.Println("\033[35m          for example:             go run main.go -m \"[5 |\\---/|]\\n[5 | o_o |]\\n[5  \\_^_/ ]\"\033[0m")

// 	fmt.Println("\n")
// 	fmt.Println("\033[44mFor encoding\033[0m")
// 	fmt.Println("\033[34mFor single line encoding:          add \"-e\" after main.go (example: go run main.go -e \"[pattern_you_wish_to_encode]\" \033[0m")
// 	fmt.Println("\033[34m             for example:          go run main.go -e \"#####-_-_-_-_-_-#####\" \033[0m")
// 	fmt.Println("\033[35mFor decoding from file:            Follow this pattern => go run main.go -e \"[file_with_code\033[34m.art.txt\033[35m]\" \033[0m")

// 	fmt.Println("\033[44mFor multiline encoding or decoding use files and add \"-e\" or \"-m\" and file name between \" \" \" like this way: \033[45mgo run main.go -m \"cats.encoded.txt\"\033[0m")
// 	fmt.Println("\n")

// }
