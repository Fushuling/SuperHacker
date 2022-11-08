package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/kbinani/screenshot"
)

func takeShot() (*image.RGBA, error) {
	n := screenshot.NumActiveDisplays()
	if n < 0 {
		return nil, errors.New("no screen detected")
	}
	bounds := image.Rectangle{image.Point{0, 0}, image.Point{1923, 1080}} //
	//screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	return img, err
}

// The main handler to deliver the image (can be called directly
func myHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Print("Hello, got a request\n")
	mimeWriter := multipart.NewWriter(w)
	mimeWriter.SetBoundary("--boundary")
	contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
	w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, pre-check=0, post-check=0, max-age=0")
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Connection", "close")
	s := time.Now()
	for {
		partHeader := make(textproto.MIMEHeader)
		partHeader.Add("Content-Type", "image/jpeg")
		partHeader.Add("X-StartTime", fmt.Sprintf("%v", s.Unix()))
		partHeader.Add("X-Timestamp", fmt.Sprintf("%v", s.Unix()))
		partWriter, _ := mimeWriter.CreatePart(partHeader)
		snapshot, _ := takeShot()
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, snapshot, nil)
		//storeImage(snapshot, "test.png")
		partWriter.Write(buf.Bytes())
	}

}

// small start function for the server (can be refactored in the main function)
func startServer() {
	s := &http.Server{
		Addr:              ":8080",
		Handler:           http.HandlerFunc(myHandler),
		TLSConfig:         nil,
		ReadTimeout:       10 * time.Hour, //.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      10 * time.Hour, //Second,
		IdleTimeout:       0,
		MaxHeaderBytes:    1 << 20,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	log.Fatal(s.ListenAndServe())
}

// serveCmd represents the serve command
func main() {
	fmt.Println("serve called")
	startServer()
	fmt.Print("finished")
}
