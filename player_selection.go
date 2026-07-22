package main

import (
	"bufio"
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

type player struct {
	ID   uuid.UUID
	Name string
}

func createPlayer(ctx context.Context, scanner *bufio.Scanner, userID uuid.UUID) (uuid.UUID, error) {
	var playerID uuid.UUID
	fmt.Println("\n==== Player Creation ====")
	for {
		fmt.Print(" Name: ")
		if scanner.Scan() {
			input := scanner.Text()
			if input != "" {
				dbUserID, err := Cfg.DBQueries.CreatePlayer(ctx, database.CreatePlayerParams{
					Name:   input,
					UserID: userID,
				})
				if err != nil {
					return playerID, err
				}
				playerID = dbUserID
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return playerID, err
		}
	}
	return playerID, nil
}

func getPlayer(ctx context.Context, scanner *bufio.Scanner, userID uuid.UUID) (bool, error) {
	// Get Players by UserID
	players, err := getPlayers(ctx, userID)
	if err != nil {
		return false, err
	}

	// Select Player
	var playerID uuid.UUID
	fmt.Println("\n==== Player Selection ====")
	for i, player := range players {
		fmt.Printf(" %d. %s\n", i+1, player.Name)
	}
	fmt.Printf(" %d. Create Character\n", len(players)+1)
	fmt.Printf(" %d. Exit\n", len(players)+2)
	for {
		fmt.Printf("Enter Selection (1-%d): ", (len(players) + 2))
		if scanner.Scan() {
			input, err := strconv.ParseInt(scanner.Text(), 10, 32)
			if err != nil {
				return false, err
			}
			if input > 0 || input < int64(len(players)+3) {
				// Exit
				if input == int64(len(players)+2) {
					return true, nil
				}

				// Create Character
				if input == int64(len(players)+1) {
					pID, err := createPlayer(ctx, scanner, userID)
					if err != nil {
						return false, err
					}
					playerID = pID
					break
				}

				// Selected Character
				playerID = players[input-1].ID
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return false, err
		}
	}
	dbPlayer, err := Cfg.DBQueries.GetPlayersByID(ctx, playerID)
	if err != nil {
		return false, err
	}
	inventory, err := getInventory(ctx, dbPlayer.ID)
	if err != nil {
		return false, err
	}
	player, err := models.InitPlayer(dbPlayer.ID, dbPlayer.Name, dbPlayer.CurrentExp, dbPlayer.CurrentLevel, dbPlayer.Gold, dbPlayer.LocationID, inventory)
	if err != nil {
		return false, err
	}
	Assets.ID = playerID
	Assets.Player = player
	return false, nil
}

func getPlayers(ctx context.Context, userID uuid.UUID) ([]player, error) {
	var players []player
	dbPlayers, err := Cfg.DBQueries.GetPlayersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, dbPlayer := range dbPlayers {
		p := player{
			ID:   dbPlayer.ID,
			Name: dbPlayer.Name,
		}
		players = append(players, p)
	}
	return players, nil
}
