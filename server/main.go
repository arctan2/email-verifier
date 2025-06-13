package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"email_verify/dbconn"
	"email_verify/webroutes"
	"email_verify/respond"
)

func HeaderMiddleware(headers map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range headers {
				w.Header().Set(key, value)
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	res := respond.ResponseStruct{ Err: false, Msg: "pong" }
	json.NewEncoder(w).Encode(&res)
}

var ADDR string = ":8000"

func printIpv4() {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Printf("\thttp://%s%s\n", ipv4, ADDR)
		}
	}
	fmt.Printf("\thttp://localhost%s\n", ADDR)
}

func main() {
	isProdDBFlag := flag.Bool("prod-db", false, "run production db")
	isProdModeFlag := flag.Bool("prod-mode", false, "run in production mode with static file serve")
	portFlag := flag.String("port", "", "specify port. default: 8000")

	flag.Parse()

	if *portFlag != "" {
		ADDR = ":" + *portFlag
	}

	var db *sql.DB

	if *isProdDBFlag {
		db = dbconn.DbConn(*isProdModeFlag)
	} else {
		db = dbconn.DbTestConn()
	}

	if db == nil {
		return
	}

	webMux := webroutes.NewWebRoutesMux(db)
	mainMux := http.NewServeMux()

	headers := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization",
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}

	responseHeaders := HeaderMiddleware(headers)

	pingHandler := http.HandlerFunc(ping)

	if *isProdModeFlag {
		staticDir := "./dist"
		fs := http.FileServer(http.Dir(staticDir))

		mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join(staticDir, r.URL.Path)
			info, err := os.Stat(path)

			if err == nil && !info.IsDir() {
				fs.ServeHTTP(w, r)
				return
			}

			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		})
	}

	mainMux.Handle("/api/ping", responseHeaders(pingHandler))
	mainMux.Handle("/api/web/", responseHeaders(http.StripPrefix("/api/web", webMux)))

	server := http.Server{
		Addr: ADDR,
		Handler: mainMux,
	}

	fmt.Println("listening on: ")
	printIpv4()
	log.Fatal(server.ListenAndServe())
}

