package services

import (
	"demo/internal/styles"
	"fmt"
	"reflect"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func getCategory(revenue float64, posCount int, ytdDays int) (string, float64) {
	annualizedRevenue := revenue
	if ytdDays > 0 {
		annualizedRevenue = (revenue / float64(ytdDays)) * 365.0
	}
	if annualizedRevenue > 1000000000 || posCount >= 401 {
		return "A", annualizedRevenue
	}
	if annualizedRevenue > 500000000 || posCount > 250 {
		return "B", annualizedRevenue
	}
	if annualizedRevenue > 100000000 || posCount >= 100 {
		return "C", annualizedRevenue
	}
	return "D", annualizedRevenue
}

func GenerateExcel(
	f *excelize.File,
	sheetName string,
	data interface{},
	model interface{},
	reportMonth string,
	startRow int,
	startCol int,
) {

	t := reflect.TypeOf(model)

	// =========================
	// 🟢 Step 1: Header
	// =========================
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 👉 header name
		header := field.Tag.Get("excel")
		if header == "" {
			header = field.Name
		}

		col, _ := excelize.ColumnNumberToName(startCol + i)
		cell := fmt.Sprintf("%s%d", col, startRow)

		f.SetCellValue(sheetName, cell, header)

		// 👉 header style
		f.SetCellStyle(sheetName, cell, cell, styles.HeaderStyle(f))

		// 👉 width
		width := 15.0
		if w := field.Tag.Get("width"); w != "" {
			if val, err := strconv.Atoi(w); err == nil {
				width = float64(val)
			}
		}
		f.SetColWidth(sheetName, col, col, width)
	}

	startRow++

	// =========================
	// 🟢 Step 2: Data
	// =========================
	vData := reflect.ValueOf(data)

	for r := 0; r < vData.Len(); r++ {
		v := vData.Index(r)

		for c := 0; c < t.NumField(); c++ {
			field := t.Field(c)

			col, _ := excelize.ColumnNumberToName(startCol + c)
			cell := fmt.Sprintf("%s%d", col, startRow+r)

			val := v.Field(c).Interface()
			f.SetCellValue(sheetName, cell, val)

			// =========================
			// 🟢 format
			// =========================
			format := field.Tag.Get("format")
			align := field.Tag.Get("align")
			styleID := styles.GetStyle(f, format, align, true)

			f.SetCellStyle(sheetName, cell, cell, styleID)

		}
	}
}
