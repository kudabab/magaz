package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kudabab/market-s/db"
	"github.com/kudabab/market-s/entity"
)

func GetProduct(category string) ([]entity.Product, error) {

	var query string

	switch category {
	case "tobacco":
		query = "SELECT * FROM tobacco"
	case "water":
		query = "SELECT * FROM water"
	case "beers":
		query = "SELECT * FROM beers"
	default:
		return nil, nil
	}

	rows, err := db.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error querying the database: %v", err)
	}
	defer rows.Close()

	var items []entity.Product

	// Обработка результатов запроса
	for rows.Next() {
		var item entity.Product
		if err := rows.Scan(&item.Id, &item.Name, &item.Price); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func GetProductByID(category string, productID int) (*entity.Product, error) {
	query := fmt.Sprintf("SELECT id, name, price FROM %s WHERE id=$1", category)

	row := db.DB.QueryRow(context.Background(), query, productID)
	fmt.Println("row: ", row)
	var product entity.Product
	err := row.Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	return &product, nil
}
