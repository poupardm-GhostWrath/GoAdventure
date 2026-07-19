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
	"github.com/howeyc/gopass"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/auth"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/config"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
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
	initCfg, err := initialization.InitializeConfig()
	if err != nil {
		log.Fatal(err)
	}
	Cfg = initCfg
	initItems, err := initialization.InitializeItems(Cfg.DBQueries)
	if err != nil {
		log.Fatal(err)
	}
	assets := config.GlobalAssets{
		Items: initItems,
	}
	Assets = &assets

	// Create Test User
	_, err = Cfg.DBQueries.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		hash, err := auth.HashPassword("test")
		if err != nil {
			log.Fatal(err)
		}
		err = Cfg.DBQueries.CreateUser(ctx, database.CreateUserParams{
			Email:        "test@example.com",
			PasswordHash: hash,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	var failedAttempts int32
	var userID uuid.UUID
	for {
		var email string
		fmt.Print("Email: ")
		if scanner.Scan() {
			email = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Print("Password: ")
		password, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal(err)
		}
		dbUser, err := Cfg.DBQueries.GetUserByEmail(ctx, email)
		if err != nil {
			log.Fatal(err)
		}
		match, err := auth.CheckPasswordHash(string(password), dbUser.PasswordHash)
		if err != nil {
			log.Fatal(err)
		}
		if !match {
			failedAttempts += 1
		} else {
			userID = dbUser.ID
			break
		}
		if failedAttempts == 3 {
			log.Fatal("Too many attempts.")
		}
	}
	dbPlayers, err := Cfg.DBQueries.GetPlayersByUserID(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}
	amountOfCharacters := 0
	fmt.Println("=== Character Selection ===")
	for _, dbPlayer := range dbPlayers {
		fmt.Printf("%d. %s\n", amountOfCharacters+1, dbPlayer.Name)
		amountOfCharacters += 1
	}
	fmt.Printf("%d. Create New Character\n", amountOfCharacters+1)
	fmt.Printf("%d. Exit\n", amountOfCharacters+2)
	for {
		fmt.Print("Input: ")
		if scanner.Scan() {
			input, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				log.Fatal(err)
			}
			if input > 0 && input < (int64(amountOfCharacters+3)) {
				if input == int64(amountOfCharacters+1) {
					// Create Player
					fmt.Print("Character Name: ")
					if scanner.Scan() {
						charName := scanner.Text()
						dbPlayer, err := Cfg.DBQueries.CreatePlayer(ctx, database.CreatePlayerParams{
							UserID: userID,
							Name:   charName,
						})
						inventory := make(map[int32]int32)
						player, err := models.NewPlayer(dbPlayer.ID, dbPlayer.Name, dbPlayer.CurrentExp.Int32, dbPlayer.CurrentLevel.Int32, dbPlayer.Gold.Int32, inventory)
						if err != nil {
							log.Fatal(err)
						}
						Assets.Player = player
						break
					}
					if err := scanner.Err(); err != nil {
						log.Fatal(err)
					}
				}
				if input == int64(amountOfCharacters+2) {
					return 0
				}
				dbPlayer, err := Cfg.DBQueries.GetPlayersByID(ctx, dbPlayers[input-1].ID)
				if err != nil {
					log.Fatal(err)
				}
				dbInventory, err := Cfg.DBQueries.GetInventoryByPlayerID(ctx, dbPlayer.ID)
				if err != nil {
					log.Fatal(err)
				}
				inventory := make(map[int32]int32)
				for _, dbInventoryItem := range dbInventory {
					inventory[dbInventoryItem.ItemID] = dbInventoryItem.Quantity
				}

				player, err := models.NewPlayer(dbPlayer.ID, dbPlayer.Name, dbPlayer.CurrentExp.Int32, dbPlayer.CurrentLevel.Int32, dbPlayer.Gold.Int32, inventory)
				if err != nil {
					log.Fatal(err)
				}
				Assets.Player = player
				break
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("======== Welcome %s =======\n", Assets.Player.GetName())
	fmt.Println("You are in the Starting Town!")
	fmt.Println("Notice: type 'exit' to exit.")
	fmt.Println("Notice: type 'help' for help menu.")

	for {
		fmt.Print("What would you like to do?: ")
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
		return false, errors.New("Look feature not implemented.")
	case "inventory", "inv":
		fmt.Println("=== Inventory ===")
		for itemID, quantity := range Assets.Player.GetInventory() {
			item := Assets.Items[itemID]
			fmt.Printf("%s: %d\n", item.GetName(), quantity)
		}
		return false, nil
	default:
		return false, errors.New("Invalid command")
	}
}
