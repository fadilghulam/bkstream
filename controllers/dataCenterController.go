package controllers

import (
	db "bkstream/config"
	"bkstream/helpers"
	"encoding/csv"
	"fmt"
	"io"
	"log"
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
