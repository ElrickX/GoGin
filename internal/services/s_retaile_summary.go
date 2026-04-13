package services

import (
	"demo/internal/models"
	"fmt"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"
)

func BuildRetailerSummary(
	f *excelize.File,
	retailer []models.Retailer,
	supplier []models.Sup,
	pos []models.POS,
	reportMonth string,
	activePOS []models.ActivePOS,
) {
	startRow := 1
	f.NewSheet("Retailer Summary")

	f.SetCellValue("Retailer Summary", fmt.Sprintf("A%d", startRow), "Category Legend")
	startRow++
	f.SetCellValue("Retailer Summary", fmt.Sprintf("A%d", startRow), "A: Annual revenue > 1B or active POS >= 401")
	startRow++
	f.SetCellValue("Retailer Summary", fmt.Sprintf("A%d", startRow), "B: Annual revenue > 500M-1B or active POS 251-400")
	startRow++
	f.SetCellValue("Retailer Summary", fmt.Sprintf("A%d", startRow), "C: Annual revenue > 100M-500M or active POS 100-250")
	startRow++
	f.SetCellValue("Retailer Summary", fmt.Sprintf("A%d", startRow), "D: Annual revenue < 100M and active POS < 100")
	startRow++
	f.SetCellValue("Retailer Summary", fmt.Sprintf("A%d", startRow), "Note: This table may not include indonesia currency")
	startRow++
	f.SetCellValue("Retailer Summary", fmt.Sprintf("A%d", startRow), "Note: This table is by Retailer Group Grouping")

	startRow += 2
	bizSummary := BuildBizSummaryByRetailer(pos, reportMonth)

	retailerGroupMap := make(map[string][]models.Retailer)
	for _, r := range retailer {
		if r.Country != "MALAYSIA" {
			continue
		}
		retailerGroupMap[r.RetailerGroup] = append(retailerGroupMap[r.RetailerGroup], r)
	}
	result := make(map[string][]RetailerSummaryRow)

	for group, retailers := range retailerGroupMap {
		aggAmount := 0.0
		aggYtdAmount := 0.0
		aggTrans := 0
		supCount := 0
		activeCount := 0

		for _, r := range retailers {
			bizSummaryRows := bizSummary[r.RetailerGUID]
			if len(bizSummaryRows) > 0 {
				last := bizSummaryRows[len(bizSummaryRows)-1] // 👉 最后一天
				aggAmount += last.AccumulateAmount
				aggTrans += last.AccumulateTrans
				aggYtdAmount += last.AccumulateAmountYTD

				for _, s := range supplier {
					if s.RetailerGUID == r.RetailerGUID {
						supCount = s.SupCount
					}
				}
			}
			activeCount += getActiveCount(activePOS, r.RetailerGUID)

		}
		ytdDays := getYTDDays(reportMonth)

		ranking, annualizedrevenue := getCategory(aggAmount, activeCount, ytdDays)

		result[group] = append(result[group], RetailerSummaryRow{
			GUID:              retailers[0].RetailerGUID,
			StartDate:         retailers[0].StartDate.Format("2006-01-01"),
			RegNo:             retailers[0].RegNo,
			Category:          ranking,
			Name:              retailers[0].RetailerGroup,
			Amount:            aggAmount,
			Trans:             aggTrans,
			YTDAmount:         aggYtdAmount,
			AnnualizedRevenue: annualizedrevenue,
			SupCount:          supCount,
			Remark:            "",
		})
	}
	var flat []RetailerSummaryRow

	for _, rows := range result {
		flat = append(flat, rows...)
	}
	sort.Slice(flat, func(i, j int) bool {
		return flat[i].AnnualizedRevenue > flat[j].AnnualizedRevenue
	})
	GenerateExcel(f, "Retailer Summary", flat, RetailerSummaryRow{}, reportMonth, startRow, 1)
}

func getYTDDays(reportMonth string) int {
	ytdDays := 365
	if t, err := time.Parse("2006-01", reportMonth); err == nil {
		// 下个月第一天
		nextMonth := t.AddDate(0, 1, 0)

		// 月底
		endOfMonth := nextMonth.AddDate(0, 0, -1)

		// 年累计天数
		ytdDays = endOfMonth.YearDay()

	}
	return ytdDays
}

func getActiveCount(activePOS []models.ActivePOS, retailerGUID string) int {
	posmap := make(map[string]map[string]struct{})
	for _, a := range activePOS {
		if _, ok := posmap[a.RetailerGUID]; !ok {
			posmap[a.RetailerGUID] = make(map[string]struct{})
		}

		posmap[a.RetailerGUID][a.PosID] = struct{}{}
	}
	return len(posmap[retailerGUID])
}
