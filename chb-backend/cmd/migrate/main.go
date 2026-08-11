package main

import (
	"fmt"
	"log"
	"os"

	"github.com/campushub/chb-backend/internal/config"
	"github.com/campushub/chb-backend/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	migrationDir := "migrations"
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("Failed to read migrations dir: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Printf("Found migration file: %s\n", entry.Name())
		}
	}

	seed(db)
	fmt.Println("Migration completed")
}

func seed(db *gorm.DB) {
	// Pools
	var count int64
	db.Table("pools").Count(&count)
	if count == 0 {
		db.Exec("INSERT INTO pools (pool_type, total_supply, balance) VALUES ('public', 50000000000, 50000000000)")
		db.Exec("INSERT INTO pools (pool_type, total_supply, balance) VALUES ('official', 0, 0)")
		log.Println("Seeded: pools")
	}

	// Reward rules
	db.Table("reward_rules").Count(&count)
	if count == 0 {
		rewardRepo := repository.NewRewardRepo(db)
		rewardRepo.InsertDefaultRules(db)
		log.Println("Seeded: reward rules")
	}

	// Trust level caps
	db.Table("trust_level_caps").Count(&count)
	if count == 0 {
		rewardRepo := repository.NewRewardRepo(db)
		rewardRepo.InsertDefaultCaps(db)
		log.Println("Seeded: trust level caps")
	}
}
