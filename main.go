package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/poupardm-GhostWrath/GoAdventure/internal/config"
	"github.com/poupardm-GhostWrath/GoAdventure/internal/initialization"
)

var Cfg *config.GlobalConfig

var Assets *config.GlobalAssets

func main() {
	fmt.Println("Hello World")
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
	fmt.Println("Loading...")
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
	fmt.Println("Loading Complete!")
	for _, item := range Assets.Items {
		fmt.Println("------ Item ------")
		fmt.Printf("Name: %s\n", item.GetName())
		fmt.Printf("Description: %s\n", item.GetDescription())
		fmt.Printf("Category: %s\n", item.GetCategory())
		fmt.Printf("Effect Description: %s\n", item.GetEffect().Description)
		fmt.Printf("Effect Target: %s\n", item.GetEffect().Target)
		fmt.Printf("Effect Value: %d\n", item.GetEffect().Value)
		fmt.Printf("Value: %d\n", item.GetValue())
		fmt.Println()
	}
	<-ctx.Done()
	_, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return 0
}
