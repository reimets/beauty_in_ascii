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

	if len(args) != 1 {
		fmt.Println("\nError:\n No input provided or too many arguments")
		displayTheUsage()
		return
	}

	var inputBytes []byte
	var input string
	var err error

	// Kui -m lipp on kasutusel koos failinimega
	if *multiLineFlag && strings.HasSuffix(args[0], ".txt") {
		inputBytes, err = os.ReadFile(args[0]) // os.ReadFile tagastab baitide massiivi ([]byte) ja vea.
		if err != nil {
			fmt.Printf("Error:\nError reading file: %s\n", err)
			return
		}
		input = string(inputBytes) // baitide massiiv ([]byte) teisendatakse stringiks.
	} else if *multiLineFlag {
		// Kui -m lipp on kasutusel, kuid sisend on antud otse käsurealt
		input = strings.Replace(args[0], "\\n", "\n", -1) // "-1" tähistab, et asendamine tuleks teha kõikide leidude puhul
		/*
			s: string, mida töödeldakse.
			old: alamstring, mida soovite asendada.
			new: alamstring, millega soovite old asendada.
			n: asenduste maksimaalne arv. -1 tähendab, et asendatakse kõik esinemised.
		*/
	} else {
		// Kui kasutatakse ainult ühte argumenti ilma -m liputa
		input = args[0]
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
		fmt.Println("Decoding not implemented yet")
	}

	// Edasine töötlus 'input' muutujaga
	fmt.Println(input)
}

func decode(input string, multiLine bool) {
	if multiLine {
		handleMultiLineInput(false)
	} else {
		// Decode single line
		decodedString, success := decodeString(input)
		if !success {
			fmt.Println(decodedString)
			return
		} else {
			fmt.Println(decodedString)
		}
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

	if !isBracketsBalanced(input) {
		displayTheUsage()
		fmt.Println("\nError:\n Square brackets are unbalanced")
		return "", false
	}

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

	fmt.Println("\033[35mFor single line decoding:          Follow this patter => go run main.go \"[\033[34m[number]\033[35m[single space]\033[34m[character(s)]\033[35m][same logic as in previous brackets][etc.]]\" \033[0m")
	fmt.Println("\033[35m             for example:          go run main.go \"[5 #][5 -_]-[5 #]\" \033[0m")
	fmt.Println("\033[34mFor decoding from file:            add \"-m\" => Follow this pattern => go run main.go -m \"[file_with_code\033[35m.encoded.txt\033[34m]\" \033[0m")
	fmt.Println("\033[35mFor multiline decoding from input: add \"-m\" => go run main.go -m \033[0m")
	fmt.Println("\033[35m          for example:             go run main.go -m \"[5 |\\---/|]\\n[5 | o_o |]\\n[5  \\_^_/ ]\"\033[0m")

	fmt.Println("\n")
	fmt.Println("\033[44mFor encoding\033[0m")
	fmt.Println("\033[34mFor single line encoding:          add \"-e\" after main.go (example: go run main.go -e \"[pattern_you_wish_to_encode]\" \033[0m")
	fmt.Println("\033[34m             for example:          go run main.go -e \"#####-_-_-_-_-_-#####\" \033[0m")
	fmt.Println("\033[35mFor decoding from file:            Follow this pattern => go run main.go -e \"[file_with_code\033[34m.art.txt\033[35m]\" \033[0m")

	fmt.Println("\033[44mFor multiline encoding or decoding use files and add \"-e\" or \"-m\" and file name between \" \" \" like this way: \033[45mgo run main.go -m \"cats.encoded.txt\"\033[0m")
	fmt.Println("\n")

}
