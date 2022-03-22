package repository

import "database/sql"

type WarehouseItem struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Min      int32  `json:"min"`
	Max      int32  `json:"max"`
	Quantity int32  `json:"quantity"`
}

func (wi *WarehouseItem) CreateWarehouseItem(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO warehouse_items (name, quantity, min, max) VALUES ($1, $2, $3, $4) RETURNING id",
		&wi.Name, &wi.Quantity, &wi.Min, &wi.Max,
	).Scan(&wi.Id)
}

func (wi *WarehouseItem) GetWarehouseItemById(db *sql.DB) error {
	return db.QueryRow(
		"SELECT * FROM warehouse_items WHERE id = $1",
		wi.Id,
	).Scan(&wi.Id, &wi.Name, &wi.Quantity, &wi.Min, &wi.Max)
}

func GetWarehouseItems(db *sql.DB, start, count int) ([]WarehouseItem, error) {
	rows, err := db.Query("SELECT * FROM warehouse_items LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []WarehouseItem{}

	for rows.Next() {
		var wi WarehouseItem

		if err := rows.Scan(&wi.Id, &wi.Name, &wi.Quantity, &wi.Min, &wi.Max); err != nil {
			return nil, err
		}

		items = append(items, wi)
	}

	return items, nil
}
