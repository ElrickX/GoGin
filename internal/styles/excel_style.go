package styles

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

var styleCache = make(map[string]int)

func HeaderStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D9E1F2"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	return style
}
func GetStyle(
	f *excelize.File,
	format string,
	align string,
	withBorder bool,
) int {

	// 🟢 key（唯一识别）
	key := fmt.Sprintf("%s|%s|%v", format, align, withBorder)

	// 🟢 如果已有，直接用
	if id, ok := styleCache[key]; ok {
		return id
	}

	// =========================
	// 🟢 build style
	// =========================
	style := &excelize.Style{}

	// 👉 format
	switch format {
	case "currency":
		style.NumFmt = 3
	case "percent":
		style.NumFmt = 10
	case "date":
		style.NumFmt = 14
	case "number":
		style.NumFmt = 3
	default:
		style.NumFmt = 49 // ✅ TEXT 格式
	}

	// 👉 align
	if align != "" {
		style.Alignment = &excelize.Alignment{
			Horizontal: align,
		}
	}

	// 👉 border
	if withBorder {
		style.Border = []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		}
	}

	// 🟢 创建 style
	styleID, _ := f.NewStyle(style)

	// 🟢 cache
	styleCache[key] = styleID

	return styleID
}
