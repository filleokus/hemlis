package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/filleokus/hemlis/internal/hemlis"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

func GeneratePDF(secret hemlis.GeneratedSecret, share hemlis.Share, pdfOptions PDFOptions) core.Document {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	parameters := PDFParams{
		HemlisVersion:    "0.1.0",
		CreationDate:     currentTime,
		NumberOfShares:   ParamRedacted,
		Threshold:        ParamRedacted,
		PublicKeyString:  ParamRedacted,
		KeyMaterialWords: []string{},
		ShareIdentifier:  share.Identifier(),
	}
	if pdfOptions.IncludeNumberOfShares {
		parameters.NumberOfShares = fmt.Sprintf("%d", secret.NumberOfShares())
	}
	if pdfOptions.IncludeThreshold {
		parameters.Threshold = fmt.Sprintf("%d", secret.Threshold())
	}
	if pdfOptions.InlcudePublicKey {
		parameters.PublicKeyString = secret.PublicKeyString()
	}
	if pdfOptions.IncludeWords {
		parameters.KeyMaterialWords = share.Words()
	}
	return CreatePDFDocument(parameters)
}

func PrintSecretToCLI(secret hemlis.GeneratedSecret) {
	shares := secret.Shares()
	for shareIndex, share := range shares {
		fmt.Printf("Share %d (%s)\n", shareIndex+1, share.Identifier())
		words := share.Words()
		const chunkSize = 5
		for chunkIndex := 0; chunkIndex < len(words); chunkIndex += chunkSize {
			end := chunkIndex + 5
			if end > len(words) {
				end = len(words)
			}
			chunk := words[chunkIndex:end]
			for wordIndex, word := range chunk {
				fmt.Printf("| %2d %-4s ", chunkIndex+wordIndex, word)
				if (wordIndex+1)%chunkSize == 0 {
					fmt.Println("|")
				}
			}
		}
		fmt.Println(("|\n--------------------------------"))
	}
}

func SaveSecretToDisk(secret hemlis.GeneratedSecret) {
	// Open a new file for writing only
	file, err := os.Create(fmt.Sprintf("%s.txt", secret.PublicKeyString()))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write strings to the file
	fmt.Fprintf(file, "# Public Key: %s\n", secret.PublicKeyString())
	fmt.Fprintf(file, "# Private Key: %s\n", secret.PrivateKeyString())
	for i, share := range secret.Shares() {
		fmt.Fprintf(file, "# Share %d (%s)\n", i, share.Identifier())
		for _, word := range share.Words() {
			fmt.Fprintf(file, "%s\n", word)
		}
		fmt.Fprintf(file, "\n")
	}
}

type PDFOptions struct {
	IncludeNumberOfShares bool
	IncludeThreshold      bool
	InlcudePublicKey      bool
	IncludeWords          bool
}

func main() {
	// Define flags
	numberOfShares := flag.Uint("shares", 0, "Number of shares to generate")
	threshold := flag.Uint("threshold", 0, "Number of shares required to reconstruct the secret")
	output := flag.String("output", "pdf,cli,txt", "Output format (pdf, txt, cli) separated by comma (default: pdf, cli, txt)")
	includeNumberOfShares := flag.Bool("pdf-include-num-shares", true, "Include the number of shares in the PDF")
	includeThreshold := flag.Bool("pdf-include-num-threshold", true, "Include the threshold number in the PDF")
	includePublicKey := flag.Bool("pdf-include-public-key", true, "Include the public key in the PDF")
	includeWords := flag.Bool("pdf-include-wordlist", false, "Include words in PDF to avoid writing by hand (not recomended)")
	// TODO: Implement this
	// outputDir := flag.String("pdf-output-dir", "./", "Directory to save the generated files (default: current directory)")
	// combinePDFs := flag.Bool("pdf-merge", false, "Combine all generated PDFs into one")
	// txtFilePath := flag.String("txt-path", "", "Path to save the private key, public key and shares as a text file (default: value of public key)")

	flag.Parse()

	if (numberOfShares == nil || threshold == nil) || (*numberOfShares < 1 || *threshold < 1) {
		fmt.Println("Number of shares must and threshold must be at least 2")
		flag.Usage()
		os.Exit(1)
	}

	secret, _ := hemlis.GenerateSecret(*numberOfShares, *threshold)
	fmt.Printf("Public Key: %s\n", secret.PublicKeyString())
	fmt.Printf("Private Key: %s\n", secret.PrivateKeyString())
	fmt.Printf("Number of shares: %d\n", secret.NumberOfShares())
	fmt.Printf("Threshold: %d\n", secret.Threshold())
	fmt.Println(("--------------------------------"))

	shouldOutputPDF := strings.Contains(*output, "pdf")
	shouldOutputTXT := strings.Contains(*output, "txt")
	shouldOutputCLI := strings.Contains(*output, "cli")

	if shouldOutputCLI {
		PrintSecretToCLI(*secret)
	}

	if shouldOutputTXT {
		fmt.Println("Saving txt file")
		SaveSecretToDisk(*secret)
	}

	if shouldOutputPDF {
		fmt.Println("Generating PDF's")
		for _, share := range secret.Shares() {
			pdf := GeneratePDF(*secret, share, PDFOptions{
				IncludeNumberOfShares: *includeNumberOfShares,
				IncludeThreshold:      *includeThreshold,
				InlcudePublicKey:      *includePublicKey,
				IncludeWords:          *includeWords,
			})
			pdf.Save(fmt.Sprintf("%s.pdf", share.Identifier()))
		}
		if !*includeWords {
			fmt.Println("Print the PDF's and manually write the words on the papers")
		} else {
			fmt.Println("Print the PDF's on a secure printer as the PDF contains the secret words")
		}
	}
}
