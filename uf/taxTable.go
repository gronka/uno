package uf

import (
	"encoding/csv"
	"os"
)

var fullTaxCsv string = "/fridayy/taxes/full.csv"

var averageSalesTax string = "0.0635"

type TaxRateRow struct {
	Postal string
	Rate   string
}

type TaxTable map[string]TaxRateRow

func createTaxTable() TaxTable {
	taxRateTable := make(map[string]TaxRateRow, 0)

	file, err := os.Open(fullTaxCsv)
	if err != nil {
		Fatal(err)
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	// data is of type [][]string
	data, err := csvReader.ReadAll()
	if err != nil {
		Fatal(err)
		panic(err)
	}

	for i, line := range data {
		if i > 0 { // omit header line
			var trr TaxRateRow
			for j, field := range line {
				switch j {
				case 0:
					trr.Postal = field
				case 1:
					trr.Rate = field
				default:
					msg := "invalid number of rows in tax rates table"
					Fatal(msg)
					panic(msg)
				}
			}
			taxRateTable[trr.Postal] = trr
		}
	}
	return taxRateTable
}

func (taxTable *TaxTable) GetRate(postal string) string {
	first5 := postal[0:5]
	taxRateRow, ok := (*taxTable)[first5]
	if ok {
		return taxRateRow.Rate
	} else {
		Warn("Unknown sales tax for postal code " + postal)
		return averageSalesTax
	}
}
