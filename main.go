package main

import (
	"art/functions"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	encode := flag.Bool("e", false, "Enable encoding mode")
	multiline := flag.Bool("m", false, "Enable multiline mode")
	inputPath := flag.String("i", "", "Input file path")
	outputPath := flag.String("o", "", "Output file path (optional)")
	xorMode := flag.Bool("xor", false, "Enable XOR encryption/decryption mode")
	xorKey := flag.String("key", "", "Key to use for XOR")
	rot13Mode := flag.Bool("rot13", false, "Apply ROT13 cipher to input")

	//parsing flags.
	flag.Parse()
	//checking for CLI arguments
	args := flag.Args()

	var inputContent string
	var err error

	//if -i flag is given try to read from given filepath, else read from CLI argument. 
	if *inputPath != "" {
		inputContent, err = functions.ReadTxtFile(*inputPath, *multiline)
		if err != nil {
			fmt.Println("Failed to read input file:", err)
			os.Exit(1)
		}
	} else if len(args) > 0 {
		inputContent = strings.Join(args, " ")
	} else { //no input given give feedback and exit.
		fmt.Println("No input provided. Use -i to read from file or pass input string as argument.")
		os.Exit(1)
	}

	// Encoding or decoding
	var result string
	switch {
		case *xorMode: 
			if *xorKey == "" {
				fmt.Println("Error: XOR key must be provided with --key")
				os.Exit(1)
			}
			result, err = functions.Xorify(inputContent, *xorKey)
			if err != nil {
				fmt.Println("XOR error: ", err)
				os.Exit(1)
			}
		case *rot13Mode:
			result = functions.Rot13ify(inputContent)

		case *encode:
			result = functions.EncodeString(inputContent, *multiline)
		
		default:
			result = functions.DecodeString(inputContent, *multiline)
	}
	// Output
	if *outputPath != "" {
		err = functions.WriteTxtFile(*outputPath, result)
		if err != nil {
			fmt.Println("Failed to write to output file:", err)
			os.Exit(1)
		}
		fmt.Println("Output written to:", *outputPath)
	} else {
		fmt.Println("Output:")
		fmt.Println(result)
	}
}


