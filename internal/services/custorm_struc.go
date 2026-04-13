package services

import "time"

type RetailerSummaryRow struct {
	GUID              string  `excel:"GUID" width:"20"`
	StartDate         string  `excel:"Start Date" width:"15" align:"center"`
	RegNo             string  `excel:"Reg No" width:"18"`
	Category          string  `excel:"Category" width:"10" align:"center"`
	Name              string  `excel:"Retailer Name" width:"25"`
	Amount            float64 `excel:"Amount" width:"18" format:"currency"`
	Trans             int     `excel:"Trans" width:"12" format:"number"`
	YTDAmount         float64 `excel:"YTD Amount" width:"18" format:"currency"`
	AnnualizedRevenue float64 `excel:"Annualized Revenue" width:"20" format:"currency"`
	SupCount          int     `excel:"Supplier Count" width:"15" format:"number"`
	Remark            string  `excel:"Remark" width:"25"`
}

type RetailerPOSBizSummaryRow struct {
	BizDate             time.Time `excel:"Biz Date" width:"25" align:"center" format:"date"`
	Amount              float64   `excel:"Amount" width:"18" format:"currency"`
	AccumulateAmount    float64   `excel:"MTD Amount" width:"18" format:"currency"`
	AccumulateAmountYTD float64   `excel:"YTD Amount" width:"18" format:"currency"`
	Trans               int       `excel:"Trans" width:"12" format:"number"`
	AccumulateTrans     int       `excel:"MTD Trans" width:"15" format:"number"`
	AccumulateTransYTD  int       `excel:"YTD Trans" width:"15" format:"number"`
}

type RetailerSalesContribute struct {
	No            string  `excel:"No" width:"6" align:"center"`
	OutletCode    string  `excel:"Outlet Code" width:"15"`
	OutletName    string  `excel:"Outlet Name" width:"25"`
	PosActive     int     `excel:"Active POS" width:"12" format:"number"`
	Amount        float64 `excel:"Amount" width:"18" format:"currency"`
	AmountPercent float64 `excel:"Amount %" width:"12" format:"percent"`
	Trans         int     `excel:"Trans" width:"12" format:"number"`
	TransPercent  float64 `excel:"Trans %" width:"12" format:"percent"`
}
