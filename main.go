package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/auth"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/config"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/development"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/initialization"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

var Cfg *config.GlobalConfig

var Assets *config.GlobalAssets

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	status := run(ctx, cancel)
	cancel()
	defer func() {
		if err := Cfg.DB.Close(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close db connection: %v\n", err)
		}
	}()
	os.Exit(status)
}

func run(ctx context.Context, cancel context.CancelFunc) int {
	// Initialize Global Config & Assets
	initCfg, err := initialization.InitializeConfig()
	if err != nil {
		log.Fatal(err)
	}
	Cfg = initCfg
	initItems, err := initialization.InitializeItems(Cfg.DBQueries)
	if err != nil {
		log.Fatal(err)
	}
	initLocations, err := initialization.InitializeLocations(Cfg.DBQueries)
	if err != nil {
		log.Fatal(err)
	}
	initStores, err := initialization.InitializeStore(initItems, initLocations)
	if err != nil {
		log.Fatal(err)
	}
	assets := config.GlobalAssets{
		Items:     initItems,
		Locations: initLocations,
		Stores:    initStores,
	}
	Assets = &assets

	// Create Test User
	if Cfg.ENV == "development" {
		err = development.CreateTestUser(ctx, Cfg.DBQueries)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Login
	scanner := bufio.NewScanner(os.Stdin)
	var userID uuid.UUID
	fmt.Println("\n=== GoAdventure ===")
	fmt.Println(" 1. Login")
	fmt.Println(" 2. Register")
	fmt.Println(" 3. Exit")
outer:
	for {
		fmt.Print(" Choice: ")
		if scanner.Scan() {
			input, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				fmt.Println("Invalid choice. Please try again.")
				continue
			}
			switch input {
			case 1:
				userID, err = auth.Login(ctx, Cfg.DBQueries, scanner)
				if err != nil {
					log.Fatal(err)
				}
				break outer
			case 2:
				err = auth.Register(ctx, Cfg.DBQueries, scanner)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Registration successful!")
				userID, err = auth.Login(ctx, Cfg.DBQueries, scanner)
				if err != nil {
					log.Fatal(err)
				}
				break outer
			case 3:
				return 0
			default:
				fmt.Println("Invalid choice. Please try again.")
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// Get Player
	exit, err := getPlayer(ctx, scanner, userID)
	if err != nil {
		log.Fatal(err)
	}
	if exit {
		return 0
	}

	fmt.Printf("\n======== Welcome %s =======\n", Assets.Player.GetName())
	fmt.Println("Notice: type 'exit' to exit.")
	fmt.Println("Notice: type 'help' for help menu.")

	for {
		fmt.Print("\nWhat would you like to do?: ")
		if scanner.Scan() {
			input := scanner.Text()
			exit, err := parseCommand(scanner, input)
			if err != nil {
				fmt.Println(err)
			}
			if exit {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	err = savePlayer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Done()
	_, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return 0
}

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
| Take        | take              |
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
		fmt.Println("\n=== Inventory ===")
		for itemID, quantity := range Assets.Player.GetInventory() {
			item := Assets.Items[itemID]
			fmt.Printf("%s: %d\n", item.GetName(), quantity)
		}
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
	if Assets.Locations[Assets.Player.GetLocation()].HasStore() {
		fmt.Println("You see a store in the corner.")
	}
	directions := Assets.Locations[Assets.Player.GetLocation()].GetDirections()
	for _, direction := range directions {
		fmt.Printf("You see %s to the %s.\n", Assets.Locations[direction.GetLocationID()].GetName(), direction.GetDirection())
	}
}

func store(scanner *bufio.Scanner) error {
	store := Assets.Stores[Assets.Player.GetLocation()]
	fmt.Printf("\n=== %s ===\n", store.GetName())
	fmt.Println(" 1. Check Store Inventory")
	fmt.Println(" 2. Check Player Inventory")
	fmt.Println(" 3. Buy Item")
	fmt.Println(" 4. Sell Item")
	fmt.Println(" 5. Exit")
outer:
	for {
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
		for itemID, quantity := range v.GetInventory() {
			fmt.Printf(" %s: %d\n", Assets.Items[itemID].GetName(), quantity)
		}
	default:
		return
	}
}
