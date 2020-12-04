package main

import (
	"log"

	"github.com/tealeg/xlsx/v3"
)

type discountItem struct {
	sku         string
	location    string
	discountVal map[string]string //percent group, amount off
}

func main() {
	// open an existing file
	file, err := xlsx.OpenFile("./data/discount_item_list.xlsx")

	// go through each row,
	// get first cell (A) this is sku
	// if any subsequent cell has a value, add that cells discount group to the tranlated xlsx document

	if err != nil {
		log.Fatalf("FileOpenErr: %+v", err)
	}

	// file now contains a reference to the workbook
	// show all the sheets in the workbook
	// log.Println("Sheets in this file:")

	// for i, sh := range file.Sheets {
	// 	log.Println(i, sh.Name)
	// }
	// log.Println("----")

	// this creates a new file, we might be able to use this to make the magento file
	// wb := xlsx.NewFile();
	sh, ok := file.Sheet["Sheet1"]

	if !ok {
		log.Fatalf("SheetDoesNotExistErr")
		return
	}

	// get the coordinates of the cell, if its
	// 1st in a cell its the sku
	// 2nd in a cell its 10 %
	// 3rd in a cell its 20 %
	// 4th in a cell its 25 %
	// 5th in a cell its 30 %
	// 6th in a cell its 35 %
	// 7th in a cell its 38 %
	// 9th in a cell its 40 %
	// make a map for this string to convert
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

	// row, err := sh.Row(1)

	// let's so something with the row ...
	// log.Println("Max row is", sh.MaxRow)

	// for each row
	// we want to get each cell value
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

				di.discountVal[code] = value

				// if we hit h aka we know we can run our translate function and move on to our next sku
				// because this is our final bound aka our final customer group
				if code == "40%" {
					translateRowToMagentoForm(*di)
				}
			}
			// log.Printf("%+v", *di)
			return err

		})

	})
}

func translateRowToMagentoForm(di discountItem) {
	// if all items are 0, don't add to new xlsx file
	cm := len(di.discountVal)
	count := 1
	addFlag := true
	for _, v := range di.discountVal {
		if v == "0" {
			count++
		}
		if cm == count {
			addFlag = false
		}
	}

	if addFlag {
		log.Printf("%+v", di)
		// create the magento sheet

		// first create a new row with all the required values
		// then insert each item as a cell in a row
	}
}
