package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/DmitriyKost/wbl0/pkg/structs"
	_ "github.com/lib/pq"
)

var db *sql.DB
var mutex sync.Mutex

// Inits database and cache in exactly this order. Whenever the pkg/database is imported (actually when application starts.)
func init() {
    initDB()
    initCache()
}


// Initializes database connection and executes necessary queries from ./config/migrations.sql file.
//
// If there's error while reading file, or executing migrations it panics.
func initDB() {
    connStr := "user=" + os.Getenv("DB_LOGIN") + 
    " dbname=" + os.Getenv("DB_NAME") + 
    " password=" + os.Getenv("DB_PASSWORD") + 
    " sslmode=disable"

    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Cannot connect to db: %v", err)
    }

    migrations, err := os.ReadFile("./config/migrations.sql")
    if err != nil {
        log.Fatalf("Failed to read migrations: %v", err)
    }
    for _, query := range strings.Split(string(migrations), "query_separator") {
        _, err := db.Exec(query)
        if err != nil {
            panic(err)
        }
    }
}

// Inserts order data into the database and caches it.
//
// If an order with the same UID already exists it does nothing.
//
// If there is an error during the database insertion or deserializing process, it returns the error.
func InsertOrder(order []byte) error {
    var orderUID structs.OrderUID
    if err := json.Unmarshal(order, &orderUID); err != nil {
        return err
    }
    
    query := "INSERT INTO orders (order_uid, order_data) VALUES ($1, $2) ON CONFLICT(order_uid) DO NOTHING"
    if res, err := db.Exec(query, orderUID.OrderUID, order); err != nil {
        return err
    } else if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
        return nil
    }
    cacheOrder(order, orderUID.OrderUID)
    return nil
}