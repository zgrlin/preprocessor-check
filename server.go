package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Dosya Boyutu
const MAX_UPLOAD_SIZE = 10240 * 1024 // 100mb set

// Cookie icin store
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

// Giris icin kullanicilar
var users = map[string]string{
	"user": "pass",
	"user2": "pass2",
}

type process struct {
	Name   string
	Status string
}

type Progress struct {
	TotalSize int64
	BytesRead int64
}

func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("Yeni bir dosya yuklendi!")
		return
	}
	fmt.Println("Su anda yukleniyor: %d\n", pr.BytesRead)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {

		grep := exec.Command("grep", "java")
		ps := exec.Command("ps", "-ef")
		pipe, _ := ps.StdoutPipe()
		defer pipe.Close()
		grep.Stdin = pipe
		ps.Start()
		res, _ := grep.Output()
		fmt.Printf("Process kontrol edildi \n")
		w.Write(res)
	} else {
		http.Error(w, "bye bye.", http.StatusForbidden)
	}
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {

		pkill := exec.Command("pkill", "-f", "process")
		out, err := pkill.CombinedOutput()
		fmt.Printf("Java durduruldu \n")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		if err != nil {
			fmt.Printf("Problem: %s\n", string(out))
		}
		w.Write(out)
	} else {
		http.Error(w, "bye bye.", http.StatusForbidden)
	}
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {

		start := exec.Command("sh", "start.sh")
		out, err := start.CombinedOutput()
		fmt.Printf("Process baslatildi \n")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		if err != nil {
			fmt.Printf("Problem: \n%s\n", (err))
		}
		w.Write(out)
	} else {
		http.Error(w, "bye bye.", http.StatusForbidden)
	}
}

// Giris icin kontrol
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Duzgun formatta gonder!",
			http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	if originalPassword, ok := users[username]; ok {
		session, _ := store.Get(r, "session.id")
		if password == originalPassword {
			session.Values["authenticated"] = true
			session.Save(r, w)
		} else {
			fmt.Printf("hatali sifre denemesi.\n")
			http.Error(w, "sifre?", http.StatusUnauthorized)
			return
		}
	} else {
		fmt.Printf("hatali giris denemesi.\n")
		http.Error(w, "kullanici adi?", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

// Cikis icin kontrol
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		session.Values["authenticated"] = false
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		fmt.Printf("kullanici cikis yapti.\n")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// Dosya Yukleme Mekanizmasi
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {

		if r.Method != "POST" {
			http.Error(w, "Olmaz isler pesindesin.", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["file"]

		for _, fileHeader := range files {
			if fileHeader.Size > MAX_UPLOAD_SIZE {
				http.Error(w, fmt.Sprintf("Dosya boyutu cok buyuk: %s. Biraz kucultun.", fileHeader.Filename), http.StatusBadRequest)
				return
			}

			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer file.Close()

			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			filetype := http.DetectContentType(buff)
			if filetype != "application/zip" && filetype != "script" {
				http.Error(w, "Dosya formati desteklenmiyor.", http.StatusBadRequest)
				return
			}

			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = os.MkdirAll("./uploads", os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			defer f.Close()

			pr := &Progress{
				TotalSize: fileHeader.Size,
			}

			_, err = io.Copy(f, io.TeeReader(file, pr))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		t, _ := template.ParseFiles("finish.gtpl")
		t.Execute(w, nil)
		return
	} else {
		fmt.Printf("Yetkisiz kullanici dosya yuklemeye calisti.\n")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// Dashboard
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {

		if r.Method != "POST" {
			t, _ := template.ParseFiles("upload.gtpl")
			t.Execute(w, nil)
			return
		}
	} else {
		fmt.Printf("Yetkisiz kullanici ana sayfaya ulasmak istiyor.\n")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// Ana Sayfa
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		if r.Method != "POST" {
			t, _ := template.ParseFiles("login.gtpl")
			t.Execute(w, nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Format Hatasi: %v", err)
			return
		}
	}
}

// PING PONG
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	process := process{Name: "process_name", Status: "process_status"}
	jsonResponse, jsonError := json.Marshal(process)

	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}

	fmt.Println(string(jsonResponse))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {

	log.Println("Sunucu baslatildi.")

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/dashboard", dashboardHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	r.HandleFunc("/process", processHandler)
	r.HandleFunc("/stop", stopHandler)
	r.HandleFunc("/start", startHandler)
	r.HandleFunc("/status", statusHandler).Methods("GET")

	fs := http.FileServer(http.Dir("log"))
	r.PathPrefix("/log/").Handler(http.StripPrefix("/log/", fs))

	fmt.Printf("/ modulu baslatildi.\n")
	fmt.Printf("/login modulu baslatildi.\n")
	fmt.Printf("/dashboard modulu baslatildi.\n")
	fmt.Printf("/upload modulu baslatildi.\n")
	fmt.Printf("/logout modulu baslatildi.\n")
	fmt.Printf("/process modulu baslatildi.\n")
	fmt.Printf("/stop modulu baslatildi.\n")
	fmt.Printf("/start modulu baslatildi.\n")
	fmt.Printf("/log/ modulu baslatildi.\n")

	tls_cfg := &tls.Config {
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	httpServer := &http.Server{
		Handler:      r,
		Addr:         "IP:PORT",
		TLSConfig:    tls_cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServeTLS("server.key", "server.pem"))
}
