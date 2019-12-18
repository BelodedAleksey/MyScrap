package main

import (
	"fmt"

	"./autochess"
)

func main() {
	//Test autochess
	fmt.Println("TOP500")
	players, err := autochess.GetTop500()
	if err != nil {
		fmt.Printf("Error GetTop500: %s", err)
	}
	for _, p := range players {
		fmt.Println("Игрок: ", p.Name+"\t"+p.ID+"\t"+p.Rank+"\t"+p.MMR+"\t"+p.IconURL)
	}

	fmt.Println("SEARCH KIRITO")
	players, err = autochess.GetPlayersByName("kirito")
	if err != nil {
		fmt.Printf("Error GetPlayersByName: %s", err)
	}
	for _, p := range players {
		fmt.Println("Игрок: ", p.Name+"\t"+p.ID+"\t"+p.Rank+"\t"+p.IconURL)
	}
}
