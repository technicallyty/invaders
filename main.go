package main

import (
	"flag"
	"fmt"
	"github.com/technicallyty/invasion/x"
	"os"
)

func main() {
	numAliens := flag.Int("n", 3, "number of aliens invading the x")
	fileName := flag.String("map", "map.txt", "file to use for map definition")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	world, err := x.LoadMapFromFile(file)
	if err != nil {
		panic(err)
	}

	err = world.SeedAliens(*numAliens)
	if err != nil {
		panic(err)
	}

	moved := true
	world.CheckBattleConditionsAndExec()
	for moved {
		moved = world.MoveAlien()
		if !moved {
			fmt.Println("nobody moved! exiting...")
			world.PrintCities()
		}
	}
}
