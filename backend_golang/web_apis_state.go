package main

import (
  "log"
  "flag"
  "apis"
  "net/http"
  "apis/state"
)

// msg type for websocket
type msg struct {
  Command string
}

// server entry
func main() {
  // handle request
  port := flag.String("p", "3002", "port")
  dir := flag.String("d", ".", "dir")
  flag.Parse()

  // REST APIs
  http.HandleFunc("/kube-apis/state", apis.DefaultHandler)
  http.HandleFunc("/kube-apis/state/", state.StateHandler)
  
  log.Printf("Serving %s Http port: %s\n", *dir, *port)
  log.Fatal(http.ListenAndServe(":" + *port, nil))
}