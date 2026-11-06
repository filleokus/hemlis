package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/filleokus/hemlis/internal/bech32"
	"github.com/filleokus/hemlis/internal/hemlis"
	"github.com/filleokus/hemlis/internal/version"
)

func main() {
	printVersion := flag.Bool("v", false, "Print version information")
	filePath := flag.String("f", "", `Path to the file containing shares. 
The file should have the following format:
- Per share, write one word per line (33 words per share)
- Leave a blank line between shares
- (Lines starting with # are ignored as comments)`)
	flag.Parse()

	if *printVersion {
		fmt.Printf("hemlis-combine %s\ncommit %s, built at %s\n", version.Version, version.Commit, version.BuildDate)
		os.Exit(0)
	}

	var combinedSecret string
	var err error

	if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		combinedSecret, err = readAndParseShares(bufio.NewScanner(file))
		if err != nil {
			fmt.Printf("Error reading or parsing shares from file: %s\n", err)
			os.Exit(1)
		}
	} else if stdinIsPiped() {
		fmt.Println("Reading shares from STDIN. Press Ctrl+D to finish.")
		combinedSecret, err = readAndParseShares(bufio.NewScanner(os.Stdin))
		if err != nil {
			fmt.Printf("Error reading or parsing shares from STDIN: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("You must specify a file path using the -f flag or pipe data into STDIN.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Combined secret:", combinedSecret)
	fmt.Println("If not enough shares were used, the above secret will be incorrect")
}

// stdinIsPiped reports whether STDIN is a pipe or file rather than a terminal.
func stdinIsPiped() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func readAndParseShares(scanner *bufio.Scanner) (string, error) {
	shares, err := parseShares(scanner)
	if err != nil {
		return "", err
	}

	combinedSecret, err := combineShares(shares)
	if err != nil {
		return "", fmt.Errorf("could not combine shares: %s", err)
	}
	return combinedSecret, nil
}

// parseShares reads shares from the scanner. Each share is a block of
// words, one per line. Blocks are separated by one or more blank lines,
// and lines starting with # are ignored. Surrounding whitespace is ignored.
func parseShares(scanner *bufio.Scanner) ([][]string, error) {
	var shares [][]string
	var current []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		if line == "" {
			if len(current) > 0 {
				shares = append(shares, current)
				current = nil
			}
			continue
		}
		current = append(current, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read shares: %w", err)
	}
	if len(current) > 0 {
		shares = append(shares, current)
	}
	if len(shares) == 0 {
		return nil, fmt.Errorf("no shares found in input")
	}
	return shares, nil
}

func combineShares(shares [][]string) (string, error) {
	for i, share := range shares {
		if len(share) != 33 {
			return "", fmt.Errorf("share %d has %d words, expected 33", i, len(share))
		}
		for j, word := range share {
			if len(word) != 4 {
				return "", fmt.Errorf("share %d, word %d: (%s) has %d characters, expected 4", i, j, word, len(word))
			}
		}
	}

	var byteSlices [][]byte

	for i, share := range shares {
		bytes, err := hemlis.DecodeWordsToBytes(share)
		if err != nil {
			return "", fmt.Errorf("share index %d, could not decode to bytes: %s", i, err)
		}
		fmt.Printf("✅ Share %d (%s) decoded\n", i, hemlis.ShareIdentifier(bytes))
		byteSlices = append(byteSlices, bytes)
	}

	secretBytes, err := hemlis.CombineSecret(byteSlices)
	if err != nil {
		return "", fmt.Errorf("could not combine secrets: %s", err)
	}
	secretString, err := bech32.Encode("AGE-SECRET-KEY-", secretBytes)
	if err != nil {
		return "", fmt.Errorf("could not bech32 encode secret: %s", err)
	}
	return secretString, nil

}
