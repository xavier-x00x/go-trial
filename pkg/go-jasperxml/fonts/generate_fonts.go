package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	fontDir := "."
	// Find gofpdf directory dynamically to get the encoding map
	out, err := exec.Command("go", "list", "-f", "{{.Dir}}", "github.com/jung-kurt/gofpdf").Output()
	if err != nil {
		fmt.Printf("Error finding gofpdf directory: %v\n", err)
		return
	}
	gofpdfDir := strings.TrimSpace(string(out))
	encodingFile := filepath.Join(gofpdfDir, "font", "cp1252.map")

	files, err := os.ReadDir(fontDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".ttf") {
			fmt.Printf("Generating definition for %s...\n", f.Name())
			// MakeFont(fontFile, encodingFile, fontDir, output, embed)
			err := gofpdf.MakeFont(filepath.Join(fontDir, f.Name()), encodingFile, fontDir, os.Stdout, true)
			if err != nil {
				fmt.Printf("Error generating %s: %v\n", f.Name(), err)
			} else {
				fmt.Printf("Successfully generated definition for %s\n", f.Name())
			}
		}
	}
}

/*
Catatan Penting:
- Untuk mendaftarkan font ke gofpdf, kita perlu menggunakan fungsi MakeFont.
- MakeFont akan menghasilkan file json dan z untuk setiap file ttf.
- File json dan z tersebut harus berada di direktori font.
- SetFontLocation digunakan untuk memberitahu gofpdf di mana mencari file json dan z.

cara menjalankan:

```sh
go run pkg/go-jasperxml/fonts/generate_fonts.go
```

Outputnya adalah file .json dan .z di folder yang sama
*/
