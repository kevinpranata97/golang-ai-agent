package models

import (
	"time"
	"database/sql"
)

// Product represents the Product entity
type Product struct {
	Id int `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Description string `json:"description"`
	Price float64 `json:"price" validate:"required"`
	Created_at time.Time `json:"created_at" validate:"required"`
}

// CreateProduct creates a new Product in the database
func CreateProduct(db *sql.DB, product *Product) error {
	query := `INSERT INTO products (name, description, price) VALUES (?, ?, ?)`
	
	result, err := db.Exec(query, product.Name, product.Description, product.Price)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.Id = int(id)
	return nil
}

// GetProductByID retrieves a Product by ID
func GetProductByID(db *sql.DB, id int) (*Product, error) {
	product := &Product{}
	query := `SELECT id, name, description, price, created_at FROM products WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Created_at)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetAllProducts retrieves all Products
func GetAllProducts(db *sql.DB) ([]Product, error) {
	query := `SELECT id, name, description, price, created_at FROM products`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		product := Product{}
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Created_at)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct updates a Product in the database
func UpdateProduct(db *sql.DB, product *Product) error {
	query := `UPDATE products SET name = ?, description = ?, price = ? WHERE id = ?`
	
	_, err := db.Exec(query, product.Name, product.Description, product.Price, product.Id)
	return err
}

// DeleteProduct deletes a Product from the database
func DeleteProduct(db *sql.DB, id int) error {
	query := `DELETE FROM products WHERE id = ?`
	
	_, err := db.Exec(query, id)
	return err
}
