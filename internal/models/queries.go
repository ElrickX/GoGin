package models

import (
	"database/sql"
	"log"
)

func QueryRetailer(db *sql.DB) []Retailer {
	query := `SELECT retailer_guid,retailer_name,

IF(b.group_code IS NULL,retailer_group,b.group_name) AS retailer_group,
	IFNULL(reg_no,'') AS reg_no,
	IFNULL(reg_name,'') AS reg_name,
	country,report_summary_included,start_date
	FROM retailer a
	LEFT JOIN ml_retailer_group b
		ON a.retailer_group=b.group_code
	WHERE is_active=1 AND is_monitor=1
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var retailers []Retailer
	for rows.Next() {
		var r Retailer
		if err := rows.Scan(&r.RetailerGUID, &r.RetailerName, &r.RetailerGroup, &r.RegNo, &r.RegName, &r.Country, &r.ReportSummaryIncluded, &r.StartDate); err != nil {
			log.Fatal(err)
		}
		retailers = append(retailers, r)
	}

	return retailers
}

func QueryPOS(db *sql.DB, bizdatefrom, bizdateto string) []POS {
	query := `
	SELECT a.bizdate,a.retailer_guid,a.retailer_name,a.outlet_guid,a.outlet_code,
	ifnull(b.branch_name,'') AS branch_name,a.pos_id,a.amount,a.trans
	FROM pos_info a
	inner join cp_set_branch b
		on a.retailer_guid=b.retailer_guid and a.outlet_code=b.outlet_code
	WHERE a.bizdate BETWEEN ? AND ?
	`

	rows, err := db.Query(query, bizdatefrom, bizdateto)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var pos []POS
	for rows.Next() {
		var p POS
		if err := rows.Scan(&p.BizDate, &p.RetailerGUID, &p.RetailerName, &p.OutletGUID, &p.OutletCode, &p.BranchName, &p.PosID, &p.Amount, &p.Trans); err != nil {
			log.Fatal(err)
		}
		pos = append(pos, p)
	}
	return pos
}

func QueryActivePOS(db *sql.DB, bizdatefrom, bizdateto string) []ActivePOS {
	query := `
	SELECT retailer_guid,outlet_guid,outlet_code,bizdate,pos_id
	FROM ts_active_pos
	WHERE bizdate BETWEEN ? AND ?`

	rows, err := db.Query(query, bizdatefrom, bizdateto)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var activePOS []ActivePOS
	for rows.Next() {
		var ap ActivePOS
		if err := rows.Scan(&ap.RetailerGUID, &ap.OutletGUID, &ap.OutletCode, &ap.BizDate, &ap.PosID); err != nil {
			log.Fatal(err)
		}
		activePOS = append(activePOS, ap)
	}

	return activePOS
}

func QueryGR(db *sql.DB, monthfrom, monthto string) []GR {
	query := `
	SELECT retailer_guid,COUNT(sup_code) AS sup_count
	FROM
	(
	SELECT retailer_guid,sup_code
	FROM gr_info
	WHERE PERIOD BETWEEN ? AND ?
	GROUP BY retailer_guid,sup_code
	)a
	GROUP BY retailer_guid
	`

	rows, err := db.Query(query, monthfrom, monthto)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var grs []GR
	for rows.Next() {
		var gr GR
		if err := rows.Scan(&gr.RetailerGUID, &gr.SupCount); err != nil {
			log.Fatal(err)
		}
		grs = append(grs, gr)
	}

	return grs
}

func QuerySup(db *sql.DB, monthfrom, monthto string) []Sup {
	query := `
	SELECT a.retailer_guid,COUNT(sup_code) AS sup_count
	FROM
	(
	SELECT a.retailer_guid,sup_code
	FROM gr_info a
	INNER JOIN supplier_info b
		ON a.retailer_guid=b.retailer_guid AND a.sup_code=b.code
	WHERE PERIOD BETWEEN ? AND ?
	GROUP BY a.retailer_guid,sup_code
	)a
	GROUP BY retailer_guid
		`

	rows, err := db.Query(query, monthfrom, monthto)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var sups []Sup
	for rows.Next() {
		var s Sup
		if err := rows.Scan(&s.RetailerGUID, &s.SupCount); err != nil {
			log.Fatal(err)
		}
		sups = append(sups, s)
	}

	return sups
}

func QueryRetailerPOSPayment(db *sql.DB, bizdatefrom, bizdateto string) []RetailerPOSPayment {
	query := `
	SELECT
	a.retailer_guid AS retailer_guid,
	a.retailer_name AS retailer_name,
	IFNULL(b.period_code, '')   AS period_code,
	IFNULL(b.total, 0)         AS total,
	IFNULL(b.Cash, 0)          AS cash,
	IFNULL(b.Credit_Card, 0)   AS credit_card,
	IFNULL(b.Debit_Card, 0)    AS debit_card,
	IFNULL(b.Other, 0)         AS other,
	IFNULL(b.E_Wallet, 0)      AS e_wallet,
	IFNULL(b.G_Wallet, 0)      AS g_wallet,
	IFNULL(b.MyKasih, 0)       AS mykasih,
	IFNULL(ROUND(((b.MyKasih / b.total) * 100),1), 0) AS MyKasih_Contribute_Pct,
	IFNULL(a.Remark, '')       AS remark
	FROM retailer a
	LEFT JOIN (SELECT
					a.retailer_guid  AS retailer_guid,
					LEFT(a.bizdate,7) AS period_code,
					ROUND(SUM(a.actualvalue),2) AS total,
					ROUND(SUM((CASE WHEN (a.paytype = 'CASH') THEN a.actualvalue ELSE 0 END)),2) AS Cash,
					ROUND(SUM((CASE WHEN (a.paytype = 'CREDIT CARD') THEN a.actualvalue ELSE 0 END)),2) AS Credit_Card,
					ROUND(SUM((CASE WHEN (a.paytype = 'ATM CARD') THEN a.actualvalue ELSE 0 END)),2) AS Debit_Card,
					ROUND(SUM((CASE WHEN (UPPER(a.paytype) NOT IN ('CASH','CREDIT CARD','ATM CARD','E-WALLET','G-WALLET')) THEN a.actualvalue ELSE 0 END)),2) AS Other,
					ROUND(SUM((CASE WHEN (UPPER(a.paytype) = 'E-WALLET') THEN a.actualvalue ELSE 0 END)),2) AS E_Wallet,
					ROUND(SUM((CASE WHEN (UPPER(a.paytype) = 'G-WALLET') THEN a.actualvalue ELSE 0 END)),2) AS G_Wallet,
					ROUND(SUM((CASE WHEN ((UPPER(a.paytype) <> 'CASH') AND (UPPER(a.cardtype) LIKE '%KASIH')) THEN a.actualvalue ELSE 0 END)),2) AS MyKasih
				FROM pos_payment_info a
				WHERE bizdate BETWEEN ? AND ?
				GROUP BY a.retailer_guid,period_code) b
		ON a.retailer_guid = b.retailer_guid
	WHERE a.is_Active = 1
		AND a.is_Monitor = 1
		and country='MALAYSIA'
	ORDER BY total DESC`

	rows, err := db.Query(query, bizdatefrom, bizdateto)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var retailerPOSPayment []RetailerPOSPayment
	for rows.Next() {
		var ap RetailerPOSPayment
		if err := rows.Scan(&ap.RetailerGUID, &ap.RetailerName, &ap.PeriodCode, &ap.Total, &ap.Cash, &ap.CreditCard, &ap.DebitCard, &ap.Other, &ap.EWallet, &ap.GWallet, &ap.MyKasih, &ap.MyKasihContributePct, &ap.Remark); err != nil {
			log.Fatal(err)
		}
		retailerPOSPayment = append(retailerPOSPayment, ap)
	}

	return retailerPOSPayment
}
