package database

import (
	"fmt"
	"log"
)

type Cache struct {
    orders map[string][]byte
}

var CachePtr *Cache


// Creates a new cache instance and fills it with orders from the database.
//
// It panics if there's error scanning orders.
func initCache() {
    CachePtr = &Cache{
        orders: make(map[string][]byte),
    }
    rows, err := db.Query("SELECT * FROM orders")
    if err != nil {
        log.Fatalf("Failed to get orders from db: %v", err)
        return
    }
    defer rows.Close()
    
    for rows.Next() {
        var orderUID string
        var order []byte
        if err := rows.Scan(&orderUID, &order); err != nil {
            log.Fatalf("Failed to scan order: %v", err)
        } else {
            CachePtr.orders[orderUID] = order
        }
    }
}

// Returns cached order data or error if the order is not in the cache.
func GetOrder(orderUID string) ([]byte, error) {
    mutex.Lock()
    defer mutex.Unlock()
    cachedOrder, ok := CachePtr.orders[orderUID]
    if !ok {
        return nil, fmt.Errorf("Order not found")
    }
    return cachedOrder, nil
}

// Stores the given order data in the cache under the provided orderUID.
// It acquires a lock to ensure thread safety.
func cacheOrder(order []byte, orderUID string) {
    mutex.Lock()
    defer mutex.Unlock()
    CachePtr.orders[orderUID] = order
}
