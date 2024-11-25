package controllers

import (
	db "bkstream/config"
	"bkstream/helpers"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gorm.io/gorm"
)

// func FetchCustomerDataAsCSVDynamic(db *gorm.DB, filePath string) error {
// 	var results []map[string]interface{}

// 	// Query the customer table dynamically
// 	if err := db.Model(&structs.Customer{}).Where("salesman_id = 781").Find(&results).Error; err != nil {
// 		return err
// 	}

// 	if len(results) == 0 {
// 		return fmt.Errorf("no data found")
// 	}

// 	// Create a CSV file
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Dynamically get the column names (keys from the first result map)
// 	header := make([]string, 0, len(results[0]))
// 	for key := range results[0] {
// 		header = append(header, key)
// 	}

// 	// Write the header
// 	if err := writer.Write(header); err != nil {
// 		return err
// 	}

// 	// Write the data rows
// 	for _, row := range results {
// 		record := make([]string, len(header))
// 		for i, colName := range header {
// 			// Convert the value to a string
// 			value := row[colName]
// 			record[i] = fmt.Sprintf("%v", value)
// 		}
// 		if err := writer.Write(record); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func FetchCustomerDataAsCSVDownload(c *fiber.Ctx, db *gorm.DB) error {
	start := time.Now()

	// var results []map[string]interface{}
	var results []*orderedmap.OrderedMap[string, interface{}]

	// Query the customer table dynamically
	// if err := db.Model(&structs.Customer{}).Where("salesman_id = 781").Find(&results).Error; err != nil {
	// 	return err
	// }

	// finalQuery := `SELECT *, id||'' as id FROM customer WHERE salesman_id = 781`
	finalQuery := `WITH s AS(
						SELECT d.id 
						FROM salesman d
						WHERE d.branch_id = 111
					)
					SELECT d.*, d.id||'' as id
					FROM s
					JOIN customer d
					ON d.salesman_id = s.id
					WHERE d.is_verifikasi != -1`

	results, err := helpers.ExecuteGORMQueryOrdered2(finalQuery)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Gagal execute query",
			Success: false,
		})
	}

	// labels := []string{}
	// for pair := results[0].Oldest(); pair != nil; pair = pair.Next() {
	// 	if pair.Key != "product" {
	// 		labels = append(labels, pair.Key)
	// 	}
	// }

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).SendString("No data found")
	}

	// Create an in-memory pipe to write CSV data to
	reader, writer := io.Pipe()
	csvWriter := csv.NewWriter(writer)

	// Run the CSV writing in a goroutine
	go func() {
		defer writer.Close()
		defer csvWriter.Flush()

		// Dynamically get the column names (keys from the first result map)
		// header := make([]string, 0, results[0].Len())
		// for key := range results[0] {
		// 	header = append(header, key)
		// }
		// fmt.Println(results)
		// fmt.Println(results[0])
		header := []string{}
		for pair := results[0].Oldest(); pair != nil; pair = pair.Next() {
			// fmt.Println(pair.Key, pair.Value)
			header = append(header, pair.Key)
		}

		// fmt.Println("Header:", header)

		// Write the header
		if err := csvWriter.Write(header); err != nil {
			log.Println("Error writing CSV header:", err)
			return
		}

		// Write the data rows
		// for _, row := range results {
		// 	record := make([]string, len(header))
		// 	for i, colName := range header {
		// 		// Convert the value to a string
		// 		value := row[colName]
		// 		record[i] = fmt.Sprintf("%v", value)
		// 	}
		// 	if err := csvWriter.Write(record); err != nil {
		// 		log.Println("Error writing CSV row:", err)
		// 		return
		// 	}
		// }

		for _, row := range results {
			record := make([]string, len(header))
			for i, colName := range header {
				// Convert the value to a string
				value, ok := row.Get(colName)
				if !ok {
					// Handle the case where the key does not exist in the OrderedMap
					// For example, you can set the value to an empty string or skip the row
					value = ""
				}
				// record[i] = fmt.Sprintf("%v", value)
				if value == nil {
					record[i] = "null"
				} else {
					record[i] = fmt.Sprintf("%v", value)
				}
			}
			if err := csvWriter.Write(record); err != nil {
				log.Println("Error writing CSV row:", err)
				return
			}
		}
	}()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Elapsed time:", elapsed)

	// Set the content headers for CSV download
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="customer_data.csv"`)

	// Stream the CSV file to the client
	return c.SendStream(reader)
}

func GetDataCustomer(c *fiber.Ctx) error {
	// if err := FetchCustomerDataAsCSVDynamic(db.DB, "customer_data.csv"); err != nil {
	// 	log.Fatal(err)
	// }

	if err := FetchCustomerDataAsCSVDownload(c, db.DB); err != nil {
		log.Fatal(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func FetchDashboardOmzet(qWhereBranchId string, qWhereBranchHolderId string, date string, mode string) (map[string]interface{}, error) {

	if date == "" {
		date = "NOW()"
	}

	var qWhereDatePenjualan, qWhereDatePenjualanLM, qWhereDatePengembalian, qWhereDatePengembalianLM string
	if mode != "" {
		qWhereDatePenjualan = `p.tanggal_penjualan
									BETWEEN NOW() - INTERVAL '5 minute' AND NOW()
								{{.QWherePbranchId}}`

		qWhereDatePenjualanLM = `p.tanggal_penjualan
									BETWEEN NOW() - INTERVAL '5 minute' AND NOW()	
								{{.QWherePbranchId}}`

		qWhereDatePengembalian = `p.tanggal_pengembalian
									BETWEEN NOW() - INTERVAL '5 minute' AND NOW()
								{{.QWherePbranchId}}`

		qWhereDatePengembalianLM = `p.tanggal_pengembalian
									 BETWEEN NOW() - INTERVAL '5 minute' AND NOW()
								{{.QWherePbranchId}}`

	} else {
		qWhereDatePenjualan = `DATE(p.tanggal_penjualan)
									BETWEEN DATE(date_trunc('month', {{.QDate}}))
											AND DATE(date_trunc('month', {{.QDate}}) + '1 month'::interval - '1 day'::interval)
								{{.QWherePbranchId}}`

		qWhereDatePenjualanLM = `DATE(p.tanggal_penjualan)
									BETWEEN DATE((date_trunc('month', {{.QDate}}) - '1 month'::interval))
											AND DATE((date_trunc('month', {{.QDate}})- '1 month'::interval) + '1 month'::interval - '1 day'::interval)

								{{.QWherePbranchId}}`

		qWhereDatePengembalian = `DATE(p.tanggal_pengembalian)
									BETWEEN DATE(date_trunc('month', {{.QDate}}))
											AND DATE(date_trunc('month', {{.QDate}}) + '1 month'::interval - '1 day'::interval)
								{{.QWherePbranchId}}`

		qWhereDatePengembalianLM = `DATE(p.tanggal_pengembalian)
									BETWEEN DATE((date_trunc('month', {{.QDate}}) - '1 month'::interval))
											AND DATE((date_trunc('month', {{.QDate}})- '1 month'::interval) + '1 month'::interval - '1 day'::interval)
								{{.QWherePbranchId}}`
	}

	templateReplaceQuery := map[string]interface{}{
		"QWhereDatePenjualan":      qWhereDatePenjualan,
		"QWhereDatePenjualanLM":    qWhereDatePenjualanLM,
		"QWhereDatePengembalian":   qWhereDatePengembalian,
		"QWhereDatePengembalianLM": qWhereDatePengembalianLM,
		"QWherePbranchId":          qWhereBranchId,
		"QWhereBranchHolderId":     qWhereBranchHolderId,
		"QDate":                    date,
	}

	queryGetOmzet := `WITH penjualan_this_month as (
							SELECT SUM((pd.harga - pd.diskon) * pd.jumlah) as total_penjualan,
											SUM(pd.jumlah) as total_pack,
											SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
											SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM penjualan p
							JOIN penjualan_detail pd
								ON p.id = pd.penjualan_id
							WHERE 
								{{.QWhereDatePenjualan}}
						), penjualan_last_month as (
							SELECT SUM((pd.harga - pd.diskon) * pd.jumlah) as total_penjualan, 
										SUM(pd.jumlah) as total_pack,
										SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
										SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM penjualan p
							JOIN penjualan_detail pd
								ON p.id = pd.penjualan_id
							WHERE 
								{{.QWhereDatePenjualanLM}}
						), pengembalian_this_month as (
							SELECT SUM(pd.harga * pd.jumlah) as total_penjualan,
										SUM(pd.jumlah) as total_pack,
										SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
										SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM pengembalian p
							JOIN pengembalian_detail pd
								ON p.id = pd.pengembalian_id
							WHERE 
								{{.QWhereDatePengembalian}}
								
						), pengembalian_last_month as (
							SELECT SUM(pd.harga * pd.jumlah) as total_penjualan,
										SUM(pd.jumlah) as total_pack,
										SUM(pd.jumlah) FILTER (WHERE pd.harga <> 0) as total_pack_omzet,
										SUM(pd.jumlah) FILTER (WHERE pd.harga = 0) as total_pack_bonus
							FROM pengembalian p
							JOIN pengembalian_detail pd
								ON p.id = pd.pengembalian_id
							WHERE 
								{{.QWhereDatePengembalianLM}}
						)

						SELECT data.otm as this_month, 
										data.olm as last_month,
										ROUND((((data.otm - data.olm)  / CASE WHEN data.olm = 0 THEN 1 ELSE data.olm END) * 100)::numeric,2) as growth
						FROM (
							SELECT COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'ptm'),0) - 
											COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'pgtm'),0) as otm,
											COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'plm'),0) - 
											COALESCE(MAX(sq.total_penjualan) FILTER (WHERE sq.flag = 'pglm'),0) as olm
							FROM (
								SELECT *, 'ptm' as flag
								FROM penjualan_this_month ptm

								UNION ALL

								SELECT *, 'plm' as flag
								FROM penjualan_last_month plm

								UNION ALL

								SELECT *, 'pgtm' as flag
								FROM pengembalian_this_month pgtm

								UNION ALL

								SELECT *, 'pglm' as flag
								FROM pengembalian_last_month pglm
							) sq
						) data
						`

	query1, err := helpers.PrepareQuery(queryGetOmzet, templateReplaceQuery)
	query1Fix, err := helpers.PrepareQuery(query1, templateReplaceQuery)

	// fmt.Println(query1Fix)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	dataOmzet, err := helpers.ExecuteQuery(query1Fix)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	queryGetPiutang := `SELECT SUM(data.pitm) as piutang_this_month, COALESCE(SUM(data.pilm),0) as piutang_last_month, ROUND((((SUM(data.pitm) - COALESCE(SUM(data.pilm),0))  / CASE WHEN SUM(data.pilm) = 0 THEN 1 ELSE SUM(data.pilm) END) * 100)::numeric,2) as growth
						FROM (
						SELECT
							SUM( CASE WHEN (DATE_PART('day',{{.QDate}}::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
								THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pitm, 
		--                     SUM( CASE WHEN (DATE_PART('day',{{.QDate}}::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
		--                          THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pilm,
												0 as pilm,
							SUM(total_piutang-COALESCE(ppd.nominal,0)) AS total
								
							FROM
								piutang pi
							LEFT JOIN
								(SELECT piutang_id, SUM(nominal) as nominal 
								FROM pembayaran_piutang pp
								JOIN pembayaran_piutang_detail ppd 
								ON ppd.pembayaran_piutang_id = pp.id
								WHERE DATE(pp.tanggal_pembayaran) <= {{.QDate}}
								GROUP BY piutang_id) ppd 
								ON ppd.piutang_id = pi.id
							JOIN penjualan p
							ON p.id = pi.penjualan_id
							JOIN customer c
							ON c.id = p.customer_id  AND c.is_kasus IN ( 0 )
							JOIN salesman se
							ON se.id = p.salesman_id

							LEFT JOIN
							area ae
							ON ae.id = p.area_id 
							LEFT JOIN branch be
							ON be.id = p.branch_id 
							LEFT JOIN rayon re 
							ON re.id = p.rayon_id 
							LEFT JOIN sr sre
							ON sre.id = p.sr_id 

							JOIN salesman sh
							ON c.salesman_id = sh.id
							LEFT JOIN area ah
							ON ah.id = ANY(sh.area_id) AND ah.id = p.area_id
							LEFT JOIN branch bh
							ON sh.branch_id = bh.id 
							LEFT JOIN rayon rh
							ON rh.id = bh.rayon_id 
							LEFT JOIN sr srh 
							ON srh.id = rh.sr_id 

							WHERE
								DATE(pi.tanggal_piutang) <= {{.QDate}} 
								AND c.is_kasus IN ( 0 )
								AND rh.id <> 501
								{{.QWhereBranchHolderId}}

							UNION ALL

							SELECT
								0 as pitm,
								SUM( CASE WHEN (DATE_PART('day',({{.QDate}}-'1 month'::interval)::timestamp -DATE(pi.tanggal_piutang)::timestamp) > 90 )
									THEN total_piutang-COALESCE(ppd.nominal,0) ELSE 0 END ) AS pilm,
								SUM(total_piutang-COALESCE(ppd.nominal,0)) AS total
								
							FROM
								piutang pi
							LEFT JOIN
								(SELECT piutang_id, SUM(nominal) as nominal 
								FROM pembayaran_piutang pp
								JOIN pembayaran_piutang_detail ppd 
								ON ppd.pembayaran_piutang_id = pp.id
								WHERE DATE(pp.tanggal_pembayaran) <= ({{.QDate}}-'1 month'::interval)
								GROUP BY piutang_id) ppd 
								ON ppd.piutang_id = pi.id
							JOIN penjualan p
							ON p.id = pi.penjualan_id
							JOIN customer c
							ON c.id = p.customer_id  AND c.is_kasus IN ( 0 )
							JOIN salesman se
							ON se.id = p.salesman_id

							LEFT JOIN
							area ae
							ON ae.id = p.area_id 
							LEFT JOIN branch be
							ON be.id = p.branch_id 
							LEFT JOIN rayon re 
							ON re.id = p.rayon_id 
							LEFT JOIN sr sre
							ON sre.id = p.sr_id 

							JOIN salesman sh
							ON c.salesman_id = sh.id
							LEFT JOIN area ah
							ON ah.id = ANY(sh.area_id) AND ah.id = p.area_id
							LEFT JOIN branch bh
							ON sh.branch_id = bh.id 
							LEFT JOIN rayon rh
							ON rh.id = bh.rayon_id 
							LEFT JOIN sr srh 
							ON srh.id = rh.sr_id 

							WHERE
								DATE(pi.tanggal_piutang) <= ({{.QDate}}-'1 month'::interval)
								AND c.is_kasus IN ( 0 )
								AND rh.id <> 501
								{{.QWhereBranchHolderId}}
						) data`

	query2, err := helpers.PrepareQuery(queryGetPiutang, templateReplaceQuery)

	// fmt.Println(query2)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	dataPiutang, err := helpers.ExecuteQuery(query2)

	if err != nil {
		fmt.Println("piutang ", err.Error())
		return nil, err
	}

	returnData := make(map[string]interface{})

	if len(dataOmzet) > 0 {
		returnData["omzet"] = dataOmzet[0]
	}

	if len(dataPiutang) > 0 {
		returnData["receiveable"] = dataPiutang[0]
	}

	return returnData, nil
}

func GetDashboardOmzet(c *fiber.Ctx) error {

	start := time.Now()

	branchId := helpers.ParamArray(c.Context().QueryArgs().PeekMulti("branchId[]"))
	date := c.Query("date")

	if date == "" {
		date = "CURRENT_DATE"
	} else {
		date = " DATE('" + date + "') "
	}

	var qWhereBranchId, QWhereBranchHolderId string
	if len(branchId) > 0 {
		qWhereBranchId = " AND p.branch_id IN (" + strings.Join(branchId, ",") + ")"
		QWhereBranchHolderId = " AND bh.id IN (" + strings.Join(branchId, ",") + ")"
		// qOnBranchId = qWhereBranchId
	}

	datas, err := FetchDashboardOmzet(qWhereBranchId, QWhereBranchHolderId, date, "a")
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Terjadi kesalahan ketika mengambil data",
			Success: false,
		})
	}

	elapsed := time.Since(start)

	type Response struct {
		Message string        `json:"message"`
		Success bool          `json:"success"`
		Data    interface{}   `json:"data"`
		Elapsed time.Duration `json:"elapsed"`
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Message: "Success",
		Success: true,
		Data:    datas,
		Elapsed: time.Duration(elapsed.Seconds()),
	})

}
