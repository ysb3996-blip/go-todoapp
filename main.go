package main

import (
	"log"
	"todoapp/app/controllers"
)

func main() {
	log.Fatalln(controllers.StartMainServer())
}
