package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	server := http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	fmt.Printf("Server listening on %s", server.Addr)

	go http.ListenAndServe(":80", http.HandlerFunc(redirectHTTP))

	crt := "/etc/letsencrypt/live/mithyagames.com/fullchain.pem"
	key := "/etc/letsencrypt/live/mithyagames.com/privkey.pem"

	err := server.ListenAndServeTLS(crt, key)

	if err != nil {
		fmt.Println(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("request " + r.URL.Path)
	fmt.Fprintln(w, "Hello Mayank!")
}

func redirectHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Method : " + r.Method + " Port : " + r.Host +
		" URI : " + r.URL.RequestURI())

	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}

	target := "https://" + stripPort(r.Host) + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusFound)
}

func stripPort(hostport string) string {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}
	return net.JoinHostPort(host, "443")
}
