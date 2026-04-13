package main

import (
	"demo/internal/services"
	"demo/panda_lib"
	"fmt"
	"time"
)

func main() {
	db := panda_lib.ConfigDB()
	defer db.Close()

	fmt.Println("Connected to database successfully!")
	reportMonth := ""
	if reportMonth == "" {
		// Default: previous month
		now := time.Now()
		prev := now.AddDate(0, -1, 0)
		reportMonth = prev.Format("2006-01")
	}
	services.Execute(db, reportMonth)
}
