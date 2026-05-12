package gojasperxml

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"go-trial/pkg/go-jasperxml/parser"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
	"github.com/jung-kurt/gofpdf"
)

type templateGenerator struct {
	report   *parser.JasperReport
	pdf      *gofpdf.Fpdf
	mm2pt    float64
	margin   float64
	pageW    float64
	pageH    float64
	rowIndex    int
	fontDir     string
	loadedFonts map[string]bool
}

func newTemplateGenerator(report *parser.JasperReport) *templateGenerator {
	return &templateGenerator{
		report:      report,
		mm2pt:       2.8346,
		fontDir:     "pkg/go-jasperxml/fonts",
		loadedFonts: make(map[string]bool),
	}
}

func (g *templateGenerator) Generate(data []map[string]interface{}, params ...interface{}) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	pageWidth := float64(g.report.PageWidth) / g.mm2pt
	pageHeight := float64(g.report.PageHeight) / g.mm2pt

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size:    gofpdf.SizeType{Wd: pageWidth, Ht: pageHeight},
	})

	leftMargin := float64(g.report.LeftMargin) / g.mm2pt
	rightMargin := float64(g.report.RightMargin) / g.mm2pt
	topMargin := float64(g.report.TopMargin) / g.mm2pt
	bottomMargin := float64(g.report.BottomMargin) / g.mm2pt

	pdf.SetMargins(leftMargin, topMargin, rightMargin)
	pdf.SetAutoPageBreak(true, bottomMargin)

	g.pdf = pdf
	if g.fontDir != "" {
		g.pdf.SetFontLocation(g.fontDir)
	}
	g.loadFonts()

	g.margin = leftMargin
	g.pageW = pageWidth
	g.pageH = pageHeight

	pdf.AddPage()

	firstRow := data[0]

	if g.report.Title != nil {
		g.renderBand(g.report.Title, firstRow, g.report.Title.Height)
	}

	if g.report.PageHeader != nil {
		g.renderBand(g.report.PageHeader, firstRow, g.report.PageHeader.Height)
	}

	if g.report.ColumnHeader != nil {
		g.renderBand(g.report.ColumnHeader, firstRow, g.report.ColumnHeader.Height)
	}

	for idx, row := range data {
		g.rowIndex = idx + 1
		if g.report.Detail != nil {
			g.renderBand(g.report.Detail, row, g.report.Detail.Height)
		}
	}

	if g.report.ColumnFooter != nil {
		g.renderBand(g.report.ColumnFooter, firstRow, g.report.ColumnFooter.Height)
	}

	if g.report.Summary != nil {
		g.renderBand(g.report.Summary, firstRow, g.report.Summary.Height)
	}

	if g.report.PageFooter != nil {
		y := pageHeight - bottomMargin - float64(g.report.PageFooter.Height)/g.mm2pt
		g.pdf.SetY(y)
		g.renderBand(g.report.PageFooter, firstRow, g.report.PageFooter.Height)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *templateGenerator) renderBand(band *parser.Band, row map[string]interface{}, height int) {
	if band == nil {
		return
	}

	type ElementData struct {
		Text             string
		CalculatedHeight float64
	}
	evaluatedElements := make([]ElementData, len(band.Elements))

	bandHeight := float64(height) / g.mm2pt
	maxY := 0.0

	for i, elem := range band.Elements {
		var text string
		if elem.Type == "textField" {
			text = g.resolveExpression(elem.Text, row, elem.Pattern)
		} else if elem.Type == "staticText" {
			text = elem.Text
		} else if elem.Type == "barcode" {
			text = g.resolveExpression(elem.BarcodeExpression, row, "")
		}

		calcHeight := float64(elem.Height) / g.mm2pt

		if elem.Stretch && text != "" && (elem.Type == "textField" || elem.Type == "staticText") {
			fontName := elem.FontName
			if fontName == "" {
				fontName = "Helvetica"
			}
			fontStyle := ""
			if strings.Contains(elem.FontStyle, "B") {
				fontStyle += "B"
			}
			if strings.Contains(elem.FontStyle, "I") {
				fontStyle += "I"
			}
			fontSize := float64(elem.FontSize)
			if fontSize == 0 {
				fontSize = 10
			}

			g.pdf.SetFont(fontName, fontStyle, fontSize)

			width := float64(elem.Width) / g.mm2pt
			lines := g.pdf.SplitLines([]byte(text), width)
			if len(lines) > 0 {
				fontLineHeight := (fontSize / g.mm2pt) * 1.15
				// Include 1.0mm top padding + 0.5mm bottom padding in the calculation
				calcHeight = fontLineHeight*float64(len(lines)) + 1.5
				if calcHeight < float64(elem.Height)/g.mm2pt {
					calcHeight = float64(elem.Height) / g.mm2pt
				}
			}
		}

		evaluatedElements[i] = ElementData{
			Text:             text,
			CalculatedHeight: calcHeight,
		}

		elemBottom := float64(elem.Y)/g.mm2pt + calcHeight
		if elemBottom > maxY {
			maxY = elemBottom
		}
	}

	if maxY > bandHeight {
		bandHeight = maxY
	}

	bottomMargin := float64(g.report.BottomMargin) / g.mm2pt

	if g.pdf.GetY()+bandHeight > g.pageH-bottomMargin+0.1 {
		g.pdf.AddPage()
	}

	currentY := g.pdf.GetY()

	for i, elem := range band.Elements {
		elemHeight := evaluatedElements[i].CalculatedHeight
		if elem.StretchType == "RelativeToTallestObject" {
			maxPossibleHeight := bandHeight - (float64(elem.Y) / g.mm2pt)
			if maxPossibleHeight > elemHeight {
				elemHeight = maxPossibleHeight
			}
		}
		g.renderElement(elem, evaluatedElements[i].Text, elemHeight, bandHeight, currentY)
	}

	g.pdf.SetY(currentY + bandHeight)
}

func (g *templateGenerator) renderElement(elem parser.Element, preEvaluatedText string, calcHeight float64, bandHeight float64, currentY float64) {
	pdf := g.pdf

	x := float64(elem.X)/g.mm2pt + g.margin
	y := float64(elem.Y)/g.mm2pt + currentY
	width := float64(elem.Width) / g.mm2pt
	totalHeight := calcHeight

	hAlignStr := ""
	switch elem.Align {
	case "Center":
		hAlignStr = "C"
	case "Right":
		hAlignStr = "R"
	default:
		hAlignStr = "L"
	}

	vAlignStr := ""
	switch elem.VAlign {
	case "Middle":
		vAlignStr = "M"
	case "Bottom":
		vAlignStr = "B"
	default:
		vAlignStr = "T"
	}

	// Handle background color
	if elem.Mode == "Opaque" && elem.BackColor != "" {
		r, g, b := hexToRGB(elem.BackColor)
		pdf.SetFillColor(r, g, b)
		if elem.Type != "ellipse" {
			pdf.Rect(x, y, width, totalHeight, "F")
		}
	}

	switch elem.Type {
	case "textField":
		fontName := g.mapFontName(elem.FontName)
		if fontName == "" {
			fontName = "Helvetica"
		}

		fontStyle := ""
		if strings.Contains(elem.FontStyle, "B") {
			fontStyle += "B"
		}
		if strings.Contains(elem.FontStyle, "I") {
			fontStyle += "I"
		}

		fontSize := float64(elem.FontSize)
		if fontSize == 0 {
			fontSize = 10
		}

		pdf.SetFont(fontName, fontStyle, fontSize)

		pdf.SetXY(x, y)
		borderStr := elem.Border
		if borderStr != "" {
			if elem.BorderWidth > 0 {
				pdf.SetLineWidth(elem.BorderWidth)
			} else {
				pdf.SetLineWidth(0.2)
			}
			pdf.CellFormat(width, totalHeight, "", borderStr, 0, "", false, 0, "")
			pdf.SetXY(x, y)
		}

		const topPad = 1.0
		const bottomPad = 1.0
		if elem.Stretch && preEvaluatedText != "" {
			fontLineHeight := (fontSize / g.mm2pt) * 1.15
			lines := pdf.SplitLines([]byte(preEvaluatedText), width)
			multicellHeight := fontLineHeight * float64(len(lines))

			drawY := y + topPad
			if vAlignStr == "M" {
				drawY = y + (totalHeight-multicellHeight)/2
			} else if vAlignStr == "B" {
				drawY = y + totalHeight - multicellHeight - bottomPad
			}

			pdf.SetXY(x, drawY)
			pdf.MultiCell(width, fontLineHeight, preEvaluatedText, "", hAlignStr, false)
		} else {
			if vAlignStr == "T" {
				pdf.SetXY(x, y+topPad)
				pdf.CellFormat(width, totalHeight-topPad, preEvaluatedText, "", 0, hAlignStr+vAlignStr, false, 0, "")
			} else if vAlignStr == "B" {
				pdf.CellFormat(width, totalHeight-bottomPad, preEvaluatedText, "", 0, hAlignStr+vAlignStr, false, 0, "")
			} else {
				pdf.CellFormat(width, totalHeight, preEvaluatedText, "", 0, hAlignStr+vAlignStr, false, 0, "")
			}
		}

	case "staticText":
		fontName := g.mapFontName(elem.FontName)
		if fontName == "" {
			fontName = "Helvetica"
		}

		fontStyle := ""
		if strings.Contains(elem.FontStyle, "B") {
			fontStyle = "B"
		}

		fontSize := float64(elem.FontSize)
		if fontSize == 0 {
			fontSize = 10
		}

		pdf.SetFont(fontName, fontStyle, fontSize)

		pdf.SetXY(x, y)
		borderStr := elem.Border
		if borderStr != "" {
			if elem.BorderWidth > 0 {
				pdf.SetLineWidth(elem.BorderWidth)
			} else {
				pdf.SetLineWidth(0.2)
			}
			pdf.CellFormat(width, totalHeight, "", borderStr, 0, "", false, 0, "")
			pdf.SetXY(x, y)
		}

		const topPad = 1.0
		const bottomPad = 1.0
		if elem.Stretch && preEvaluatedText != "" {
			fontLineHeight := (fontSize / g.mm2pt) * 1.15
			lines := pdf.SplitLines([]byte(preEvaluatedText), width)
			multicellHeight := fontLineHeight * float64(len(lines))

			drawY := y + topPad
			if vAlignStr == "M" {
				drawY = y + (totalHeight-multicellHeight)/2
			} else if vAlignStr == "B" {
				drawY = y + totalHeight - multicellHeight - bottomPad
			}

			pdf.SetXY(x, drawY)
			pdf.MultiCell(width, fontLineHeight, preEvaluatedText, "", hAlignStr, false)
		} else {
			if vAlignStr == "T" {
				pdf.SetXY(x, y+topPad)
				pdf.CellFormat(width, totalHeight-topPad, preEvaluatedText, "", 0, hAlignStr+vAlignStr, false, 0, "")
			} else if vAlignStr == "B" {
				pdf.CellFormat(width, totalHeight-bottomPad, preEvaluatedText, "", 0, hAlignStr+vAlignStr, false, 0, "")
			} else {
				pdf.CellFormat(width, totalHeight, preEvaluatedText, "", 0, hAlignStr+vAlignStr, false, 0, "")
			}
		}

	case "barcode":
		if preEvaluatedText != "" {
			var b barcode.Barcode
			var err error

			switch elem.BarcodeType {
			case "Code128":
				b, err = code128.Encode(preEvaluatedText)
			case "QRCode":
				b, err = qr.Encode(preEvaluatedText, qr.M, qr.Auto)
			case "EAN13":
				b, err = ean.Encode(preEvaluatedText)
			}

			if err == nil {
				// Scale the barcode
				imgW := int(width * 10) // scale up for resolution
				imgH := int(totalHeight * 10)
				b, _ = barcode.Scale(b, imgW, imgH)

				// Convert to NRGBA to ensure 8-bit depth for gofpdf
				bounds := b.Bounds()
				img := image.NewNRGBA(bounds)
				draw.Draw(img, bounds, b, bounds.Min, draw.Src)

				var buf bytes.Buffer
				png.Encode(&buf, img)

				imageName := fmt.Sprintf("barcode_%v_%v", elem.X, elem.Y)
				pdf.RegisterImageOptionsReader(imageName, gofpdf.ImageOptions{ImageType: "PNG"}, &buf)
				pdf.ImageOptions(imageName, x, y, width, totalHeight, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
			}
		}

	case "line":
		x1 := x
		y1 := y
		x2 := x + width
		y2 := y

		if elem.Height > 0 {
			y2 = y + float64(elem.Height)/g.mm2pt
		}

		pdf.SetLineWidth(0.5)
		pdf.Line(x1, y1, x2, y2)

	case "rect":
		if elem.BorderWidth > 0 {
			pdf.SetLineWidth(elem.BorderWidth)
		} else {
			pdf.SetLineWidth(0.5)
		}
		pdf.Rect(x, y, width, totalHeight, "D")

	case "ellipse":
		centerX := x + width/2
		centerY := y + totalHeight/2
		radiusX := width / 2
		radiusY := totalHeight / 2

		if elem.BorderWidth > 0 {
			pdf.SetLineWidth(elem.BorderWidth)
		} else {
			pdf.SetLineWidth(0.5)
		}
		style := "D"
		if elem.Mode == "Opaque" && elem.BackColor != "" {
			style = "FD"
		}
		pdf.Ellipse(centerX, centerY, radiusX, radiusY, 0, style)
	}
}

func (g *templateGenerator) resolveExpression(expr string, row map[string]interface{}, pattern string) string {
	expr = strings.TrimSpace(expr)
	expr = strings.ReplaceAll(expr, "<![CDATA[", "")
	expr = strings.ReplaceAll(expr, "]]>", "")
	expr = strings.ReplaceAll(expr, "$V{REPORT_COUNT}", strconv.Itoa(g.rowIndex))

	// Resolve ternary expressions first
	ternaryRE := regexp.MustCompile(`\$F\{(\w+)\} \!= null \? \$F\{(\w+)\} : "?([^":]*)"?`)
	expr = ternaryRE.ReplaceAllStringFunc(expr, func(match string) string {
		parts := strings.Split(match, " != null ? ")
		if len(parts) != 2 {
			return match
		}

		fieldName := strings.Trim(parts[0], "$F{} ")
		fieldName = strings.Trim(fieldName, "}")

		rest := parts[1]
		if strings.Contains(rest, " : ") {
			trueAndFalse := strings.Split(rest, " : ")
			if len(trueAndFalse) == 2 {
				truePart := strings.TrimSpace(trueAndFalse[0])
				trueField := strings.Trim(truePart, "$F{}")

				if val, ok := row[trueField]; ok && val != nil && val != "" {
					return g.formatNumber(val)
				}

				falseVal := strings.Trim(trueAndFalse[1], "\" ")
				return falseVal
			}
		}

		if val, ok := row[fieldName]; ok && val != nil && val != "" {
			return g.formatNumber(val)
		}
		return ""
	})

	// Now evaluate concatenation and fields
	var result strings.Builder
	inQuote := false
	var currentToken strings.Builder

	flushToken := func() {
		tok := strings.TrimSpace(currentToken.String())
		currentToken.Reset()
		if tok == "" || tok == "+" {
			return
		}
		if strings.HasPrefix(tok, "$F{") && strings.HasSuffix(tok, "}") {
			fieldName := tok[3 : len(tok)-1]
			if val, ok := row[fieldName]; ok && val != nil {
				if pattern == "#,##0" || pattern == "" {
					result.WriteString(g.formatNumber(val))
				} else {
					result.WriteString(formatValue(val, pattern))
				}
			}
		} else {
			result.WriteString(tok)
		}
	}

	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if c == '"' {
			if inQuote {
				inQuote = false
				result.WriteString(currentToken.String())
				currentToken.Reset()
			} else {
				flushToken()
				inQuote = true
			}
		} else if c == '+' && !inQuote {
			flushToken()
		} else {
			currentToken.WriteByte(c)
		}
	}
	flushToken()

	return result.String()
}

func (g *templateGenerator) formatNumber(v interface{}) string {
	switch val := v.(type) {
	case string:
		if n, err := strconv.ParseFloat(val, 64); err == nil {
			return g.formatNumberFloat(n)
		}
		return val
	case float64:
		return g.formatNumberFloat(val)
	case int:
		return strconv.FormatFloat(float64(val), 'f', 0, 64)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (g *templateGenerator) formatNumberFloat(n float64) string {
	if n == float64(int64(n)) {
		return formatWithSeparator(int64(n))
	}

	s := fmt.Sprintf("%.2f", n)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]

	if len(intPart) <= 3 {
		return intPart + "." + decPart
	}

	var result strings.Builder
	length := len(intPart)
	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteByte(intPart[i])
	}
	return result.String() + "." + decPart
}

func formatWithSeparator(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}

	var result strings.Builder
	length := len(s)
	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteByte(s[i])
	}
	return result.String()
}

func formatValue(v interface{}, pattern string) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%v", val)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (g *templateGenerator) mapFontName(name string) string {
	if name == "" {
		return ""
	}
	
	lowerName := strings.ToLower(name)
	if g.loadedFonts[lowerName] {
		return name // Use custom font if loaded
	}

	if strings.Contains(lowerName, "arial") || strings.Contains(lowerName, "sans") {
		return "Helvetica"
	}
	if strings.Contains(lowerName, "times") || strings.Contains(lowerName, "serif") {
		return "Times"
	}
	if strings.Contains(lowerName, "courier") || strings.Contains(lowerName, "mono") {
		return "Courier"
	}
	return name
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return int(r), int(g), int(b)
}

func (g *templateGenerator) loadFonts() {
	if g.fontDir == "" {
		return
	}

	files, err := os.ReadDir(g.fontDir)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name())
		if strings.ToLower(ext) == ".json" {
			name := strings.TrimSuffix(f.Name(), ext)
			style := ""
			if strings.HasSuffix(name, "b") {
				style = "B"
				name = strings.TrimSuffix(name, "b")
			} else if strings.HasSuffix(name, "i") {
				style = "I"
				name = strings.TrimSuffix(name, "i")
			} else if strings.HasSuffix(name, "bi") {
				style = "BI"
				name = strings.TrimSuffix(name, "bi")
			}

			// Add font to gofpdf
			// Since SetFontLocation is set, we only need the filename
			g.pdf.AddFont(name, style, f.Name())
			g.loadedFonts[strings.ToLower(name)] = true
		}
	}
}
