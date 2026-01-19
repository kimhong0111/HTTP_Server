package main

import (
	"log"
	"net/http"
	"myApp/server"
)



func main(){
  server:=&server.BankingServer{
	Request: server.NewMemoryStorage()}
  log.Fatal(http.ListenAndServe(":3000",server))
}