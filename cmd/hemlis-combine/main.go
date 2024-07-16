package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/filleokus/hemlis/internal/bech32"
	"github.com/filleokus/hemlis/internal/hemlis"
)

func main() {
	filePath := flag.String("f", "", `Path to the file containing shares. 
The file should have the following format:
- Per share, write one word per line (33 words per share)
- Leave a blank line between shares
- (Lines starting with # are ignored as comments)`)
	flag.Parse()
	var combinedSecret string
	var err error

	// Check if there is data in STDIN
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("Reading shares from STDIN. Press Ctrl+D to finish.")
		combinedSecret, err = readAndParseShares(bufio.NewScanner(os.Stdin))
		if err != nil {
			fmt.Printf("Error reading or parsing shares from STDIN: %s\n", err)
			os.Exit(1)
		}
	} else if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			fmt.Printf("error opening file: %w\n", err)
			os.Exit(1)
		}
		defer file.Close()
		combinedSecret, err = readAndParseShares(bufio.NewScanner(file))
		if err != nil {
			fmt.Printf("Error reading or parsing shares from file: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("You must specify a file path using the -file flag or pipe data into STDIN.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Combined secret:", combinedSecret)
	fmt.Println("If not enough shares were used, the above secret will be incorrect")
}

func readAndParseShares(scanner *bufio.Scanner) (string, error) {
	var shares [][]string
	shareIndex := 0
	shares = append(shares, []string{})
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if line == "" {
			shareIndex++
			shares = append(shares, []string{})
			continue
		}
		shares[shareIndex] = append(shares[shareIndex], strings.TrimSpace(line))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
	}
	err := verifyShares(shares)
	if err != nil {
		return "", fmt.Errorf("invalid share (aborting): %s", err)
	}

	bytes, err := decodeSharesToBytes(shares)
	if err != nil {
		return "", fmt.Errorf("invalid share (aborting): %s", err)
	}
	combinedSecret, err := combineShareBytes(bytes)
	if err != nil {
		return "", fmt.Errorf("could not combine shares: %s", err)
	}
	return combinedSecret, nil
}

func verifyShares(shares [][]string) error {
	for i, share := range shares {
		if len(share) != 33 {
			return fmt.Errorf("share %d has %d words, expected 33", i, len(share))
		}
		for j, word := range share {
			if len(word) != 4 {
				return fmt.Errorf("share %d, word %d: (%s) has %d characters, expected 4", i, j, word, len(word))
			}
		}
	}
	return nil
}

func decodeSharesToBytes(shares [][]string) ([][]byte, error) {
	var byteSlices [][]byte

	for i, share := range shares {
		bytes, err := hemlis.DecodeWordsToBytes(share)
		if err != nil {
			return nil, fmt.Errorf("share index %d, could not decode to bytes: %s", i, err)
		}
		fmt.Printf("✅ Share %d (%s) decoded\n", i, hemlis.ShareIdentifier(bytes))
		byteSlices = append(byteSlices, bytes)
	}

	return byteSlices, nil
}

func combineShareBytes(shares [][]byte) (string, error) {
	secretBytes, err := hemlis.CombineSecret(shares)
	if err != nil {
		return "", fmt.Errorf("could not combine secrets: %s", err)
	}
	secretString, err := bech32.Encode("AGE-SECRET-KEY-", secretBytes)
	if err != nil {
		return "", fmt.Errorf("could not bech32 encode secret: %s", err)
	}
	return secretString, nil
}
