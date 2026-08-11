package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

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

	// Execute migration SQL files
	migrationDir := "migrations"
	pattern := filepath.Join(migrationDir, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Failed to glob migration files: %v", err)
	}
	sort.Strings(files)

	for _, f := range files {
		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			log.Printf("Failed to read %s: %v", f, err)
			continue
		}
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			log.Printf("Migration %s: %v (may be already applied)", filepath.Base(f), err)
		} else {
			log.Printf("Migration %s: OK", filepath.Base(f))
		}
	}

	seed(db)
	fmt.Println("Migration completed")
}

func seed(db *gorm.DB) {
	var count int64
	db.Table("pools").Count(&count)
	if count == 0 {
		db.Exec("INSERT INTO pools (pool_type, total_supply, balance) VALUES ('public', 50000000000, 50000000000)")
		db.Exec("INSERT INTO pools (pool_type, total_supply, balance) VALUES ('official', 0, 0)")
		log.Println("Seeded: pools")
	}

	db.Table("reward_rules").Count(&count)
	if count == 0 {
		rewardRepo := repository.NewRewardRepo(db)
		rewardRepo.InsertDefaultRules(db)
		log.Println("Seeded: reward rules")
	}

	db.Table("trust_level_caps").Count(&count)
	if count == 0 {
		rewardRepo := repository.NewRewardRepo(db)
		rewardRepo.InsertDefaultCaps(db)
		log.Println("Seeded: trust level caps")
	}
}