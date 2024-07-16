package main

import (
	"fmt"
	"os"
	"time"

	"github.com/filleokus/hemlis/internal/hemlis"
)

func GeneratePDF(secret hemlis.GeneratedSecret, pdfOptions PDFOptions) {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	parameters := PDFParams{
		HemlisVersion:    "0.1.0",
		CreationDate:     currentTime,
		NumberOfShares:   ParamRedacted,
		Threshold:        ParamRedacted,
		PublicKeyString:  ParamRedacted,
		KeyMaterialWords: []string{},
		ShareIdentifier:  "",
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

	for _, share := range secret.Shares() {
		if pdfOptions.DangerousIncludePrivateKey {
			parameters.KeyMaterialWords = share.Words
		}
		document := CreatePDFDocument(parameters)
		document.Save(fmt.Sprintf("share-%s.pdf", share.Identifier))
	}
}

func PrintSecretToCLI(secret hemlis.GeneratedSecret) {
	fmt.Printf("Public Key: %s\n", secret.PublicKeyString())
	fmt.Printf("Private Key: %s\n", secret.PrivateKeyString())
	fmt.Printf("Number of shares: %d\n", secret.NumberOfShares())
	fmt.Printf("Threshold: %d\n", secret.Threshold())
	fmt.Println(("--------------------------------"))
	shares := secret.Shares()
	for shareIndex, share := range shares {
		fmt.Printf("Share %d (%s)\n", shareIndex+1, share.Identifier)
		words := share.Words
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
		fmt.Fprintf(file, "# Share %d (%s)\n", i, share.Identifier)
		for _, word := range share.Words {
			fmt.Fprintf(file, "%s\n", word)
		}
		fmt.Fprintf(file, "\n")
	}
}

type PDFOptions struct {
	IncludeNumberOfShares      bool
	IncludeThreshold           bool
	InlcudePublicKey           bool
	DangerousIncludePrivateKey bool
}

func main() {
	secret, _ := hemlis.GenerateSecret(6, 4)
	PrintSecretToCLI(*secret)
	GeneratePDF(*secret, PDFOptions{IncludeNumberOfShares: true, IncludeThreshold: true, InlcudePublicKey: true, DangerousIncludePrivateKey: true})
	SaveSecretToDisk(*secret)
}
