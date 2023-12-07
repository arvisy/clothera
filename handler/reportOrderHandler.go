package handler

import (
	"database/sql"
	"fmt"
	"pair-project/entity"
	"strings"
)

func OrderReportMenu(db *sql.DB) error {
	fmt.Println("1 -> Total Revenue & Quantity Sold")
	fmt.Println("2 -> Rental Revenue by Costume")
	fmt.Println("3 -> Sales Revenue by Clothes")
	fmt.Println("4 -> Back to Main Menu")
	fmt.Println("")

	var choice int
	fmt.Print("Choice: ")
	_, err := fmt.Scan(&choice)
	if err != nil {
		return fmt.Errorf("OrderReportMenu: %w", err)
	}
	fmt.Println("")

	switch choice {
	case 1:
		err = AllRevenue(db)
		if err != nil {
			return fmt.Errorf("OrderReportMenu: %w", err)
		}
		err = TotalQuantity(db)
		if err != nil {
			return fmt.Errorf("OrderReportMenu: %w", err)
		}
	case 2:
		err = RentalRevenueByCostume(db)
		if err != nil {
			return fmt.Errorf("OrderReportMenu: %w", err)
		}
	case 3:
		err = SalesRevenueByClothes(db)
		if err != nil {
			return fmt.Errorf("OrderReportMenu: %w", err)
		}
	case 4:
		return nil
	}
	return nil

}

func TotalQuantity(db *sql.DB) error {
	query := `SELECT
	SUM(rents.Quantity) AS rentQuantity,
    SUM(sales.Quantity) AS saleQuantity
FROM orders
JOIN rents ON orders.OrderID = rents.OrderID
JOIN sales ON orders.OrderID = sales.OrderID`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("TotalOrders: %w", err)
	}
	defer rows.Close()

	var rentQuantity, saleQuantity int

	for rows.Next() {
		err := rows.Scan(&rentQuantity, &saleQuantity)
		if err != nil {
			return fmt.Errorf("RentPrice: %w", err)
		}
	}

	fmt.Println("Showing Total Quantity Sold...")
	fmt.Printf("\n%-20s| %-15s| %-15s\n", "Total Products Sold", "Total Rentals", "Total Sales")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("%-20d| %-15d| %-15d\n", (rentQuantity + saleQuantity), rentQuantity, saleQuantity)

	fmt.Println("")
	fmt.Println("")

	return nil
}

func AllRevenue(db *sql.DB) error {
	// get values
	totalRev, err := TotalRevenue(db)
	if err != nil {
		return fmt.Errorf("AllRevenue: %w", err)
	}

	rentRev, err := TotalRentRevenue(db)
	if err != nil {
		return fmt.Errorf("AllRevenue: %w", err)
	}

	salesRev, err := TotalSalesRevenue(db)
	if err != nil {
		return fmt.Errorf("AllRevenue: %w", err)
	}

	// print
	fmt.Println("Showing Total Revenue...")
	fmt.Printf("\n%-15s| %-15s| %-15s\n", "Total Revenue", "Rental Revenue", "Sales Revenue")
	fmt.Println("-----------------------------------------------")
	fmt.Printf("%-15.2f| %-15.2f| %-15.2f\n", totalRev, rentRev, salesRev)

	fmt.Println("")

	return nil
}

func TotalRevenue(db *sql.DB) (float64, error) {
	query := `SELECT SUM(orders.totalPrice) AS totalRevenue FROM orders`

	rows, err := db.Query(query)
	if err != nil {
		return 0, fmt.Errorf("TotalRevenue: %w", err)
	}
	defer rows.Close()

	var totalRev float64

	for rows.Next() {
		err := rows.Scan(&totalRev)
		if err != nil {
			return 0, fmt.Errorf("TotalRevenue: %w", err)
		}
	}
	return totalRev, nil
}

func TotalRentRevenue(db *sql.DB) (float64, error) {
	query := `SELECT SUM(rents.rentPrice) AS totalRentRevenue FROM rents`

	rows, err := db.Query(query)
	if err != nil {
		return 0, fmt.Errorf("TotalRentRevenue: %w", err)
	}
	defer rows.Close()

	var totalRev float64

	for rows.Next() {
		err := rows.Scan(&totalRev)
		if err != nil {
			return 0, fmt.Errorf("TotalRentRevenue: %w", err)
		}
	}
	return totalRev, nil
}

func TotalSalesRevenue(db *sql.DB) (float64, error) {
	query := `SELECT SUM(sales.Quantity*clothes.ClothesPrice) AS totalSalesRevenue
				FROM sales JOIN clothes ON sales.ClothesID = clothes.ClothesID`

	rows, err := db.Query(query)
	if err != nil {
		return 0, fmt.Errorf("TotalSalesRevenue: %w", err)
	}
	defer rows.Close()

	var totalRev float64

	for rows.Next() {
		err := rows.Scan(&totalRev)
		if err != nil {
			return 0, fmt.Errorf("TotalSalesRevenue: %w", err)
		}
	}
	return totalRev, nil
}

func RentalRevenueByCostume(db *sql.DB) error {
	query := `SELECT 
	costumes.CostumeName, 
	SUM(rents.Quantity) AS Quantity, 
	SUM(rents.RentPrice) AS TotalRentPrice
FROM rents
JOIN costumes ON rents.CostumeID = costumes.CostumeID
JOIN orders ON rents.OrderID = orders.OrderID
GROUP BY costumes.CostumeName`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("RentalRevenueByCostume: %w", err)
	}
	defer rows.Close()

	var rentals []entity.RevenueByCostume

	for rows.Next() {
		var r entity.RevenueByCostume
		err := rows.Scan(&r.CostumeName, &r.Quantity, &r.TotalRevenue)
		if err != nil {
			return fmt.Errorf("RentalRevenueByCostume: %w", err)
		}
		rentals = append(rentals, r)
	}

	fmt.Println("Showing Rental Revenue by Costume...")
	fmt.Printf("\n%-12s| %-18s| %-15s\n", "Costume", "Total Quantity", "Total Revenue")
	fmt.Println(strings.Repeat("-", 50))

	for _, r := range rentals {
		fmt.Printf("%-12s| %-18d| %-15.2f\n", r.CostumeName, r.Quantity, r.TotalRevenue)
	}

	fmt.Println("")
	fmt.Println("")

	return nil
}

func SalesRevenueByClothes(db *sql.DB) error {
	query := `SELECT 
    clothes.ClothesName, 
	(sales.Quantity) AS Quantity, 
	SUM(sales.Quantity*clothes.ClothesPrice) AS TotalSalesPrice
FROM sales
JOIN clothes ON sales.ClothesID = clothes.ClothesID
GROUP BY sales.ClothesID
ORDER BY TotalSalesPrice DESC`

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("SalesRevenueByClothes: %w", err)
	}
	defer rows.Close()

	var sales []entity.RevenueByClothes

	for rows.Next() {
		var s entity.RevenueByClothes
		err := rows.Scan(&s.Name, &s.Quantity, &s.TotalRevenue)
		if err != nil {
			return fmt.Errorf("SalesRevenueByClothes: %w", err)
		}
		sales = append(sales, s)
	}

	fmt.Println("Showing Sales Revenue by Clothes...")
	fmt.Printf("\n%-17s| %-12s| %-15s\n", "Clothes", "Quantity", "Total Revenue")
	fmt.Println(strings.Repeat("-", 50))

	for _, s := range sales {
		fmt.Printf("%-17s| %-12d| %-15.2f\n", s.Name, s.Quantity, s.TotalRevenue)
	}
	fmt.Println("")
	fmt.Println("")

	return nil
}
