package database

import (
	"gorm.io/gorm"
)

type Dependency struct {
	gorm.Model
	// Keep numeric ID as the primary key and enforce uniqueness on (Name,Version)
	Name     string     `gorm:"not null;uniqueIndex:idx_name_version" json:"name"`
	Version  string     `gorm:"not null;uniqueIndex:idx_name_version" json:"version"`
	Releases []*Release `gorm:"many2many:release_dependencies;"`
	// self-referential many2many for related dependencies (optional)
	Related []*Dependency `gorm:"many2many:dependency_relations" json:"related"`
}

type Repo struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Url      *string
	Path     *string
	Releases []*Release
}

type Release struct {
	gorm.Model
	Version      string
	Commit       string
	Dependencies []*Dependency `gorm:"many2many:release_dependencies;"`
	RepoID       uint
}
