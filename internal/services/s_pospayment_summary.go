package services

import (
	"demo/internal/models"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func BuildPOSPaymentSummary(f *excelize.File, retailerPOSPayment []models.RetailerPOSPayment, reportMonth string) {
	startRow := 1
	f.SetCellValue("POS Payment Summary", fmt.Sprintf("A%d", startRow), "Retailer POS Payment Summary")
	startRow += 2
	GenerateExcel(f, "POS Payment Summary", retailerPOSPayment, models.RetailerPOSPayment{}, reportMonth, startRow, 1)

}
