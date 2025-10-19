package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func uploadFile(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Uploading file\n")
	req.ParseMultipartForm(50 << 20) // 50 mb

	file, handler, err := req.FormFile("videoFile")
	if err != nil {
		fmt.Println("Error getting video")
		fmt.Println(err)
		return
	}
	defer file.Close() // will close file when func returns... whoa
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	dst, err := os.Create("videos/" + handler.Filename)
	if err != nil {
		log.Println("error creating file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "uploaded file")
	http.Redirect(w, req, "/index", http.StatusSeeOther)
}

func videoSplit() {
	files, err := os.ReadDir("videos/")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range files {
		inPath := v.Name()
		outPath := strings.Split(v.Name(), ".")[0] + "/"
		cmd := exec.Command("ffmpeg", "-i", inPath, "-c", "copy", "-map", "0", "-segment_time", "00:00:30", "-f", "segment", "-reset_timestamps", "1", outPath+"/segment%03d.mp4")
		out, err := cmd.Output()
		fmt.Println(out)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Files split")
	fmt.Println(os.ReadDir("./videos"))
}

func main() {
	go videoSplit()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/upload", uploadFile)
	// http.HandleFunc("/process", videoSplit)
	http.ListenAndServe(":8090", nil)
	entries, err := os.ReadDir("videos/")
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	fmt.Printf("Files in /videos:\n")
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}

// pass this command to os.exec.Command(). That's it.
// ffmpeg -i video1.mp4 -c copy -map 0 -segment_time 00:00:10 -f segment -reset_timestamps 1 output%03d.mp4D
