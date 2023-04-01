package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ganeshdipdumbare/alien_attack/attack"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		log.Println("Please provide no of aliens as an argument")
		return
	}

	noOfAliens, err := strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		log.Println("Please provide valid no of aliens as an argument")
		return
	}

	world, err := attack.CreateWorld("cities.txt")
	if err != nil {
		log.Println("Error while creating world: ", err)
		return
	}
	if noOfAliens > world.GetNoOfCities() {
		log.Println("No of aliens should be less than or equal to no of cities")
		return
	}

	fmt.Println("World before attack: ")
	world.PrintMap()

	world.UnleashAliens(noOfAliens)
	fmt.Println()

	fmt.Println("World after attack: ")
	world.PrintMap()
}
