package main

import (
	"fmt"
	"log"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func CreatePDFDocument(parameters PDFParams) core.Document {
	m := initDocument(parameters)
	document, err := m.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
	return document
}

type PDFParams struct {
	ShareIdentifier  string
	NumberOfShares   string
	Threshold        string
	HemlisVersion    string
	CreationDate     string
	PublicKeyString  string
	KeyMaterialWords []string
}

const ParamRedacted = "REDACTED"

func initDocument(parameters PDFParams) core.Maroto {
	m := maroto.New(config.NewBuilder().WithMaxGridSize(95).Build())

	m.AddRow(12, text.NewCol(45, "Secret Share", props.Text{Size: 28, Align: align.Left, Style: fontstyle.Bold, Family: "Helvetica"}))
	m.AddRow(12,
		text.NewCol(45, parameters.ShareIdentifier, props.Text{Size: 23, Align: align.Left, Style: fontstyle.Normal, Family: "Courier"}),
		col.New(5),
		text.NewCol(45, "Parameters", props.Text{Size: 25, Align: align.Left, Style: fontstyle.Bold, Family: "Helvetica"}),
	)
	m.AddRow(3,
		line.NewCol(45, props.Line{Thickness: 1, SizePercent: 100}),
		col.New(5),
		line.NewCol(45, props.Line{Thickness: 1, SizePercent: 100}),
	)

	m.AddRow(85,
		col.New(45).Add(
			text.New("This is a share of a private encryption key, keep it safe.", props.Text{Size: 20, Align: align.Left, Style: fontstyle.Normal, Family: "Helvetica", Color: &props.Color{Red: 232, Green: 36, Blue: 4}}),
			text.New("This document was produced by Hemlis (github.com/filleokus/hemlis). It contains a share of an age private key split using Samir Secret Sharing as implemented by Hashicorp's Vault. The share is encoded into words from the Bytewords wordlist.", props.Text{Size: 12, Align: align.Left, Style: fontstyle.Normal, Family: "Helvetica", VerticalPadding: 1.25, Top: 20}),
			text.New("Using Hemlis (or another implementation) the private key can be recovered if enough shares are collected.", props.Text{Size: 12, Align: align.Left, Style: fontstyle.Normal, Family: "Helvetica", VerticalPadding: 1.25, Top: 20 + 35}),
		),
		col.New(5),
		parameterSection(parameters),
	)
	m.AddRow(12,
		text.NewCol(95, "Key material", props.Text{Size: 25, Align: align.Left, Style: fontstyle.Bold, Family: "Helvetica"}),
	)
	m.AddRow(5,
		line.NewCol(95, props.Line{Thickness: 1, SizePercent: 100}),
	)
	for _, row := range keyMaterialSection(parameters) {
		m.AddRows(row)
	}

	return m
}

func parameterStyleWithTop(top float64) props.Text {
	return props.Text{Size: 10, Align: align.Left, Style: fontstyle.Normal, Family: "Courier", VerticalPadding: 1.25, Top: top}
}

func parameterSection(parameters PDFParams) core.Col {
	lineMaxLength := 40
	paddingTop := 15.0
	col := col.New(45)
	col.Add(
		text.New("During the creation of the secret shares the following parameters were used", props.Text{Size: 12, Align: align.Left, Style: fontstyle.Normal, Family: "Helvetica", VerticalPadding: 1.25}),
		text.New(fmt.Sprintf("Hemlis version: %s", parameters.HemlisVersion), parameterStyleWithTop(paddingTop)),
		text.New(fmt.Sprintf("Creation date: %s", parameters.CreationDate), parameterStyleWithTop(paddingTop+5*1)),
		text.New(fmt.Sprintf("Number of shares: %s", parameters.NumberOfShares), parameterStyleWithTop(paddingTop+5*2)),
		text.New(fmt.Sprintf("Treshold: %s", parameters.Threshold), parameterStyleWithTop(paddingTop+5*3)),
	)
	if parameters.PublicKeyString != ParamRedacted {
		col.Add(
			text.New(fmt.Sprintf("Public Key: %s", parameters.PublicKeyString[0:lineMaxLength-len("Public Key: ")]), parameterStyleWithTop(paddingTop+5*4)),
			text.New(parameters.PublicKeyString[lineMaxLength-len("Public Key: "):], parameterStyleWithTop(paddingTop+5*5)),
			code.NewQr(parameters.PublicKeyString, props.Rect{Center: false, Top: paddingTop + 5*7, Left: 30, Percent: 40}),
			text.New("Public Key", props.Text{Size: 8, Align: align.Center, Style: fontstyle.Normal, Family: "Courier", Top: paddingTop + 5*7 + 45, Left: 5}),
		)
	} else {
		col.Add(text.New(fmt.Sprintf("Public Key: %s", parameters.PublicKeyString), parameterStyleWithTop(paddingTop+5*4)))
	}
	return col
}

func keyMaterialSection(parameters PDFParams) []core.Row {
	rows := []core.Row{}
	styling := props.Text{Size: 16, Align: align.Left, Family: "Courier", Top: 0}
	words := make([]string, 33)
	if parameters.KeyMaterialWords == nil || len(parameters.KeyMaterialWords) == 0 {
		for i := 0; i < 33; i++ {
			words[i] = "__   __   __   __"
		}
	} else {
		copy(words, parameters.KeyMaterialWords)
	}

	for i := 0; i < 16; i++ {
		row := row.New(8)
		row.Add(
			text.NewCol(40, fmt.Sprintf("%3d: %s", i, words[i]), styling),
			col.New(5),
			text.NewCol(40, fmt.Sprintf("%3d: %s", i+16, words[i+16]), styling),
		)
		rows = append(rows, row)
	}
	rows = append(rows, row.New(8).Add(col.New(45), text.NewCol(40, fmt.Sprintf("%3d: %s", 32, words[32]), styling)))
	return rows
}
