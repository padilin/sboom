package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(databasePath string) {
	// assign to package-level db (avoid shadowing)
	db = OpenDB(databasePath)

	if err := db.AutoMigrate(
		&Repo{},
		&Release{},
		&Dependency{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func OpenDB(databasePath string) *gorm.DB {
	d, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	db = d
	return d
}

func GetDB() *gorm.DB {
	return db
}

func GetRepos(limit int) ([]Repo, error) {
	var repos []Repo
	ctx := context.Background()
	result := db.WithContext(ctx).
		Preload("Releases").
		Preload("Releases.Dependencies").
		Limit(limit).
		Find(&repos)
	return repos, result.Error
}

func GetRepo(id uint) (Repo, error) {
	var repo Repo
	ctx := context.Background()
	result := db.WithContext(ctx).
		Preload("Releases").
		Preload("Releases.Dependencies").
		Where("id = ?", id).
		First(&repo)
	return repo, result.Error
}

func SaveRepo(repo Repo) (Repo, error) {
	// Ensure every dependency exists (dedupe by name+version) so we don't create
	// duplicate dependency rows with new IDs when creating associations.
	for _, r := range repo.Releases {
		if r == nil {
			continue
		}
		for _, d := range r.Dependencies {
			if d == nil {
				continue
			}
			if err := ensureDependency(d); err != nil {
				return repo, fmt.Errorf("failed to ensure dependency %s@%s: %w", d.Name, d.Version, err)
			}
		}
	}

	var res *gorm.DB
	if repo.ID == 0 {
		res = db.Create(&repo)
	} else {
		res = db.Save(&repo)
	}
	if res.Error != nil {
		return repo, fmt.Errorf("failed to save repo: %w", res.Error)
	}
	return GetRepo(repo.ID)
}

// ensureDependency finds an existing dependency by name+version and reuses it,
// or creates it if missing. On success the passed dependency will have its ID set.
func ensureDependency(dep *Dependency) error {
    if dep == nil {
        return errors.New("nil dependency")
    }
    var existing Dependency
    result := db.Where("name = ? AND version = ?", dep.Name, dep.Version).First(&existing)
    if result.Error == nil {
        // reuse existing ID so GORM treats this as an existing record
        dep.ID = existing.ID
        return nil
    }
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            // Log the debug message instead of returning an error
            log.Printf("Dependency %s@%s not found, creating a new one.", dep.Name, dep.Version)
        } else {
            return result.Error
        }
    }
    // Not found -> create the dependency row
    if err := db.Create(dep).Error; err != nil {
        return err
    }
    return nil
}
