package networking

import (
	"MDPreview/internal"
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var initLayout string

var messageChannel = make(chan string)
var scrollChannel = make(chan string)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	spaData := strings.ReplaceAll(initLayout, "\t", "    ")
    fmt.Fprintf(w, spaData)
}

func events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no");
	w.Header().Set("Cache-Control", "no-cache");
	w.Header().Set("Connection", "keep-alive")

	var sendData = func(d string, ev string) {
		escData := strings.ReplaceAll(d, "\n", "<br>")
		spaData := strings.ReplaceAll(escData, "\t", "    ")

		fmt.Fprintf(w, "event: %s\n", ev)
		fmt.Fprintf(w, "data: %s\n\n", spaData)

		w.(http.Flusher).Flush()
	}

	for {
		select {
			case msg := <-messageChannel:
				sendData(msg, "reload")

			case scroll := <-scrollChannel:
				sendData(scroll, "scroll")

			case <-r.Context().Done():
				return
		}
	}
}

func httpServer() {
	http.HandleFunc("/", index)
	http.HandleFunc("/events", events)

	err := http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server Closed\n")
	}
}

func StartServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("error: failed to listen on tcp port 8080: %s", err)
	}

	fmt.Println("sig_start")
	go httpServer()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error: failed to accept tcp connection: %s", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	val, isSet := os.LookupEnv("BROWSER")
	if isSet {
		browser := exec.Command(val, "127.0.0.1:8000")
		browser.Start()
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Bytes()
		json := internal.EventsJSON(message)

		switch json.Event {
		case "init":
			html := internal.ToMarkdown(json.Content)
			fmt.Printf("%s\n", html)
			initLayout = internal.LayoutPage(html)
		case "scroll":
			fmt.Printf("scrolling: %s\n", json.Content)
			scrollChannel<-json.Content
		case "reload":
			html := internal.ToMarkdown(json.Content)
			fmt.Printf("%s\n", html)
			messageChannel<-html
		}
	}
}
