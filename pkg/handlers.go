package pkg

import (
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/DmitriyKost/wbl0/config"
	"github.com/DmitriyKost/wbl0/pkg/database"
)

// Seeking for templates
var templates = template.Must(template.ParseGlob("static/templates/*.html"))

// Returns index HTML template
func Index(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index.html")
}

// If the order is not found or if there's an error retrieving the order, it sends a bad request (HTTP 400) response.
//
// If the order is found, it writes the order data as JSON to the response body.
func GetOrder(w http.ResponseWriter, r *http.Request) {
    orderUID := r.PathValue("order_uid")
    if order, err := database.GetOrder(orderUID); err != nil {
        http.Error(w, "There's no order with this id :(", http.StatusBadRequest)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(order)
    }
}

// Renders the specified HTML template.
//
// If there's error during template rendering, it sends an HTTP 500 status code.
func renderTemplate(w http.ResponseWriter, tmpl string) {
    err := templates.ExecuteTemplate(w, tmpl, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// Middleware function that logs information about incoming HTTP requests.
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        recorder := &responseRecorder{w, http.StatusOK}
        next.ServeHTTP(recorder, r)
        ip := r.RemoteAddr
        log.Printf("%s - - [%s] \"%s %s %s\" %d",
            ip,
            start.Format("02/Jan/2006:15:04:05 -0700"),
            r.Method,
            r.URL.RequestURI(),
            r.Proto,
            recorder.status,
        )
    })
}


// responseRecorder is a custom implementation of http.ResponseWriter that captures the HTTP response status code.
type responseRecorder struct {
    http.ResponseWriter
    status int
}

// WriteHeader records the HTTP response status code before calling the underlying ResponseWriter's WriteHeader method.
func (r *responseRecorder) WriteHeader(statusCode int) {
    r.status = statusCode
    r.ResponseWriter.WriteHeader(statusCode)
}
