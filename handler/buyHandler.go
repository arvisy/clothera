package handler

import (
	"database/sql"
	"fmt"
	"pair-project/entity"
)

func AddClothes(db *sql.DB, clothes entity.Clothes, customer entity.Customer, orderID int) (int, error) {
	for {
		fmt.Print("Enter quantity: ")
		var quantity int
		_, err := fmt.Scan(&quantity)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid quantity.")
			continue
		}

		available, err := CheckClothesAvailability(db, clothes.ClothesID, quantity)
		if err != nil {
			return 0, fmt.Errorf("error checking clothes availability: %v", err)
		}

		if !available {
			fmt.Println("Insufficient stock for the selected clothes. Please enter a valid quantity.")
			continue
		}

		saleQuery := `
			INSERT INTO Sales (OrderID, ClothesID, Quantity)
			VALUES (?, ?, ?)
		`
		_, err = db.Exec(saleQuery, orderID, clothes.ClothesID, quantity)
		if err != nil {
			return 0, fmt.Errorf("error inserting sale record: %v", err)
		}

		err = UpdateStock(db, clothes.ClothesID, quantity)
		if err != nil {
			return 0, fmt.Errorf("error updating stock: %v", err)
		}

		fmt.Println("Clothes added to the order successfully.")
		return quantity, nil
	}
}

func CheckClothesAvailability(db *sql.DB, clothesID int, quantity int) (bool, error) {
	stockQuery := `
		SELECT ClothesStock
		FROM Clothes
		WHERE ClothesID = ?
	`
	var stock int
	err := db.QueryRow(stockQuery, clothesID).Scan(&stock)
	if err != nil {
		return false, err
	}

	return stock >= quantity, nil
}

func UpdateStock(db *sql.DB, clothesID int, quantity int) error {
	updateStockQuery := `
		UPDATE Clothes
		SET ClothesStock = ClothesStock - ?
		WHERE ClothesID = ?
	`

	_, err := db.Exec(updateStockQuery, quantity, clothesID)
	if err != nil {
		return err
	}

	return nil
}

func GetClothesByID(db *sql.DB, clothesID int) (*entity.Clothes, error) {
	query := `
		SELECT ClothesID, ClothesName, ClothesCategory, ClothesPrice, ClothesStock
		FROM Clothes
		WHERE ClothesID = ?
	`

	var clothes entity.Clothes
	err := db.QueryRow(query, clothesID).Scan(
		&clothes.ClothesID,
		&clothes.ClothesName,
		&clothes.ClothesCategory,
		&clothes.ClothesPrice,
		&clothes.ClothesStock,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get clothes by ID: %v", err)
	}

	return &clothes, nil
}

func ByName(db *sql.DB, clothesName string) (int, error) {
	query := `SELECT ClothesID FROM Clothes WHERE ClothesName = ?`
	var clothesID int
	err := db.QueryRow(query, clothesName).Scan(&clothesID)

	if err != nil {
		return -1, err
	}

	return clothesID, nil
}

func GetPriceClothes(db *sql.DB, clothesID int) (float64, error) {
	query := `SELECT ClothesPrice FROM Clothes WHERE ClothesID = ?`

	var clothesPrice float64
	err := db.QueryRow(query, clothesID).Scan(&clothesPrice)

	if err != nil {
		return -1, err
	}

	return clothesPrice, nil
}
