package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

// ExportItemsToCSV exports items to a CSV file
func ExportItemsToCSV(items []Item, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ID", "Type", "Data"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write items to CSV
	for _, item := range items {
		row := []string{item.ID, item.Type, item.Data}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ImportItemsFromCSV imports items from a CSV file
func ImportItemsFromCSV(filePath string) ([]Item, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Parse CSV records into items
	var items []Item
	for _, record := range records {
		if len(record) >= 3 {
			item := Item{
				ID:   record[0],
				Type: record[1],
				Data: record[2],
			}
			items = append(items, item)
		}
	}

	return items, nil
}

// Example usage:
// filePath := "exported_data.csv"
// if err := ExportItemsToCSV(items, filePath); err != nil {
//     log.Println("Error exporting items to CSV:", err)
// }
// importedItems, err := ImportItemsFromCSV(filePath)
// if err != nil {
//     log.Println("Error importing items from CSV:", err)
// }
