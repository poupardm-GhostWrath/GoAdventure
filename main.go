package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/auth"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/config"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/development"
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

//lint:ignore SA4009 Not in use
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
	initEnemies, err := initialization.InitializeEnemies(Cfg.DBQueries)
	if err != nil {
		log.Fatal(err)
	}
	initEnemiesLocation, err := initialization.InitializeEnemiesLocation(initEnemies, initLocations)
	if err != nil {
		log.Fatal(err)
	}

	assets := config.GlobalAssets{
		Items:           initItems,
		Locations:       initLocations,
		Stores:          initStores,
		Enemies:         initEnemies,
		EnemiesLocation: initEnemiesLocation,
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
	fmt.Println("\n=== Go Adventure ===")
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

	ClearScreen()

	fmt.Printf("======== Welcome %s =======\n", Assets.Player.GetName())
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

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
