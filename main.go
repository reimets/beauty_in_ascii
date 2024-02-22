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

	if len(os.Args) != 2 && len(os.Args) != 3 && len(os.Args) != 4 {
		displayTheUsage()
		return
	}

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
	// kõiki argumente, mis ei olnud lipud. Näiteks, kui programm käivitatakse go run main.go -e -m input.txt, siis args sisaldab ["input.txt"]
	if *helpFlag {
		displayTheUsage()
		return
	}

	if !isBracketsBalanced(args[0]) {
		displayTheUsage()
		fmt.Println("\nError:\n Square brackets are unbalanced")

		return
	}

	if *encodeFlag {
		// Handle encoding
		if len(args) == 1 {
			encode(args[0], *multiLineFlag) // args[0]: esimene argument käsurealt, mis ei ole lipp.
			// Eeldatavasti on see tekst, mida kasutaja soovib kodeerida või failinimi, mille sisu tuleb kodeerida.
			// *multiLineFlag: boolean väärtus, mis näitab, kas mitmerealine režiim on aktiveeritud.
			// Kui käsureal oli kasutusel lipp -m, siis multiLineFlag on true, muidu on see false.
		} else {
			displayTheUsage()
			return
		}
	} else {
		// Handle decoding
		if len(args) == 1 { // kontrollib, kas peale lippude (flags) on antud veel üks argument ja selleks on siis args[0]
			decode(args[0], *multiLineFlag)
		} else {
			displayTheUsage()
			return
		}
	}
}

func decode(input string, multiLine bool) {
	if multiLine {
		handleMultiLineInput(false)
	} else {
		// Decode single line
		decodedString, _ := decodeString(input)
		// if !success {
		// 	fmt.Println(decodedString)
		// 	return
		// } else {
		fmt.Println(decodedString)
		// }
	}
}

func encode(input string, multiLine bool) {
	if multiLine {
		// Handle multi-line encoding from STDIN or a file
		handleMultiLineInput(true)
	} else {
		// Encode single line
		fmt.Println(encodeString(input))
	}
}

func decodeString(input string) (string, bool) {
	// var result string
	if strings.Contains(input, "[]") {
		displayTheUsage()
		return "\nError:\n There is no arguments between square brackets", false
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
				return "\nError:\n Incorrect input between brackets: One parameter is missing or no space between two parameters\n", false
			}

			number, err := strconv.Atoi(arguments[0])
			if err != nil || arguments[1] == "" {
				displayTheUsage()
				return "\nError:\n First argument between brackets is not a number or missing second parameter\n", false
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

func handleMultiLineInput(isEncoding bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if isEncoding {
			fmt.Println(encodeString(line))
		} else {
			fmt.Println(decodeString(line))
		}
	}
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

	fmt.Println("\033[41m Usage instructions here: \033[0m")
	fmt.Println("\n")
	fmt.Println("\033[45mFor decoding\033[0m")

	fmt.Println("\033[35mFor single line decoding:          Follow this patter => go run main.go \"[[number][single space][character(s)]][same logic as in previous brackets][etc.]]\" \033[0m")
	fmt.Println("\033[35m             for example:          go run main.go \"[5 #][5 -_]-[5 #]\" \033[0m")
	fmt.Println("\033[34mFor decoding from file:            Follow this pattern => go run main.go [file_with_code.txt] \033[0m")
	fmt.Println("\033[35mFor multiline decoding from input: add \"-m\" => go run main.go -m \033[0m")
	fmt.Println("\033[35mand enter for example:             \n[5 |\\---/|]\\n[5 | o_o |]\\n[5  \\_^_/ ] \n or \n[5 |\\---/|]\n[5 | o_o |]\n[5  \\_^_/ ]\033[0m")
	fmt.Println("\033[45m\033[1m NB! After completing the multi-line input in the terminal, please enter the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. \033[0m\033[22m")

	fmt.Println("\n")
	fmt.Println("\033[44mFor encoding\033[0m")
	fmt.Println("\033[34mFor single line encoding:          add \"-e\" after main.go (example: go run main.go -e \"[pattern_you_wish_to_encode]\" \033[0m")
	fmt.Println("\033[34m             for example:          go run main.go -e \"#####-_-_-_-_-_-#####\" \033[0m")
	fmt.Println("\033[35mFor decoding from file:            Follow this pattern => go run main.go -e [file_with_code.txt] \033[0m")

	fmt.Println("\033[34mFor multiline encoding from input: add \"-e\" & \"-m\" (example: go run main.go -e -m)\033[0m")
	fmt.Println("\033[34mand next lines enter for example:  \n" +
		"|\\---/||\\---/||\\---/||\\---/||\\---/|\n" +
		"| o_o || o_o || o_o || o_o || o_o |\n" +
		" \\_^_/  \\_^_/  \\_^_/  \\_^_/  \\_^_/ \033[0m")
	fmt.Println("\033[44m\033[1m NB! After completing the multi-line input in the terminal, please enter the EOF (End Of File) character by pressing CTRL+D on Linux/MacOS systems or CTRL+Z on Windows systems. This signals to the program that input reading is finished. \033[0m\033[22m")

}
