package internal

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/xsadia/secred/repository"
)

const badCSVHeadersError = "bad headers on CSV file"

// TODO upgrade to Go 1.18 to use generics for school items as well
func newWarehouseItemListFromCSV(lines [][]string) ([]repository.WarehouseItem, error) {
	wil := make([]repository.WarehouseItem, len(lines)-1)
	indexes := make(map[string]int, len(lines[0]))

	for i, v := range lines[0] {

		lowerV := strings.ToLower(v)

		if _, ok := indexes[lowerV]; ok {
			return nil, errors.New(badCSVHeadersError)
		}

		indexes[lowerV] = i
	}

	_, nameOk := indexes["name"]
	_, quantityOk := indexes["quantity"]
	_, minOk := indexes["min"]
	_, maxOk := indexes["max"]

	if !nameOk || !quantityOk || !minOk || !maxOk {
		return nil, errors.New(badCSVHeadersError)
	}

	// Sadly O(n2) because I have to use Atoi and ToLower
	for i := 1; i < len(lines); i++ {

		var wi repository.WarehouseItem

		quantity, _ := strconv.Atoi(lines[i][indexes["quantity"]])
		min, _ := strconv.Atoi(lines[i][indexes["min"]])
		max, _ := strconv.Atoi(lines[i][indexes["max"]])

		wi.Name = strings.ToLower(lines[i][indexes["name"]])
		wi.Quantity = int32(quantity)
		wi.Min = int32(min)
		wi.Max = int32(max)

		wil[i-1] = wi
	}

	return wil, nil
}

func ParseCSV(fileName string) ([]repository.WarehouseItem, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(f)
	rec, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	wil, err := newWarehouseItemListFromCSV(rec)

	if err != nil {
		return nil, err
	}

	return wil, nil
}
