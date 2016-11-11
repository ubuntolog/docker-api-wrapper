package main

import (
	"fmt"
	//"io"
	"net"
	"net/http"
	"os"
	"io"
)

func dieIf(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

const sock = "/var/run/docker.sock"

func server() {
	l, err := net.Listen("unix", sock)
	dieIf(err)
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello: %v\n", r.RequestURI)
		}),
	}
	srv.Serve(l)
}

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", sock)
}

func main() {
	defer os.Remove(sock)
	//go server()

	tr := &http.Transport{
		Dial: fakeDial,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("http://localhost:4243/containers/json")
	dieIf(err)
	io.Copy(os.Stdout, resp.Body)
	resp.Body.Close()

}
