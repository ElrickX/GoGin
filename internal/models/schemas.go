package models

import "time"

type Retailer struct {
	RetailerGUID          string
	RetailerName          string
	RetailerGroup         string
	RegNo                 string
	RegName               string
	Country               string
	ReportSummaryIncluded bool
	StartDate             time.Time
}

type POS struct {
	BizDate      time.Time
	RetailerGUID string
	RetailerName string
	OutletGUID   string
	OutletCode   string
	BranchName   string
	PosID        string
	Amount       float64
	Trans        int
}

type ActivePOS struct {
	RetailerGUID string
	OutletGUID   string
	OutletCode   string
	BizDate      string
	PosID        string
}

type GR struct {
	RetailerGUID string
	SupCount     int
}

type Sup struct {
	RetailerGUID string
	SupCount     int
}

type RetailerPOSPayment struct {
	RetailerGUID         string  `excel:"Retailer GUID" width:"20"`
	RetailerName         string  `excel:"Retailer Name" width:"25"`
	PeriodCode           string  `excel:"Period Code" width:"15"`
	Total                float64 `excel:"Total" width:"18" format:"currency"`
	Cash                 float64 `excel:"Cash" width:"18" format:"currency"`
	CreditCard           float64 `excel:"Credit Card" width:"18" format:"currency"`
	DebitCard            float64 `excel:"Debit Card" width:"18" format:"currency"`
	Other                float64 `excel:"Other" width:"18" format:"currency"`
	EWallet              float64 `excel:"E-Wallet" width:"18" format:"currency"`
	GWallet              float64 `excel:"G-Wallet" width:"18" format:"currency"`
	MyKasih              float64 `excel:"MyKasih" width:"18" format:"currency"`
	MyKasihContributePct float64 `excel:"MyKasih Contribute %" width:"18" format:"percent"`
	Remark               string  `excel:"Remark" width:"25"`
}
