package services

import (
	"database/sql"
	"demo/internal/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"

	_ "github.com/go-sql-driver/mysql"
)

func Execute(db *sql.DB, ReportMonth string) {
	fmt.Println("Executing...")

	// Calculate the first and last day of the month
	t, _ := time.Parse("2006-01", ReportMonth)

	bizdatefrom := ReportMonth[:4] + "-01-01"
	bizdateto := t.AddDate(0, 1, -1).Format("2006-01-02")
	monthfrom := ReportMonth[:4] + "-01"
	monthto := ReportMonth
	fmt.Println("Querying for Retailer")
	retailers := models.QueryRetailer(db)
	fmt.Println("Retailers:", len(retailers))
	fmt.Println("Querying for POS Sales")
	pos := models.QueryPOS(db, bizdatefrom, bizdateto)
	fmt.Println("POS:", len(pos))
	fmt.Println("Querying for Active POS")
	activePOS := models.QueryActivePOS(db, bizdatefrom, bizdateto)
	fmt.Println("Active POS:", len(activePOS))
	fmt.Println("Querying for GR")
	gr := models.QueryGR(db, monthfrom, monthto)
	fmt.Println("GR:", len(gr))
	fmt.Println("Querying for Sup")
	sup := models.QuerySup(db, monthfrom, monthto)
	fmt.Println("Sup:", len(sup))
	fmt.Println("Querying for Retailer POS Payment")
	retailerPOSPayment := models.QueryRetailerPOSPayment(db, bizdatefrom, bizdateto)
	fmt.Println("Retailer POS Payment:", len(retailerPOSPayment))

	// -- Generate Excel --
	f := excelize.NewFile()
	defer f.Close()

	// -- POS Payment Summary --
	f.SetSheetName("Sheet1", "POS Payment Summary")
	BuildPOSPaymentSummary(f, retailerPOSPayment, ReportMonth)

	// -- Retailer Summary --
	f.SetSheetName("Sheet2", "Retailer Summary")
	BuildRetailerSummary(f, retailers, sup, pos, ReportMonth, activePOS)

	// -- Retail detail --
	EachRetailerSummary(f, retailers, pos, activePOS, ReportMonth)

	// ── Save the file to <exe_dir>/excel/ ───────────────
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	outputDir := filepath.Join(exeDir, "excel")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	outputFile := filepath.Join(outputDir, fmt.Sprintf("demo_%s.xlsx", ReportMonth))
	if err := f.SaveAs(outputFile); err != nil {
		log.Fatalf("Error saving Excel file: %v", err)
	}

	fmt.Printf("\n✅ Report saved: %s\n", outputFile)

}
