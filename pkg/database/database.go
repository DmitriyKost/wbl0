package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DmitriyKost/wbl0/pkg/structs"
	_ "github.com/lib/pq"
)

var db *sql.DB
var mutex sync.Mutex

// Inits database and cache in exactly this order, whenever the pkg/database is imported (actually when application starts).
func init() {
    initDB()
}


// Initializes database connection and executes necessary queries from ./config/migrations.sql file.
//
// If there's error while reading file, or executing migrations it panics.
func initDB() {
    var connStr string
    var err error

    // Check if environment variable for remote database connection is set
    if os.Getenv("REMOTE_DB_HOST") != "" {
        connStr = "host=" + os.Getenv("REMOTE_DB_HOST") +
            " port=" + os.Getenv("REMOTE_DB_PORT") +
            " user=" + os.Getenv("REMOTE_DB_LOGIN") +
            " password=" + os.Getenv("REMOTE_DB_PASSWORD") +
            " dbname=" + os.Getenv("REMOTE_DB_NAME") +
            " sslmode=disable"
    } else { // Use local database configuration
        connStr = "user=" + os.Getenv("DB_LOGIN") +
            " dbname=" + os.Getenv("DB_NAME") +
            " password=" + os.Getenv("DB_PASSWORD") +
            " sslmode=disable"
    }    
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Cannot connect to db: %v", err)
    }

    // Reads the table schemas and executes CREATE TABLE queries, if table not exists.
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
    OrderCache.Set(orderUID.OrderUID, order, time.Hour)
    return nil
}

func getOrder(orderUID string) (error, []byte) {
    var order []byte
    query := "SELECT order_data FROM orders WHERE order_uid=$1";
    row := db.QueryRow(query, orderUID)
    if err := row.Scan(&order); err != nil {
        log.Printf("Error getting order from db: %v", err)
        return err, nil
    }
    return nil, order
}
