package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/emeraldls/platnova-task/internal/types"
	"github.com/unidoc/unipdf/v3/contentstream/draw"
	"github.com/unidoc/unipdf/v3/creator"
)

func ReadJSONFile(filePath string) (*types.AccountStatement, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	fileReader := bytes.NewReader(fileData)
	var statement types.AccountStatement
	err = json.NewDecoder(fileReader).Decode(&statement)
	if err != nil {
		return nil, err
	}

	return &statement, nil
}

func NewLine(c *creator.Creator, x1, y1, x2, y2 float64, isAbsolute, fillWidth bool, lineWidth float64,
	color creator.Color, isDashed bool, dashArray []int64, dashPhase int64, opacity float64, margins creator.Margins) *creator.Line {
	positioning := creator.PositionRelative
	if isAbsolute {
		positioning = creator.PositionAbsolute
	}

	fitMode := creator.FitModeNone
	if fillWidth {
		fitMode = creator.FitModeFillWidth
	}

	style := draw.LineStyleSolid
	if isDashed {
		style = draw.LineStyleDashed
	}

	line := c.NewLine(x1, y1, x2, y2)
	line.SetLineWidth(lineWidth)
	line.SetMargins(margins.Left, margins.Right, margins.Top, margins.Bottom)
	line.SetPositioning(positioning)
	line.SetFitMode(fitMode)
	line.SetColor(color)
	line.SetStyle(style)
	line.SetDashPattern(dashArray, dashPhase)
	line.SetOpacity(opacity)
	return line
}

func NewMargins(left, right, top, bottom float64) creator.Margins {
	return creator.Margins{
		Left:   left,
		Right:  right,
		Top:    top,
		Bottom: bottom,
	}
}

func NewPara(c *creator.Creator, text string, textStyle creator.TextStyle) *creator.StyledParagraph {
	p := c.NewStyledParagraph()
	p.Append(text).Style = textStyle
	p.SetEnableWrap(false)
	return p
}

func DrawCell(table *creator.Table, content creator.VectorDrawable, cellStyle types.CellStyle) error {
	var cell = table.NewCell()

	err := cell.SetContent(content)
	if err != nil {
		return err
	}
	cell.SetHorizontalAlignment(cellStyle.HAlignment)
	if cellStyle.Indent > 0 {
		cell.SetIndent(cellStyle.Indent)
	}

	return nil

}
