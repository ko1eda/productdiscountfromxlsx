package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/tealeg/xlsx/v3"
)

// i know records should be a passed by reference from main func but for sake of time, global state
var (
	records = [][]string{
		{"sku", "tier_price_website", "tier_price_customer_group", "tier_price_qty", "tier_price", "tier_price_value_type"},
	}
)

type discountItem struct {
	sku         string
	discountVal map[string]string //percent group name ex 25%, actual discount amount
}

// Quick program to convert a companies xsls discount spreadsheet into a magento 2 readible csv for import
func main() {
	// open an existing file
	file, err := xlsx.OpenFile("./data/discount_item_list.xlsx")

	if err != nil {
		log.Fatalf("FileOpenErr: %+v", err)
	}

	// close our sheet when we're done
	sh, ok := file.Sheet["Sheet1"]
	defer sh.Close()

	if !ok {
		log.Fatalf("SheetDoesNotExistErr")
		return
	}

	// maps cell letters to discount code
	cellToDiscount := map[string]string{
		"A": "SKU",
		"B": "10%",
		"C": "20%",
		"D": "25%",
		"E": "30%",
		"F": "35%",
		"G": "38%",
		"H": "40%",
	}

	// for each row in our default sheet we make a discount item
	// when we hit a sku cell, we add data to our item
	// when we hit a data sku we add data to our map related to what cell we're in
	sh.ForEachRow(func(r *xlsx.Row) error {
		di := &discountItem{}
		return r.ForEachCell(func(c *xlsx.Cell) error {
			coor := xlsx.GetCellIDStringFromCoords(c.GetCoordinates())
			switch code := cellToDiscount[coor[0:1]]; {
			case code == "SKU":
				value, err := c.FormattedValue()

				if err != nil {
					log.Fatalln(err.Error())
				}

				di = &discountItem{
					sku:         value,
					discountVal: make(map[string]string, 7),
				}
			default:
				value, err := c.FormattedValue()

				if err != nil {
					log.Fatalln(err.Error())
				}

				// to remove 0 values in general
				// if (val != 0) then run di.discountVal[code] = val
				di.discountVal[code] = value

				// if we hit h aka 40% we know we can run our translate function and move on to our next sku
				// because this is our final bound aka our final customer group
				if code == "40%" && di.sku != "" {
					translateRowToMagentoForm(*di)
				}
			}
			return err
		})
	})

	// here we create our exported file
	// and write all our translated data from the xlsx into the csv
	csvFile, err := os.Create("./data/export_advanced_pricing.csv")

	if err != nil {
		log.Fatalf("CSVFileCreateErr: %+v", err)
	}

	defer csvFile.Close()

	// we create a new csvwriter and pass our csv file
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// for each slice in records we write it to our csv file
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// check the writer for errors before we close
	if err := writer.Error(); err != nil {
		log.Fatalf("CSVWriterHasErrorErr: %+v", err)
	}
}

// This function takes a discount item, takes its values and adds it to our records slice
// which is a csv readible format, which we use to write to csv file at the end of the program
func translateRowToMagentoForm(di discountItem) {
	cm := len(di.discountVal)
	count := 1
	addFlag := true

	// if all items are 0, don't add to new xlsx file
	for _, v := range di.discountVal {
		if v == "0" {
			count++
		}
		if cm == count {
			addFlag = false
		}
	}

	// if all is good, for each item we have in our map we add a row to our eventual converted csv
	// each map entry represents one row on the csv, and each row on the csv reprsents one item and its price
	// at a given discount group EX: sku 909, discount group 40% tier price 30 --> this may be the actual percentage
	// off for certainm items even if they are in 40% group -- weird business logic
	if addFlag {
		for k, v := range di.discountVal {
			records = append(records, []string{di.sku, "All Websites [USD]", k, "1.000", v, "Discount"})
		}
	}
}
