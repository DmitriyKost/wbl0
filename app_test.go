package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DmitriyKost/wbl0/pkg"
	"github.com/DmitriyKost/wbl0/pkg/database"
)

func TestIndex(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(pkg.Index)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}

func TestGetOrder(t *testing.T) {
    req, err := http.NewRequest("GET", "/get/hello/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(pkg.GetOrder)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func TestInsertOrder(t *testing.T) {
    err := database.InsertOrder([]byte(stringOrderData))
    if err != nil {
        t.Errorf("Error inserting test order to cache and db: %v", err)
    }
}


func TestGetOrderFromCache(t *testing.T) {
    orderData := database.OrderCache.Get("b563feb7b2b84b6test")
    if orderData == nil {
        t.Errorf("Error getting test order from cache")
    }
    if string(orderData) != stringOrderData {
        t.Fatalf("Got wrong test order data from cache: %s", string(orderData))
    }
}

var stringOrderData string = `{"entry": "WBIL", "items": [{"rid": "ab4219087a764ae0btest", "name": "Mascaras", "sale": 30, "size": "0", "brand": "Vivienne Sabo", "nm_id": 2389212, "price": 453, "status": 202, "chrt_id": 9934930, "total_price": 317, "track_number": "WBILMTESTTRACK"}], "sm_id": 99, "locale": "en", "payment": {"bank": "alpha", "amount": 1817, "currency": "USD", "provider": "wbpay", "custom_fee": 0, "payment_dt": 1637907727, "request_id": "", "goods_total": 317, "transaction": "b563feb7b2b84b6test", "delivery_cost": 1500}, "delivery": {"zip": "2639809", "city": "Kiryat Mozkin", "name": "Test Testov", "email": "test@gmail.com", "phone": "+9720000000", "region": "Kraiot", "address": "Ploshad Mira 15"}, "shardkey": "9", "oof_shard": "1", "order_uid": "b563feb7b2b84b6test", "customer_id": "test", "date_created": "2021-11-26T06:22:19Z", "track_number": "WBILMTESTTRACK", "delivery_service": "meest", "internal_signature": ""}`
