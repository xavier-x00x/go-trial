package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var ErrInvalidJRXML = errors.New("invalid JRXML format")

type JasperReport struct {
	Name         string `xml:"name,attr"`
	PageWidth    int    `xml:"pageWidth,attr"`
	PageHeight   int    `xml:"pageHeight,attr"`
	ColumnWidth  int    `xml:"columnWidth,attr"`
	LeftMargin   int    `xml:"leftMargin,attr"`
	RightMargin  int    `xml:"rightMargin,attr"`
	TopMargin    int    `xml:"topMargin,attr"`
	BottomMargin int    `xml:"bottomMargin,attr"`

	Title        *Band `xml:"title>band"`
	PageHeader   *Band `xml:"pageHeader>band"`
	ColumnHeader *Band `xml:"columnHeader>band"`
	Detail       *Band `xml:"detail>band"`
	ColumnFooter *Band `xml:"columnFooter>band"`
	PageFooter   *Band `xml:"pageFooter>band"`
	Summary      *Band `xml:"summary>band"`
	Background   *Band `xml:"background>band"`

	Fields     []Field     `xml:"field"`
	Parameters []Parameter `xml:"parameter"`
	Variables  []Variable  `xml:"variable"`
}

type Band struct {
	Height    int     `xml:"height,attr"`
	SplitType string  `xml:"splitType,attr"`
	YAxis     float64 `xml:"-"`
	Elements  []Element
}

type Element struct {
	Type        string  `json:"type"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Text        string  `json:"txt,omitempty"`
	Field       string  `json:"field,omitempty"`
	Pattern     string  `json:"pattern,omitempty"`
	FontName    string  `json:"fontName,omitempty"`
	FontSize    int     `json:"fontSize,omitempty"`
	FontStyle   string  `json:"fontStyle,omitempty"`
	Align       string  `json:"align,omitempty"`
	VAlign      string  `json:"vAlign,omitempty"`
	Border      string  `json:"border,omitempty"`
	BorderWidth float64 `json:"borderWidth,omitempty"`
	Fill        bool    `json:"fill,omitempty"`
	Link        string  `json:"link,omitempty"`
	PrintWhen   string  `json:"printWhenExpression,omitempty"`
	Stretch     bool    `json:"stretch,omitempty"`
	StretchType string  `json:"stretchType,omitempty"`
	BackColor         string  `json:"backcolor,omitempty"`
	Mode              string  `json:"mode,omitempty"`
	BarcodeType       string  `json:"barcodeType,omitempty"`
	BarcodeExpression string  `json:"barcodeExpression,omitempty"`
}

type Field struct {
	Name  string `xml:"name,attr"`
	Class string `xml:"class,attr"`
}

type Parameter struct {
	Name  string `xml:"name,attr"`
	Class string `xml:"class,attr"`
}

type Variable struct {
	Name        string `xml:"name,attr"`
	Class       string `xml:"class,attr"`
	Calculation string `xml:"calculation,attr"`
	Target      string
}

func ParseJRXML(data []byte) (*JasperReport, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))

	report := &JasperReport{
		Fields:     make([]Field, 0),
		Parameters: make([]Parameter, 0),
		Variables:  make([]Variable, 0),
	}

	type stackItem struct {
		parent string
		band   *Band
	}

	var stack []stackItem
	var currentBand *Band
	var currentElement *Element

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "jasperReport":
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "name":
						report.Name = attr.Value
					case "pageWidth":
						report.PageWidth, _ = strconv.Atoi(attr.Value)
					case "pageHeight":
						report.PageHeight, _ = strconv.Atoi(attr.Value)
					case "columnWidth":
						report.ColumnWidth, _ = strconv.Atoi(attr.Value)
					case "leftMargin":
						report.LeftMargin, _ = strconv.Atoi(attr.Value)
					case "rightMargin":
						report.RightMargin, _ = strconv.Atoi(attr.Value)
					case "topMargin":
						report.TopMargin, _ = strconv.Atoi(attr.Value)
					case "bottomMargin":
						report.BottomMargin, _ = strconv.Atoi(attr.Value)
					}
				}

			case "title", "pageHeader", "columnHeader", "detail", "columnFooter", "pageFooter", "summary", "background":
				currentBand = &Band{
					Elements: make([]Element, 0),
				}
				stack = append(stack, stackItem{parent: t.Name.Local, band: currentBand})

			case "band":
				if currentBand != nil {
					for _, attr := range t.Attr {
						if attr.Name.Local == "height" {
							currentBand.Height, _ = strconv.Atoi(attr.Value)
						} else if attr.Name.Local == "splitType" {
							currentBand.SplitType = attr.Value
						}
					}
				}

			case "reportElement":
				if currentElement != nil {
					for _, attr := range t.Attr {
						switch attr.Name.Local {
						case "x":
							currentElement.X, _ = strconv.Atoi(attr.Value)
						case "y":
							currentElement.Y, _ = strconv.Atoi(attr.Value)
						case "width":
							currentElement.Width, _ = strconv.Atoi(attr.Value)
						case "height":
							currentElement.Height, _ = strconv.Atoi(attr.Value)
						case "stretchType":
							currentElement.StretchType = attr.Value
						case "backcolor":
							currentElement.BackColor = attr.Value
						case "mode":
							currentElement.Mode = attr.Value
						}
					}
				}

			case "textField":
				currentElement = &Element{Type: "textField"}
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "pattern":
						currentElement.Pattern = attr.Value
					case "isStretchWithOverflow":
						if attr.Value == "true" {
							currentElement.Stretch = true
						}
					case "textAdjust":
						if attr.Value == "StretchHeight" {
							currentElement.Stretch = true
						}
					}
				}

			case "staticText":
				currentElement = &Element{Type: "staticText"}
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "textAdjust":
						if attr.Value == "StretchHeight" {
							currentElement.Stretch = true
						}
					}
				}

			case "line":
				currentElement = &Element{Type: "line"}

			case "rectangle":
				currentElement = &Element{Type: "rect"}

			case "ellipse":
				currentElement = &Element{Type: "ellipse"}

			case "textElement":
				if currentElement != nil {
					for _, attr := range t.Attr {
						switch attr.Name.Local {
						case "textAlignment":
							currentElement.Align = attr.Value
						case "verticalAlignment", "vAlign":
							currentElement.VAlign = attr.Value
						}
					}
				}

			case "componentElement":
				currentElement = &Element{Type: "barcode"}

			case "Code128", "Code39", "QRCode", "DataMatrix", "EAN128", "EAN13", "EAN8", "UPCA", "UPCE", "PDF417":
				if currentElement != nil {
					currentElement.BarcodeType = t.Name.Local
				}

			case "codeExpression":
				if currentElement != nil {
					t, err := decoder.Token()
					if err != nil {
						return nil, err
					}
					if cd, ok := t.(xml.CharData); ok {
						currentElement.BarcodeExpression = strings.TrimSpace(string(cd))
					}
				}

			case "font":
				if currentElement != nil {
					for _, attr := range t.Attr {
						switch attr.Name.Local {
						case "fontName":
							currentElement.FontName = attr.Value
						case "size":
							currentElement.FontSize, _ = strconv.Atoi(attr.Value)
						case "isBold":
							if attr.Value == "true" {
								currentElement.FontStyle += "B"
							}
						case "isItalic":
							if attr.Value == "true" {
								currentElement.FontStyle += "I"
							}
						}
					}
				}

			case "box":
				break

			case "pen":
				if currentElement != nil {
					for _, attr := range t.Attr {
						if attr.Name.Local == "lineWidth" {
							val, _ := strconv.ParseFloat(attr.Value, 64)
							if val > 0 {
								currentElement.Border = "1"
								currentElement.BorderWidth = val
							}
						}
					}
				}

			case "topPen":
				if currentElement != nil && currentElement.Border != "1" {
					for _, attr := range t.Attr {
						if attr.Name.Local == "lineWidth" && attr.Value != "0.0" && attr.Value != "0" {
							if !strings.Contains(currentElement.Border, "T") {
								currentElement.Border += "T"
							}
						}
					}
				}

			case "bottomPen":
				if currentElement != nil && currentElement.Border != "1" {
					for _, attr := range t.Attr {
						if attr.Name.Local == "lineWidth" && attr.Value != "0.0" && attr.Value != "0" {
							if !strings.Contains(currentElement.Border, "B") {
								currentElement.Border += "B"
							}
						}
					}
				}

			case "leftPen":
				if currentElement != nil && currentElement.Border != "1" {
					for _, attr := range t.Attr {
						if attr.Name.Local == "lineWidth" && attr.Value != "0.0" && attr.Value != "0" {
							if !strings.Contains(currentElement.Border, "L") {
								currentElement.Border += "L"
							}
						}
					}
				}

			case "rightPen":
				if currentElement != nil && currentElement.Border != "1" {
					for _, attr := range t.Attr {
						if attr.Name.Local == "lineWidth" && attr.Value != "0.0" && attr.Value != "0" {
							if !strings.Contains(currentElement.Border, "R") {
								currentElement.Border += "R"
							}
						}
					}
				}

			case "field":
				f := Field{}
				for _, attr := range t.Attr {
					if attr.Name.Local == "name" {
						f.Name = attr.Value
					} else if attr.Name.Local == "class" {
						f.Class = attr.Value
					}
				}
				report.Fields = append(report.Fields, f)

			case "parameter":
				p := Parameter{}
				for _, attr := range t.Attr {
					if attr.Name.Local == "name" {
						p.Name = attr.Value
					} else if attr.Name.Local == "class" {
						p.Class = attr.Value
					}
				}
				report.Parameters = append(report.Parameters, p)

			case "variable":
				v := Variable{}
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "name":
						v.Name = attr.Value
					case "class":
						v.Class = attr.Value
					case "calculation":
						v.Calculation = attr.Value
					}
				}
				report.Variables = append(report.Variables, v)
			}

		case xml.CharData:
			content := strings.TrimSpace(string(t))
			if content != "" {
				content = strings.TrimPrefix(content, "<![CDATA[")
				content = strings.TrimSuffix(content, "]]>")
				if currentElement != nil && currentElement.Text == "" {
					currentElement.Text = content
				}
			}

		case xml.EndElement:
			switch t.Name.Local {
			case "textField", "staticText", "line", "rectangle", "ellipse", "componentElement":
				if currentElement != nil && currentBand != nil {
					currentBand.Elements = append(currentBand.Elements, *currentElement)
				}
				currentElement = nil

			case "title", "pageHeader", "columnHeader", "detail", "columnFooter", "pageFooter", "summary", "background":
				if len(stack) > 0 {
					item := stack[len(stack)-1]
					stack = stack[:len(stack)-1]

					switch item.parent {
					case "title":
						report.Title = item.band
					case "pageHeader":
						report.PageHeader = item.band
					case "columnHeader":
						report.ColumnHeader = item.band
					case "detail":
						report.Detail = item.band
					case "columnFooter":
						report.ColumnFooter = item.band
					case "pageFooter":
						report.PageFooter = item.band
					case "summary":
						report.Summary = item.band
					case "background":
						report.Background = item.band
					}
				}
				currentBand = nil
			}
		}
	}

	return report, nil
}

func ParseJRXMLFile(filePath string) (*JasperReport, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return ParseJRXML(data)
}

func (jr *JasperReport) GetField(name string) *Field {
	for i := range jr.Fields {
		if jr.Fields[i].Name == name {
			return &jr.Fields[i]
		}
	}
	return nil
}

func (jr *JasperReport) GetParameter(name string) *Parameter {
	for i := range jr.Parameters {
		if jr.Parameters[i].Name == name {
			return &jr.Parameters[i]
		}
	}
	return nil
}
