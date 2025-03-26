package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// PageData holds dynamic content for templates
type PageData struct {
	Title  string
	Movie  string
	User   string
	Logged bool
}

// renderTemplate loads and executes the template with data
func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	tmplPath := filepath.Join("templates", tmpl+".html")
	basePath := filepath.Join("templates", "base.html")

	log.Println("🔹 Rendering:", tmplPath)

	t, err := template.ParseFiles(basePath, tmplPath)
	if err != nil {
		log.Println("❌ Template error:", err)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("❌ Execution error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// Check if user is logged in
func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	return err == nil && cookie.Value == "logged_in"
}

// Remove .mp4 extension from movie names
func trimExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// Serve video files with correct MIME type
func serveVideo(w http.ResponseWriter, r *http.Request) {
	videoFile := "static/videos" + r.URL.Path[len("/videos/"):]
	log.Println("🎬 Video requested:", videoFile)

	// Ensure the request is for an MP4 file
	if !strings.HasSuffix(videoFile, ".mp4") {
		http.Error(w, "Invalid video format", http.StatusForbidden)
		return
	}

	// Set correct Content-Type
	w.Header().Set("Content-Type", "video/mp4")

	http.ServeFile(w, r, videoFile)
}

func main() {
	// Login Page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if isAuthenticated(r) {
			http.Redirect(w, r, "/catalog", http.StatusFound)
			return
		}
		renderTemplate(w, "login", PageData{Title: "Login"})
	})

	// Handle Login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username == "admin" && password == "1234" { // Simple check (Replace with DB validation)
				http.SetCookie(w, &http.Cookie{
					Name:  "session",
					Value: "logged_in",
					Path:  "/",
				})
				http.Redirect(w, r, "/catalog", http.StatusFound)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// Movie Catalog (Only accessible after login)
	http.HandleFunc("/catalog", func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		renderTemplate(w, "catalog", PageData{Title: "Movie Catalog"})
	})

	// Movie Player (Only accessible after login)
	http.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		movie := r.URL.Query().Get("movie")
		if movie == "" {
			movie = "default.mp4" // Default video if none is selected
		}
		movieTitle := trimExtension(movie) // Remove .mp4
		renderTemplate(w, "player", PageData{Title: "Now Playing", Movie: movieTitle})
	})

	// Logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // Expire the cookie
		})
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// Serve static files (CSS, JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", http.FileServer(http.Dir("thumbnails"))))
	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir("static/videos"))))

	// Start server
	log.Println("🚀 Server started at http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("❌ Server error:", err)
	}
}
