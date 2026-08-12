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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
+--------------+--------------------+
| Action       | Command            |
+--------------+--------------------+
| Exit         | exit               |
|              | quit               |
+--------------+--------------------+
| Help Menu    | help               |
+--------------+--------------------+
| Move         | move <direction>   |
|              | go <direction>     |
+--------------+--------------------+
| Look Around  | look               |
+--------------+--------------------+
| Display      | inventory          |
| Inventory    | inv                |
+--------------+--------------------+
| Display Stat | stat               |
+--------------+--------------------+
| Display Gear | gear               |
+--------------+--------------------+
| Equip Gear   | equip <itemName>   |
+--------------+--------------------+
| Unequip Gear | unequip <gearSlot> |
+--------------+--------------------+
| Access Store | store              |
+--------------+--------------------+
| Clear Screen | clear              |
+--------------+--------------------+
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
	case "gear":
		Assets.Player.DisplayGear(Assets.Items)
		return false, nil
	case "equip":
		if len(parts) < 2 {
			return false, errors.New("Equip what? 'usage: equip <item name>'")
		}
		itemName := strings.Join(parts[1:], " ")
		var itemID int32
		var found = false
		for key, item := range Assets.Items {
			if strings.ToLower(item.GetName()) == itemName {
				itemID = key
				found = true
				break
			}
		}
		if !found {
			return false, errors.New("Invalid item name")
		}
		switch Assets.Items[itemID].GetCategory() {
		case "Weapon":
			err := Assets.Player.EquipWeapon(Assets.Items, itemID)
			if err != nil {
				return false, err
			}
		case "Armor":
			err := Assets.Player.EquipArmor(Assets.Items, itemID)
			if err != nil {
				return false, err
			}
		case "Accessory":
			err := Assets.Player.EquipAccessory(Assets.Items, itemID)
			if err != nil {
				return false, err
			}
		default:
			return false, errors.New("Item is not an equipable item.")
		}
		fmt.Printf("%s was equipped.\n", Assets.Items[itemID].GetName())
		return false, nil
	case "unequip":
		if len(parts) < 2 {
			return false, errors.New("Un-equip what? 'usage: unequip <item slot>'")
		}
		switch parts[1] {
		case "weapon":
			err := Assets.Player.UnequipWeapon()
			if err != nil {
				return false, err
			}
		case "armor":
			err := Assets.Player.UnequipArmor()
			if err != nil {
				return false, err
			}
		case "accessory":
			err := Assets.Player.UnequipAccessory()
			if err != nil {
				return false, err
			}
		default:
			return false, errors.New("Invalid slot: Valid inputs are 'weapon', 'armor' or 'accessory'")
		}
		fmt.Printf("You un-equipped your %s.\n", parts[1])
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
	case "clear":
		ClearScreen()
		return false, nil
	default:
		return false, errors.New("Invalid command")
	}
}

func savePlayer(ctx context.Context) error {
	err := Cfg.DBQueries.UpdatePlayerByID(ctx, database.UpdatePlayerByIDParams{
		ID:            Assets.ID,
		CurrentExp:    Assets.Player.GetCurrentExp(),
		CurrentLevel:  Assets.Player.GetLevel(),
		Gold:          Assets.Player.GetGold(),
		LocationID:    Assets.Player.GetLocation(),
		WeaponGear:    Assets.Player.GetWeapon(),
		ArmorGear:     Assets.Player.GetArmor(),
		AccessoryGear: Assets.Player.GetAccessory(),
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
	fmt.Println("\n=== Store ===")
	fmt.Printf("Welcome to %s!\n", store.GetName())
	fmt.Println("What can I get you today?")
	store.RefreshStore(Assets.Items)
outer:
	for {
		fmt.Println()
		fmt.Print("What are we doing?> ")
		if scanner.Scan() {
			input := scanner.Text()
			parts := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
			switch parts[0] {
			case "inv":
				if len(parts) > 1 {
					if parts[1] == "store" {
						displayInventory(store)
					} else {
						displayInventory(Assets.Player)
					}
				} else {
					displayInventory(Assets.Player)
				}
			case "buy":
				if len(parts) < 3 {
					fmt.Println("I'm sorry. What are you trying to buy exactly and how much?")
					continue
				}
				quantity, err := strconv.ParseInt(parts[len(parts)-1], 10, 32)
				if err != nil || quantity < 1 {
					fmt.Println("Sorry. How many did you want to buy?")
					continue
				}
				itemName := strings.Join(parts[1:len(parts)-1], " ")
				itemFound := false
				for key, item := range Assets.Items {
					if strings.ToLower(item.GetName()) == itemName {
						itemFound = true
						n, err := store.BuyItem(Assets.Items, key, int32(quantity), Assets.Player)
						if err != nil {
							fmt.Println(err)
							break
						}
						fmt.Printf("You have bought %d %s.\n", n, item.GetName())
					}
				}
				if !itemFound {
					fmt.Printf("I'm sorry. I don't know what '%s' is.\n", cases.Title(language.English).String(itemName))
				}
			case "sell":
				if len(parts) < 3 {
					fmt.Println("I'm sorry. What are you trying to sell exactly and how much?")
					continue
				}
				quantity, err := strconv.ParseInt(parts[len(parts)-1], 10, 32)
				if err != nil || quantity < 1 {
					fmt.Println("Sorry. How may did you want to sell?")
					continue
				}
				itemName := strings.Join(parts[1:len(parts)-1], " ")
				itemFound := false
				for key, item := range Assets.Items {
					if strings.ToLower(item.GetName()) == itemName {
						itemFound = true
						n, err := store.SellItem(Assets.Items, key, int32(quantity), Assets.Player)
						if err != nil {
							fmt.Println(err)
							break
						}
						fmt.Printf("You have sold %d %s for %d gold.\n", quantity, item.GetName(), n)
					}
				}
				if !itemFound {
					fmt.Printf("I'm sorry. I don't know what '%s' is.\n", cases.Title(language.English).String(itemName))
				}

			case "help":
				help_menu := `
+----------------------------+--------------------------+
| Command                    | Description              |
+----------------------------+--------------------------+
| inv [store]                | Check inventory          |
+----------------------------+--------------------------+
| buy <itemName> <quantity>  | Buy from the store       |
+----------------------------+--------------------------+
| sell <itemName> <quantity> | Sell from your inventory |
+----------------------------+--------------------------+
| exit                       | Exit store               |
+----------------------------+--------------------------+
`
				fmt.Println(help_menu)
			case "exit":
				fmt.Println("Of course. Come back again!")
				break outer
			default:
				fmt.Println("I'm sorry I don't understand.")
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
		fmt.Printf(" Gold: %d\n", v.GetGold())
		for itemID, quantity := range v.GetInventory() {
			fmt.Printf(" %s: %d\n", Assets.Items[itemID].GetName(), quantity)
		}
	case *models.Player:
		fmt.Printf("\n=== %s Inventory ===\n", v.GetName())
		if v.GetGold() == 0 {
			fmt.Println(" Gold: Broke")
		} else {
			fmt.Printf(" Gold: %d\n", v.GetGold())
		}
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
