package networking

import (
	"MDPreview/internal"
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"
	"os/exec"
	"net"
	"net/http"
	"strings"
)

const LocalHost = "localhost"
const ChannelPort = ":8080"
const WebServerPort = ":8000"

var initLayout string

var messageChannel = make(chan string)
var scrollChannel  = make(chan string)

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

func launchBrowser() {
	// TODO: list of browsers in case $BROWSER env var isn't set
	// display error message in neovim if no browser is found
	val, isSet := os.LookupEnv("BROWSER")
	if isSet {
		exec.Command(val, LocalHost + WebServerPort).Start()
	}
}

func httpServer() {
	http.HandleFunc("/", index)
	http.HandleFunc("/events", events)

	// TODO: simulate lag; remove
	// time.Sleep(14500*time.Millisecond)

	err := http.ListenAndServe(WebServerPort, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Failed to connect; server closed\n")
	}
}

func waitForServer(address string, retries int, delay time.Duration) error {
	for i := 0; i < retries; i++ {
		conn, err := net.DialTimeout("tcp", address, delay)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("server not available at %s after %d retries", address, retries)
}

func StartServer() {
	ln, err := net.Listen("tcp", ChannelPort)
	if err != nil {
		fmt.Printf("error: failed to listen on tcp port %s: %s\n", ChannelPort, err)
	}
	fmt.Printf("\nsig_start\n")

	go httpServer()

	er := waitForServer(LocalHost + WebServerPort, 10, 250*time.Millisecond)
	if er != nil {
		fmt.Printf("Failed to connect to webserver: %s\n", er)
		os.Exit(1)
	}
	launchBrowser()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error: failed to accept tcp connection: %s\n", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Bytes()
		json := internal.EventsJSON(message)
		html := internal.ToMarkdown(json.Content)

		switch json.Event {
		case "init":
			initLayout = internal.LayoutPage(html)
		case "scroll":
			scrollChannel<-json.Content
		case "reload":
			messageChannel<-html
		}
	}
}
