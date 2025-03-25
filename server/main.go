package main

import (
	"fmt"
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

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	fmt.Println(os.Getwd())
	video_split("video1.mp4", "videos")
	// http.ListenAndServe(":8090", nil)
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
