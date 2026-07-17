package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/database"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/models"
)

var Cfg GlobalConfig

type GlobalConfig struct {
	ENV       string
	DB        *pgx.Conn
	DBQueries *database.Queries
	Logger    *slog.Logger
}

var ItemsList = make(map[string]*models.Item)

func init() {
	// Get Environmental Variables
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME not set")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER not set")
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		log.Fatal("DB_PASS not set")
	}

	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		log.Fatal("DB_ADDR not set")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT not set")
	}

	// Connect to DB
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbAddr,
		dbPort,
		dbName)

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v\n", err)
	}

	dbQueries := database.New(db)

	Cfg = GlobalConfig{
		ENV:       env,
		DB:        db,
		DBQueries: dbQueries,
	}

	// Initialize Items
	err = dbQueries.DeleteItems(context.Background())
	log.Println("Clearing items from DB...")
	if err != nil {
		log.Fatalf("Failed to clear items in DB: %v\n", err)
	}
	log.Println("Items cleared from DB")

	type ItemData struct {
		Description  string `json:"description"`
		EffectTarget string `json:"effect_target"`
		EffectValue  int32  `json:"effect_value"`
		Value        int32  `json:"value"`
	}

	log.Println("Loading Items...")
	itemFiles := []string{"Medicine", "Meal"}
	for _, itemFileName := range itemFiles {
		dbItemCategory, err := dbQueries.GetItemCategoriesByName(context.Background(), itemFileName)
		if err != nil {
			log.Fatalf("Failed to get category id: %v\n", err)
		}

		fileName := fmt.Sprintf("%s.json", itemFileName)
		filePath := path.Join("./internal/config", fileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to open file %s: %v\n", filePath, err)
		}
		var dataMap map[string]json.RawMessage
		err = json.Unmarshal(data, &dataMap)
		if err != nil {
			log.Fatalf("Failed to unmarshal data: %v\n", err)
		}

		for itemName, rawItemData := range dataMap {
			var itemData ItemData
			err := json.Unmarshal(rawItemData, &itemData)
			if err != nil {
				log.Fatalf("Failed to unmarshal item data: %v\n", err)
			}
			dbItem, err := dbQueries.CreateItems(context.Background(), database.CreateItemsParams{
				Name:        itemName,
				Description: itemData.Description,
				CategoryID:  dbItemCategory.ID,
				EffectTarget: pgtype.Text{
					String: itemData.EffectTarget,
					Valid:  (itemData.EffectTarget != ""),
				},
				EffectValue: pgtype.Int4{
					Int32: itemData.EffectValue,
					Valid: (itemData.EffectValue != 0),
				},
				Value: itemData.Value,
			})
			if err != nil {
				log.Fatalf("Failed to create item in DB: %v\n", err)
			}
			var itemCategory models.ItemCategory
			switch fileName {
			case "medicine":
				itemCategory = models.Medicine
			case "meal":
				itemCategory = models.Meal
			}
			item, err := models.NewItem(
				dbItem.Name,
				dbItem.Description,
				itemCategory,
				models.Effect{
					Target: dbItem.EffectTarget.String,
					Value:  dbItem.EffectValue.Int32,
				},
				dbItem.Value,
			)
			if err != nil {
				log.Fatalf("Failed to create item in memory: %v\n", err)
			}

			ItemsList[dbItem.Name] = item
		}

		log.Println("Loading Complete!")
	}
	for _, item := range ItemsList {
		log.Println("------- Item -------")
		log.Printf("Name: %s\n", item.GetName())
		log.Printf("Description: %s\n", item.GetDescription())
		log.Printf("Effect Target: %s\n", item.GetEffect().Target)
		log.Printf("Effect Value: %d\n", item.GetEffect().Value)
		log.Printf("Value: %d\n", item.GetValue())
		log.Println()
	}
}
