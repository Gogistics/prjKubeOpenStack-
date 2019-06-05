package main

import (
  "log"
  "flag"
  "apis"
  "net/http"
  "apis/config"
)

// msg type for websocket
type msg struct {
  Command string
}

// server entry
func main() {
  // handle request
  port := flag.String("p", "3001", "port")
  dir := flag.String("d", ".", "dir")
  flag.Parse()

  // REST APIs
  http.HandleFunc("/kube-apis/config", apis.DefaultHandler)
  http.HandleFunc("/kube-apis/config/", config.ConfigHandler)
  
  log.Printf("Serving %s Http port: %s\n", *dir, *port)
  log.Fatal(http.ListenAndServe(":" + *port, nil))
}