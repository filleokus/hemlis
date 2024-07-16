package main

import (
	"strings"
	"testing"
)

func TestGenerateNormalPdf(t *testing.T) {
	document := CreatePDFDocument(PDFParams{ShareIdentifier: "ae7e4", NumberOfShares: "3", Threshold: "2", HemlisVersion: "TEST", CreationDate: "2024-07-15T05:41:03Z", PublicKeyString: "age1sq54yaaevqc85ry6qgjravhtrvm0nshtxv9wd4cpgh9zag09ds9sz5arjl"})
	document.Save("pdf_output/ae7e4.pdf")
}

func TestGenerateRedactedPdf(t *testing.T) {
	document := CreatePDFDocument(PDFParams{ShareIdentifier: "fe551", NumberOfShares: ParamRedacted, Threshold: ParamRedacted, HemlisVersion: "TEST", CreationDate: "2024-07-15T05:41:03Z", PublicKeyString: ParamRedacted})
	document.Save("pdf_output/fe551.pdf")
}

func TestGenerateDangerousPDF(t *testing.T) {
	words := strings.Split("warm gray redo fact ugly vibe knob iris diet wave leaf hope city mint time fuel each glow undo tuna cyan easy dark hope grim stub monk cost play brew webs saga jugs", " ")
	document := CreatePDFDocument(PDFParams{ShareIdentifier: "38bde", NumberOfShares: "3", Threshold: "2", HemlisVersion: "TEST", CreationDate: "2024-07-15T05:41:03Z", PublicKeyString: "age1sq54yaaevqc85ry6qgjravhtrvm0nshtxv9wd4cpgh9zag09ds9sz5arjl", KeyMaterialWords: words})
	document.Save("pdf_output/38bde.pdf")
}
