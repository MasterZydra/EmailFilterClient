package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// startWebServer starts an HTTP server to serve the log file
func startWebServer(port string, basicAuthPassword string) {
	http.HandleFunc("/", authMiddleware(basicAuthPassword, indexHandler))

	// Route to serve the log file
	http.HandleFunc("/log", authMiddleware(basicAuthPassword, logHandler))

	// Route to clear the log file
	http.HandleFunc("/log/clear", authMiddleware(basicAuthPassword, logClearHandler))

	// Route to view config file
	http.HandleFunc("/config", authMiddleware(basicAuthPassword, configHandler))

	// Route to update config file
	http.HandleFunc("/config/update", authMiddleware(basicAuthPassword, configUpdateHandler))

	// Route to view config file
	http.HandleFunc("/blacklist", authMiddleware(basicAuthPassword, blacklistHandler))

	// Route to update blacklist file
	http.HandleFunc("/blacklist/update", authMiddleware(basicAuthPassword, blacklistUpdateHandler))

	fmt.Printf("Starting web server on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}

// authMiddleware is a middleware that checks for the secret key
func authMiddleware(basicAuthPassword string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if basicAuthPassword == "" {
			next(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode the Basic Auth credentials
		username, password, ok := r.BasicAuth()
		if !ok || username != "mailadmin" || password != basicAuthPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the next handler if the credentials are valid
		next(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Write the response
	w.Write(getHtml(`
		<p>
			<a href="/log">Show log file</a><br>
			<a href="/log/clear">Clear log file</a><br>
		</p>

		<p>
			<a href="/config">Show config file</a><br>
			<a href="/blacklist">Show blacklist file</a><br>
		</p>
	`))
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read all lines from the log file
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		http.Error(w, "Could not read log file", http.StatusInternalServerError)
		return
	}

	// Write to the response
	w.Write(getHtml(`
		<h2>info.log</h2>
		<form action="/log/clear" method="get" style="height:100%">
			<input type="submit" value="Clear" />
			<textarea style="width:100%; min-height:95%;">` + string(content) + `</textarea>
		</form>
	`))
}

func logClearHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Truncate the log file
	err := os.Truncate(logFilePath, 0)
	if err != nil {
		http.Error(w, "Could not clear log file", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Write(getHtml(`
		<h2>info.log</h2>
		<p>Log file cleared successfully</p>
		<a href="/log">Back to log</a>
	`))
	log.Println("Log file cleared via /clear-log route")
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read all lines from the config file
	configData, err := os.ReadFile(configJsonPath)
	if err != nil {
		http.Error(w, "Could not read config file", http.StatusInternalServerError)
		return
	}

	// Write the config data to the response
	w.Write(getHtml(`
		<h2>config.json</h2>
		<form action="/config/update" method="post" style="height:100%">
			<input type="submit" value="Update" />
			<textarea name="config" style="width:100%; min-height:95%;">` + string(configData) + `</textarea>
		</form>
	`))
}

func configUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the updated config from the form
	updatedConfig := r.FormValue("config")
	if updatedConfig == "" {
		http.Error(w, "Config data is empty", http.StatusBadRequest)
		return
	}

	// Write the updated config back to the file
	err = os.WriteFile(configJsonPath, []byte(updatedConfig), 0644)
	if err != nil {
		http.Error(w, "Failed to write config file", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write(getHtml(`
		<h2>config.json</h2>
		<p>Config file updated successfully</p>
		<a href="/config">Back to config</a>
	`))
	log.Println("Config file updated via /config/update route")
}

func blacklistHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read all lines from the blacklist file
	blacklistData, err := os.ReadFile(blacklistJsonPath)
	if err != nil {
		http.Error(w, "Could not read blacklist file", http.StatusInternalServerError)
		return
	}

	// Write the blacklist data to the response
	w.Write(getHtml(`
		<h2>blacklist.json</h2>
		<form action="/blacklist/update" method="post" style="height:100%">
			<input type="submit" value="Update" />
			<textarea name="blacklist" style="width:100%; min-height:95%;">` + string(blacklistData) + `</textarea>
		</form>
	`))
}

func blacklistUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the updated blacklist from the form
	updatedBlacklist := r.FormValue("blacklist")
	if updatedBlacklist == "" {
		http.Error(w, "Blacklist data is empty", http.StatusBadRequest)
		return
	}

	// Write the updated blacklist back to the file
	err = os.WriteFile(blacklistJsonPath, []byte(updatedBlacklist), 0644)
	if err != nil {
		http.Error(w, "Failed to write blacklist file", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Write(getHtml(`
		<h2>blacklist.json</h2>
		<p>Blacklist file updated successfully</p>
		<a href="/blacklist">Back to blacklist</a>
	`))
	log.Println("Blacklist file updated via /blacklist/update route")
}

func getHtml(body string) []byte {
	return []byte(`<html>
		<head>
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<style>
				h1 { font-size: 1.5em; }
				h2 { font-size: 1.25em; }
			</style>
		</head>
		<body>
			<h1><a href="/" style="text-decoration: none;color: inherit;">Email Filter Client</a></h1>
			` + body + `
			<br>
		</body>
		<footer>
			<small>Version ` + version + `</small>
		</footer>
	</html>`)
}
