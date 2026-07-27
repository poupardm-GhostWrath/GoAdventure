package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

func parseCommand(scanner *bufio.Scanner, cmd string) (bool, error) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(cmd)), " ")

	if len(parts) == 0 {
		return false, errors.New("invalid command")
	}
	verb := parts[0]
	if verb == "exit" || verb == "quit" {

	}
	switch verb {
	case "exit", "quit":
		fmt.Println("Thank you for playing. See you next time.")
		return true, nil
	case "help":
		fmt.Println("=== Help Menu ===")
		help_menu := `
+-------------+-------------------+
| Action      | Command           |
+-------------+-------------------+
| Exit        | exit              |
|             | quit              |
+-------------+-------------------+
| Help        | help              |
+-------------+-------------------+
| Move        | move <direction>  |
|             | go <direction>    |
+-------------+-------------------+
| Look        | look              |
+-------------+-------------------+
| Inventory   | inventory         |
|             | inv               |
+-------------+-------------------+
| Stat        |	stat              |
+-------------+-------------------+
| Store       | store             |
+-------------+-------------------+
`
		fmt.Println(help_menu)
		return false, nil
	case "go", "move":
		if len(parts) < 2 {
			return false, errors.New("Go where?")
		} else {
			for _, direction := range Assets.Locations[Assets.Player.GetLocation()].GetDirections() {
				if direction.GetDirection() == parts[1] {
					fmt.Printf("Moving to %s...\n", Assets.Locations[direction.GetLocationID()].GetName())
					Assets.Player.SetLocation(direction.GetLocationID())
					look()
					return false, nil
				}
			}
			return false, errors.New("Invalid Direction.")
		}
	case "look":
		look()
		return false, nil
	case "inventory", "inv":
		displayInventory(Assets.Player)
		return false, nil
	case "stat":
		Assets.Player.DisplayStats()
		return false, nil
	case "store":
		if !Assets.Locations[Assets.Player.GetLocation()].HasStore() {
			return false, errors.New("This area doesn't have a store.")
		}
		err := store(scanner)
		if err != nil {
			return false, err
		}
		return false, nil
	case "eat":
		if len(parts) < 2 {
			return false, errors.New("usage: eat <item name>")
		}
		var itemName string
		if len(parts) > 2 {
			itemName = strings.Join(parts[1:], " ")
		}
		var itemID int32
		for key, item := range Assets.Items {
			if strings.ToLower(item.GetName()) == itemName {
				itemID = key
				break
			}
		}
		err := Assets.Player.EatFood(Assets.Items, itemID)
		if err != nil {
			return false, err
		}
		return false, nil
	case "test":
		if Cfg.ENV == "development" {
			if len(parts) < 2 {
				return false, errors.New("usage: test <target> [itemID] [amount]")
			}
			var amount int32 = 1
			switch parts[1] {
			case "gold":
				if len(parts) > 2 {
					num, err := strconv.ParseInt(parts[2], 10, 32)
					if err != nil {
						return false, err
					}
					amount = int32(num)
				}
				Assets.Player.AddGold(int32(amount))
				fmt.Printf("Added %d gold.\n", amount)
			case "item":
				if len(parts) < 3 {
					return false, errors.New("usage: test item <itemID> [amount]")
				}
				if len(parts) > 3 {
					num, err := strconv.ParseInt(parts[3], 10, 32)
					if err != nil {
						return false, err
					}
					amount = int32(num)
				}
				itemID, err := strconv.ParseInt(parts[2], 10, 32)
				if err != nil {
					return false, err
				}
				Assets.Player.AddItem(int32(itemID), int32(amount))
				fmt.Printf("Added %d %s.\n", amount, Assets.Items[int32(itemID)].GetName())
			}
		}
		return false, nil
	default:
		return false, errors.New("Invalid command")
	}
}

func savePlayer(ctx context.Context) error {
	err := Cfg.DBQueries.UpdatePlayerByID(ctx, database.UpdatePlayerByIDParams{
		ID:           Assets.ID,
		CurrentExp:   Assets.Player.GetCurrentExp(),
		CurrentLevel: Assets.Player.GetLevel(),
		Gold:         Assets.Player.GetLevel(),
		LocationID:   Assets.Player.GetLocation(),
	})
	if err != nil {
		return fmt.Errorf("failed to save player: %v\n", err)
	}
	err = Cfg.DBQueries.DeleteInventoryItemByPlayerID(ctx, Assets.ID)
	if err != nil {
		return fmt.Errorf("failed to delete inventory for player: %v\n", err)
	}
	for itemID, quantity := range Assets.Player.GetInventory() {
		err = Cfg.DBQueries.CreateInventoryItem(ctx, database.CreateInventoryItemParams{
			ItemID:   itemID,
			PlayerID: Assets.ID,
			Quantity: quantity,
		})
		if err != nil {
			return fmt.Errorf("failed to save inventory: %v\n", err)
		}
	}
	return nil
}

func look() {
	fmt.Println("\nYou look around...")
	fmt.Printf("You are currently in %s.\n", Assets.Locations[Assets.Player.GetLocation()].GetName())
	fmt.Println(Assets.Locations[Assets.Player.GetLocation()].GetDescription())
	directions := Assets.Locations[Assets.Player.GetLocation()].GetDirections()
	for _, direction := range directions {
		fmt.Printf("You see %s to the %s.\n", Assets.Locations[direction.GetLocationID()].GetName(), direction.GetDirection())
	}
	if Assets.Locations[Assets.Player.GetLocation()].HasStore() {
		fmt.Println("You see a store in the corner.")
	} else {
		fmt.Println("Enemies:")
		if len(Assets.EnemiesLocation[Assets.Player.GetLocation()]) != 0 {
			for enemyID, quantity := range Assets.EnemiesLocation[Assets.Player.GetLocation()] {
				fmt.Printf("You see %d %s.\n", quantity, Assets.Enemies[enemyID].GetName())
			}
		} else {
			fmt.Println("You don't see any enemies.")
			if generateJoke() {
				fmt.Println("When is the last time you took a bath...")
			}
		}
	}
}

func store(scanner *bufio.Scanner) error {
	store := Assets.Stores[Assets.Player.GetLocation()]
	fmt.Printf("\n=== %s ===", store.GetName())

outer:
	for {
		fmt.Println()
		fmt.Println(" 1. Check Store Inventory")
		fmt.Println(" 2. Check Player Inventory")
		fmt.Println(" 3. Buy Item")
		fmt.Println(" 4. Sell Item")
		fmt.Println(" 5. Exit")
		fmt.Print(" Choice: ")
		if scanner.Scan() {
			input, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				fmt.Println("Invalid choice. Please try again.")
				continue
			}
			switch input {
			case 1:
				displayInventory(store)
			case 2:
				displayInventory(Assets.Player)
			case 3:
				fmt.Println("Feature not implemented.")
			case 4:
				fmt.Println("Feature not implemented.")
			case 5:
				break outer
			default:
				fmt.Println("Invalid choice. Please try again.")
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

func displayInventory(t any) {
	switch v := t.(type) {
	case *models.Store:
		fmt.Printf("\n=== %s Inventory ===\n", v.GetName())
		for itemID, quantity := range v.GetInventory() {
			fmt.Printf(" %s: %d\n", Assets.Items[itemID].GetName(), quantity)
		}
	case *models.Player:
		fmt.Printf("\n=== %s Inventory ===\n", v.GetName())
		if len(v.GetInventory()) == 0 {
			fmt.Println("Your inventory is empty...")
			if generateJoke() {
				time.Sleep(2 * time.Second)
				fmt.Println("Wait!")
				time.Sleep(5 * time.Second)
				fmt.Println("Never mind still empty...")
			}
		}
		for itemID, quantity := range v.GetInventory() {
			fmt.Printf(" %s: %d\n", Assets.Items[itemID].GetName(), quantity)
		}
	default:
		return
	}
}

func generateJoke() bool {
	randNum := rand.IntN(100) + 1
	if randNum < 26 {
		return true
	}
	return false
}
