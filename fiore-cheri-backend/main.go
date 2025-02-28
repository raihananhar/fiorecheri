package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Koneksi ke database MySQL
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/fiore_cheri")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

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

		var menus []Menu
		for rows.Next() {
			var menu Menu
			if err := rows.Scan(&menu.ID, &menu.Name, &menu.Price); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			menus = append(menus, menu)
		}
		return c.JSON(menus)
	})

	// Add Menu
	app.Post("/menu", func(c *fiber.Ctx) error {
		var menu Menu
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
		var menu Menu
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


	// Start server
	log.Fatal(app.Listen(":3000"))
}