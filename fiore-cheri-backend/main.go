package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/go-sql-driver/mysql"
	"fiore-cheri-backend/models" // Import struct Menu dari menu.go
)

var db *sql.DB

func main() {
	// Koneksi ke database MySQL
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/fiore_cheri")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Tes koneksi
	err = db.Ping()
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}
	fmt.Println("Connected to MySQL!")

	// Setup Fiber
	app := fiber.New()

	// Endpoint simple test
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Fiore Chéri - Coffee Shop POS API is running!")
	})

	// Get All Menu
	app.Get("/menu", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, name, price FROM menu")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var menus []models.Menu
		for rows.Next() {
			var menu models.Menu
			if err := rows.Scan(&menu.ID, &menu.Name, &menu.Price); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			menus = append(menus, menu)
		}
		return c.JSON(menus)
	})

	// Add Menu
	app.Post("/menu", func(c *fiber.Ctx) error {
		var menu models.Menu
		if err := c.BodyParser(&menu); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}

		_, err := db.Exec("INSERT INTO menu (name, price) VALUES (?, ?)", menu.Name, menu.Price)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Menu added"})
	})

	// Update Menu
	app.Put("/menu/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var menu models.Menu
		if err := c.BodyParser(&menu); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}

		_, err := db.Exec("UPDATE menu SET name = ?, price = ? WHERE id = ?", menu.Name, menu.Price, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Menu updated"})
	})

	// Delete Menu
	app.Delete("/menu/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := db.Exec("DELETE FROM menu WHERE id = ?", id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Menu deleted"})
	})
	
	app.Post("/transactions", func(c *fiber.Ctx) error {
		var details []models.TransactionDetail
		if err := c.BodyParser(&details); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
	
		// Hitung total harga
		var totalPrice float64
		for _, item := range details {
			var price float64
			err := db.QueryRow("SELECT price FROM menu WHERE id = ?", item.MenuID).Scan(&price)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Menu not found"})
			}
			item.Subtotal = price * float64(item.Qty)
			totalPrice += item.Subtotal
		}
	
		// Simpan transaksi
		result, err := db.Exec("INSERT INTO transactions (total_price) VALUES (?)", totalPrice)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	
		// Dapatkan ID transaksi baru
		transactionID, _ := result.LastInsertId()
	
		// Simpan detail transaksi
		for _, item := range details {
			_, err := db.Exec("INSERT INTO transaction_details (transaction_id, menu_id, qty, subtotal) VALUES (?, ?, ?, ?)",
				transactionID, item.MenuID, item.Qty, item.Subtotal)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}
	
		return c.JSON(fiber.Map{"message": "Transaction successful", "transaction_id": transactionID, "total_price": totalPrice})
	})
	
	// Get All Transactions
	app.Get("/transactions", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, total_price, timestamp FROM transactions")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()
	
		var transactions []models.Transaction
		for rows.Next() {
			var transaction models.Transaction
			if err := rows.Scan(&transaction.ID, &transaction.TotalPrice, &transaction.Timestamp); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			transactions = append(transactions, transaction)
		}
		return c.JSON(transactions)
	})
	
	// Get Transaction Details
	app.Get("/transactions/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		rows, err := db.Query("SELECT menu.name, td.qty, td.subtotal FROM transaction_details td JOIN menu ON td.menu_id = menu.id WHERE td.transaction_id = ?", id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()
	
		var details []fiber.Map
		for rows.Next() {
			var name string
			var qty int
			var subtotal float64
			if err := rows.Scan(&name, &qty, &subtotal); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			details = append(details, fiber.Map{"menu": name, "qty": qty, "subtotal": subtotal})
		}
		return c.JSON(details)
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
