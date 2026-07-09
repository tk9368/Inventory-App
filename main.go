package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// 1. UI Handler: Browser par Scanner Dashboard dikhane ke liye
func uiHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>IMS Dashboard</title>
		<style>
			body { font-family: 'Segoe UI', Arial, sans-serif; margin: 40px; background-color: #eef2f3; text-align: center; }
			.container { background: white; padding: 40px; border-radius: 12px; box-shadow: 0px 4px 15px rgba(0,0,0,0.05); display: inline-block; width: 350px; }
			h2 { color: #333; margin-bottom: 25px; }
			input, select, button { padding: 12px; margin: 10px 0; width: 100%; font-size: 15px; border-radius: 6px; border: 1px solid #ccc; box-sizing: border-box; }
			button { background-color: #28a745; color: white; border: none; font-weight: bold; cursor: pointer; transition: 0.2s; }
			button:hover { background-color: #218838; }
			select { background-color: #fff; }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>🏭 Workshop IMS Scanner</h2>
			<form action="/scan" method="POST">
				<input type="text" name="emp_id" placeholder="Scan Employee ID (e.g. EMP101)" required>
				<input type="text" name="part_id" placeholder="Scan Part Barcode (e.g. PART-BRAKE-01)" required>
				<select name="action">
					<option value="out">⚠️ Stock OUT (Use Part)</option>
					<option value="in">📥 Stock IN (Add to Inventory)</option>
				</select>
				<button type="submit">Submit Transaction</button>
			</form>
		</div>
	</body>
	</html>`
	fmt.Fprintf(w, html)
}

// 2. Backend Process Handler: Stock IN/OUT ka logical process
func inventoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	empID := r.FormValue("emp_id")
	partID := r.FormValue("part_id")
	action := r.FormValue("action")

	w.Header().Set("Content-Type", "text/html")

	// Employee verify karna
	var userExists bool
	db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE emp_id=$1)", empID).Scan(&userExists)
	if !userExists {
		fmt.Fprintf(w, "<h3 style='color:red; text-align:center;'>❌ Error: Invalid Employee ID (%s)!</h3><br><div style='text-align:center;'><a href='/'>Go Back</a></div>", empID)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if action == "out" {
		var currentStock int
		err := db.QueryRow("SELECT stock_count FROM parts WHERE part_id=$1", partID).Scan(&currentStock)
		if err == sql.ErrNoRows {
			tx.Rollback()
			fmt.Fprintf(w, "<h3 style='color:red; text-align:center;'>❌ Part (%s) Not Found!</h3><br><div style='text-align:center;'><a href='/'>Go Back</a></div>", partID)
			return
		} else if currentStock <= 0 {
			tx.Rollback()
			fmt.Fprintf(w, "<h3 style='color:orange; text-align:center;'>⚠️ Out of Stock! Cannot issue part.</h3><br><div style='text-align:center;'><a href='/'>Go Back</a></div>")
			return
		}

		tx.Exec("UPDATE parts SET stock_count = stock_count - 1 WHERE part_id=$1", partID)
		tx.Exec("INSERT INTO part_logs (emp_id, part_id, action_type, quantity) VALUES ($1, $2, 'OUT', 1)", empID, partID)
		fmt.Fprintf(w, "<h3 style='color:green; text-align:center;'>✅ Stock OUT Success! Remaining Stock: %d</h3><br><div style='text-align:center;'><a href='/'>Go Back</a></div>", currentStock-1)

	} else if action == "in" {
		var currentStock int
		err := db.QueryRow("SELECT stock_count FROM parts WHERE part_id=$1", partID).Scan(&currentStock)
		
		if err == sql.ErrNoRows {
			tx.Exec("INSERT INTO parts (part_id, part_name, stock_count) VALUES ($1, 'New Scanned Part', 1)", partID)
			currentStock = 0
		} else {
			tx.Exec("UPDATE parts SET stock_count = stock_count + 1 WHERE part_id=$1", partID)
		}

		tx.Exec("INSERT INTO part_logs (emp_id, part_id, action_type, quantity) VALUES ($1, $2, 'IN', 1)", empID, partID)
		fmt.Fprintf(w, "<h3 style='color:green; text-align:center;'>📥 Stock IN Success! New Stock: %d</h3><br><div style='text-align:center;'><a href='/'>Go Back</a></div>", currentStock+1)
	}

	tx.Commit()
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is missing!")
	}
	
	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", uiHandler)
	http.HandleFunc("/scan", inventoryHandler)

	fmt.Println("🏭 IMS Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
