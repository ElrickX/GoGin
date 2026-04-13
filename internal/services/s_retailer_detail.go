package services

import (
	"demo/internal/models"
	"demo/internal/styles"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func BuildBizSummaryByRetailer(
	pos []models.POS,
	reportMonth string,
) map[string][]RetailerPOSBizSummaryRow {
	t, _ := time.Parse("2006-01", reportMonth)

	start := t
	end := t.AddDate(0, 1, 0) // next month
	// reportYear := reportMonth[:4]
	reportYearInt, _ := strconv.Atoi(reportMonth[:4])

	// 🟢 Step 1: group by retailer → date
	type Daily struct {
		Amount float64
		Trans  int
	}

	retailerMap := make(map[string]map[time.Time]*Daily)

	for _, p := range pos {
		if p.BizDate.Year() != reportYearInt {
			continue
		}

		if _, ok := retailerMap[p.RetailerGUID]; !ok {
			retailerMap[p.RetailerGUID] = make(map[time.Time]*Daily)
		}

		date := p.BizDate.Truncate(24 * time.Hour) // 🔥 去掉时间部分

		if _, ok := retailerMap[p.RetailerGUID][date]; !ok {
			retailerMap[p.RetailerGUID][date] = &Daily{}
		}

		d := retailerMap[p.RetailerGUID][date]
		d.Amount += p.Amount
		d.Trans += p.Trans
	}

	// 🟢 Step 2: build result per retailer
	result := make(map[string][]RetailerPOSBizSummaryRow)

	for retailerGUID, dateMap := range retailerMap {

		// sort dates
		// 👉 用 time 排序
		var dates []time.Time
		for d := range dateMap {
			dates = append(dates, d)
		}

		sort.Slice(dates, func(i, j int) bool {
			return dates[i].Before(dates[j])
		})

		var runningAmount float64
		var runningTrans int
		var ytdAmount float64
		var ytdTrans int

		var rows []RetailerPOSBizSummaryRow

		for _, d := range dates {
			day := dateMap[d]

			// YTD (all months)
			ytdAmount += day.Amount
			ytdTrans += day.Trans

			// only output report month
			// 👉 只输出 reportMonth
			if d.Before(start) || !d.Before(end) {
				continue
			}

			// month accumulate
			runningAmount += day.Amount
			runningTrans += day.Trans

			rows = append(rows, RetailerPOSBizSummaryRow{
				BizDate:             d,
				Amount:              day.Amount,
				AccumulateAmount:    runningAmount,
				AccumulateAmountYTD: ytdAmount,
				Trans:               day.Trans,
				AccumulateTrans:     runningTrans,
				AccumulateTransYTD:  ytdTrans,
			})
		}

		result[retailerGUID] = rows
	}

	return result
}

func BuildOutletSummaryByRetailer(
	pos []models.POS,
	reportMonth string,
) map[string][]RetailerSalesContribute {

	reportYM := reportMonth

	// 🟢 retailer → outlet aggregation
	type OutletAgg struct {
		Amount     float64
		Trans      int
		PosSet     map[string]struct{} // for PosActive count
		OutletCode string
		OutletName string
	}

	retailerMap := make(map[string]map[string]*OutletAgg)

	for _, p := range pos {
		if p.BizDate.Format("2006-01") != reportYM {
			continue
		}

		if _, ok := retailerMap[p.RetailerGUID]; !ok {
			retailerMap[p.RetailerGUID] = make(map[string]*OutletAgg)
		}

		if _, ok := retailerMap[p.RetailerGUID][p.OutletGUID]; !ok {
			retailerMap[p.RetailerGUID][p.OutletGUID] = &OutletAgg{
				PosSet:     make(map[string]struct{}),
				OutletCode: p.OutletCode,
				OutletName: p.BranchName,
			}
		}

		d := retailerMap[p.RetailerGUID][p.OutletGUID]

		d.Amount += p.Amount
		d.Trans += p.Trans
		d.PosSet[p.PosID] = struct{}{} // unique POS
	}

	// 🟢 build result
	result := make(map[string][]RetailerSalesContribute)

	for retailerGUID, outletMap := range retailerMap {

		var totalAmount float64
		var totalTrans int

		for _, o := range outletMap {
			totalAmount += o.Amount
			totalTrans += o.Trans
		}

		var rows []RetailerSalesContribute

		for _, o := range outletMap {

			var amtPct float64
			var transPct float64

			if totalAmount > 0 {
				amtPct = (o.Amount / totalAmount)
			}
			if totalTrans > 0 {
				transPct = (float64(o.Trans) / float64(totalTrans))
			}

			rows = append(rows, RetailerSalesContribute{
				OutletCode:    o.OutletCode,
				OutletName:    o.OutletName,
				PosActive:     len(o.PosSet),
				Amount:        o.Amount,
				AmountPercent: amtPct,
				Trans:         o.Trans,
				TransPercent:  transPct,
			})
		}

		// 🟢 sort by Amount desc (ranking)
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].Amount > rows[j].Amount
		})

		// 🟢 assign ranking number
		for i := range rows {
			rows[i].No = fmt.Sprintf("%d", i+1)
		}

		result[retailerGUID] = rows
	}

	return result
}

func EachRetailerSummary(
	f *excelize.File,
	retailer []models.Retailer,
	pos []models.POS,
	activePOS []models.ActivePOS,
	reportMonth string,
) {
	bizSummary := BuildBizSummaryByRetailer(pos, reportMonth)
	outletSummary := BuildOutletSummaryByRetailer(pos, reportMonth)

	sort.Slice(retailer, func(i, j int) bool {
		return strings.ToLower(retailer[i].RetailerName) < strings.ToLower(retailer[j].RetailerName)
	})
	for _, r := range retailer {
		startRow := 0
		sheetName := r.RetailerName
		if len(sheetName) > 31 {
			sheetName = sheetName[:31]
		}
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "Retailer Group")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", startRow), r.RetailerGroup)
		startRow++
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "Retailer Name")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", startRow), r.RetailerName)
		startRow++
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "Report Month")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", startRow), reportMonth)
		startRow++
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "Curren Month Accumulate")
		rowCurrentMonthAccumulate := startRow
		startRow++
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "YTD Accumulate")
		rowYTDAccumulate := startRow
		startRow++
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "Total Active Outlet")
		rowTotalActiveOutlet := startRow
		startRow++
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "Total Active POS")
		rowTotalActivePOS := startRow
		// 🟢 🔥 关键：拿 summary
		bizSummaryRows := bizSummary[r.RetailerGUID]

		if len(bizSummaryRows) > 0 {
			last := bizSummaryRows[len(bizSummaryRows)-1] // 👉 最后一天

			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowCurrentMonthAccumulate), last.AccumulateAmount)
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowCurrentMonthAccumulate), fmt.Sprintf("B%d", rowCurrentMonthAccumulate), styles.GetStyle(f, "number", "right", true))
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowYTDAccumulate), last.AccumulateAmountYTD)
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowYTDAccumulate), fmt.Sprintf("B%d", rowYTDAccumulate), styles.GetStyle(f, "number", "right", true))
		} else {
			// 没数据（避免 panic）
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowCurrentMonthAccumulate), 0)
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowCurrentMonthAccumulate), fmt.Sprintf("B%d", rowCurrentMonthAccumulate), styles.GetStyle(f, "number", "right", true))
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowYTDAccumulate), 0)
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowYTDAccumulate), fmt.Sprintf("B%d", rowYTDAccumulate), styles.GetStyle(f, "number", "right", true))
		}

		OutletSummaryRows := outletSummary[r.RetailerGUID]
		if len(OutletSummaryRows) > 0 {
			last := OutletSummaryRows[len(OutletSummaryRows)-1] // 👉 最后一天

			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowTotalActiveOutlet), last.No)
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowTotalActiveOutlet), fmt.Sprintf("B%d", rowTotalActiveOutlet), styles.GetStyle(f, "number", "right", true))
		} else {
			// 没数据（避免 panic）
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowTotalActiveOutlet), 0)
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowTotalActiveOutlet), fmt.Sprintf("B%d", rowTotalActiveOutlet), styles.GetStyle(f, "number", "right", true))
		}

		// 🟢 🔥 关键：拿 activePOS 数量
		posmap := make(map[string]map[string]struct{})

		for _, a := range activePOS {
			if _, ok := posmap[a.RetailerGUID]; !ok {
				posmap[a.RetailerGUID] = make(map[string]struct{})
			}

			posmap[a.RetailerGUID][a.PosID] = struct{}{}
		}
		activeCount := len(posmap[r.RetailerGUID])

		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowTotalActivePOS), activeCount)

		startRow += 2

		startCell := fmt.Sprintf("B%d", startRow)
		endCell := fmt.Sprintf("G%d", startRow)
		headerStyle := styles.HeaderStyle(f)

		f.MergeCell(sheetName, startCell, endCell)
		f.SetCellValue(sheetName, startCell, "POS")
		f.SetCellStyle(sheetName, startCell, endCell, headerStyle)

		startCell = fmt.Sprintf("K%d", startRow)
		endCell = fmt.Sprintf("P%d", startRow)

		f.MergeCell(sheetName, startCell, endCell)
		f.SetCellValue(sheetName, startCell, "Sales Contribue")
		f.SetCellStyle(sheetName, startCell, endCell, headerStyle)

		startRow++
		GenerateExcel(f, sheetName, bizSummary[r.RetailerGUID], RetailerPOSBizSummaryRow{}, reportMonth, startRow, 1)

		GenerateExcel(f, sheetName, outletSummary[r.RetailerGUID], RetailerSalesContribute{}, reportMonth, startRow, 9)
	}

}
