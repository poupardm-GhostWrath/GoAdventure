package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/config"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/initialization"
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
	assets := config.GlobalAssets{
		Items:     initItems,
		Locations: initLocations,
	}
	Assets = &assets

	// Create Test User
	err = createTestUser(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Login
	scanner := bufio.NewScanner(os.Stdin)
	userID, err := login(ctx, scanner)
	if err != nil {
		log.Fatal(err)
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
	fmt.Println("You are in the Starting Town!")
	fmt.Println("Notice: type 'exit' to exit.")
	fmt.Println("Notice: type 'help' for help menu.")

	for {
		fmt.Print("\nWhat would you like to do?: ")
		if scanner.Scan() {
			input := scanner.Text()
			exit, err := parseCommand(input)
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

func parseCommand(cmd string) (bool, error) {
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
`
		fmt.Println(help_menu)
		return false, nil
	case "go", "move":
		if len(parts) < 2 {
			return false, errors.New("Go where?")
		} else {
			return false, errors.New("Move feature not implemented.")
		}
	case "look":
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
		return false, nil
	case "inventory", "inv":
		fmt.Println("\n=== Inventory ===")
		for itemID, quantity := range Assets.Player.GetInventory() {
			item := Assets.Items[itemID]
			fmt.Printf("%s: %d\n", item.GetName(), quantity)
		}
		return false, nil
	case "stat":
		displayStat()
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
