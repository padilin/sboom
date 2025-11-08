package main

import (
	"log"
	"sboom/database"
)


func main() {
    database.InitDB("sboom.db")
    repo := database.Repo{
        Name: "tester",
        Releases: []*database.Release{
            {
                Version: "1.0.0",
                Dependencies: []*database.Dependency{
                    {
                        Name:    "GORM",
                        Version: "1",
                    },
                    {
                        Name:    "Git",
                        Version: "1",
                    },
                },
            },
            {
                Version: "2.0.0",
                Dependencies: []*database.Dependency{
                    {
                        Name:    "GORM",
                        Version: "2",
                    },
                    {
                        Name:    "Git",
                        Version: "1",
                    },
                },
            },
        },
    }

	savedRepo, err := database.SaveRepo(repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(savedRepo.CreatedAt)
}

