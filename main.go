package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DevMehta22/mongoapi/routes"
)

func main()  {
	fmt.Println("MongoDB API")
	r := routes.Router()
	fmt.Println("Server is getting started...")
	log.Fatal(http.ListenAndServe(":4000",r))
	fmt.Println("Listening at port 4000!")
}