package gojasperxml

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"go-trial/pkg/go-jasperxml/parser"
)

var ErrNoData = errors.New("no data provided")
var ErrNoTemplate = errors.New("no template provided")

type JasperXML struct {
	data     []map[string]interface{}
	params   map[string]interface{}
	template *parser.JasperReport
}

func NewFromJSON(data []byte) (*JasperXML, error) {
	var parsed []map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	if len(parsed) == 0 {
		return nil, ErrNoData
	}

	return &JasperXML{
		data:     parsed,
		params:   make(map[string]interface{}),
		template: nil,
	}, nil
}

func NewFromFile(filePath string) (*JasperXML, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewFromJSON(data)
}

func NewFromJRXML(templatePath string) (*JasperXML, error) {
	report, err := parser.ParseJRXMLFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JRXML: %w", err)
	}

	return &JasperXML{
		data:     []map[string]interface{}{},
		params:   make(map[string]interface{}),
		template: report,
	}, nil
}

func NewFromJRXMLWithData(templatePath string, dataPath string) (*JasperXML, error) {
	report, err := parser.ParseJRXMLFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JRXML: %w", err)
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse data JSON: %w", err)
	}

	if len(parsed) == 0 {
		return nil, ErrNoData
	}

	return &JasperXML{
		data:     parsed,
		params:   make(map[string]interface{}),
		template: report,
	}, nil
}

func (j *JasperXML) SetData(data []map[string]interface{}) {
	j.data = data
}

func (j *JasperXML) GetData() []map[string]interface{} {
	return j.data
}

func (j *JasperXML) SetParameter(key string, value interface{}) {
	j.params[key] = value
}

func (j *JasperXML) SetParameters(params map[string]interface{}) {
	for k, v := range params {
		j.params[k] = v
	}
}

func (j *JasperXML) GeneratePDF() ([]byte, error) {
	if j.template != nil {
		gen := newTemplateGenerator(j.template)
		return gen.Generate(j.data, j.params)
	}

	if len(j.data) == 0 {
		return nil, ErrNoData
	}

	return nil, ErrNoTemplate

	// gen := newBasicGenerator()
	// return gen.Generate(j.data, j.params)
}

func (j *JasperXML) SavePDF(outputPath string) error {
	pdfBytes, err := j.GeneratePDF()
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, pdfBytes, 0644)
}
