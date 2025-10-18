package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "You have reached the index page. I haven't done anything here. Go to /hello or /headers :P\n")
	// http.
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

	// next step, make chunks and make those persistent. can discard original. i think. maybe. :O

	tempFile, err := os.CreateTemp("temp-files", "upload-*.mp4") // should change the format to accept more filetypes
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "Successfully uploaded File\n")

}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8090", nil)
}

// pass this command to os.exec.Command(). That's it.
// ffmpeg -i video1.mp4 -c copy -map 0 -segment_time 00:00:10 -f segment -reset_timestamps 1 output%03d.mp4D
func video_split(inPath string, outPath string) {
	cmd := exec.Command("ffmpeg", "-i", inPath, "-c", "copy", "-map", "0", "-segment_time", "00:00:30", "-f", "segment", "-reset_timestamps", "1", outPath+"/output%03d.mp4")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	// fmt.Print(output)
	// assert.Nil(t, err)
}
