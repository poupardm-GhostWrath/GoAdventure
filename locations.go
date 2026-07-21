package main

import "fmt"

func Look() {
	fmt.Println("\nYou look around...")
	fmt.Printf("You are currently in %s.\n", Assets.Locations[Assets.Player.GetLocation()].GetName())
	fmt.Println(Assets.Locations[Assets.Player.GetLocation()].GetDescription())
	if Assets.Locations[Assets.Player.GetLocation()].HasStore() {
		fmt.Println("You see a store in the corner.")
	}
	directions := Assets.Locations[Assets.Player.GetLocation()].GetDirections()
	for _, direction := range directions {
		fmt.Printf("You see %s to the %s.\n", Assets.Locations[direction.TargetLocationID].GetName(), direction.Direction)
	}
}
