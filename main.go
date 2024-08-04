package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const patchUpdateDelay = 48 * time.Hour
const scrapingDelay = 30 * time.Second

func main() {
	log.Println("Application starting...")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	log.Printf("Connecting to database: %s", dbURL)

	db, err := NewDB(dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	if err := db.CreateTables(); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	log.Println("Database tables created/verified")

	// Start scraping in a separate goroutine
	log.Println("Starting scraping process in background...")
	go startScraping(db)

	// Set up REST API
	log.Println("Setting up REST API...")
	r := gin.Default()

	r.GET("/matchups/:champion/:role", func(c *gin.Context) {
		champion := c.Param("champion")
		role := c.Param("role")
		limit := c.DefaultQuery("limit", "8")
		limitInt, _ := strconv.Atoi(limit)

		log.Printf("Received request for /matchups/%s/%s with limit %d", champion, role, limitInt)

		status, err := db.GetScrapingStatus()
		if err != nil {
			log.Printf("Error getting scraping status: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Scraping status: CurrentPatch=%s, LastScrapedPatch=%s, IsUpdating=%v",
			status.CurrentPatch, status.LastScrapedPatch, status.IsUpdating)

		patch := status.LastScrapedPatch
		if patch == "" {
			patch = status.CurrentPatch
			log.Printf("LastScrapedPatch is empty, using CurrentPatch: %s", patch)
		}

		if status.IsUpdating {
			c.Header("X-Patch-Updating", "true")
		}

		log.Printf("Calling GetTopMatchups with champion=%s, role=%s, limit=%d, patch=%s",
			champion, role, limitInt, patch)

		matchups, err := db.GetTopMatchups(champion, role, limitInt, patch)
		if err != nil {
			log.Printf("Error getting top matchups: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Printf("GetTopMatchups returned %d matchups", len(matchups))

		if len(matchups) == 0 {
			log.Printf("No matchups found for %s in %s role", champion, role)
			c.JSON(404, gin.H{"error": "No matchups found", "patch": patch})
			return
		}

		log.Printf("Returning %d matchups for %s in %s role", len(matchups), champion, role)
		c.JSON(200, gin.H{"patch": patch, "matchups": matchups})
	})

	r.GET("/matchups/:champion/:role/all", func(c *gin.Context) {
		champion := c.Param("champion")
		role := c.Param("role")

		log.Printf("Received request for /matchups/%s/%s/all", champion, role)

		status, err := db.GetScrapingStatus()
		if err != nil {
			log.Printf("Error getting scraping status: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Scraping status: CurrentPatch=%s, LastScrapedPatch=%s, IsUpdating=%v",
			status.CurrentPatch, status.LastScrapedPatch, status.IsUpdating)

		patch := status.LastScrapedPatch
		if patch == "" {
			patch = status.CurrentPatch
			log.Printf("LastScrapedPatch is empty, using CurrentPatch: %s", patch)
		}

		if status.IsUpdating {
			c.Header("X-Patch-Updating", "true")
		}

		log.Printf("Calling GetAllMatchups with champion=%s, role=%s, patch=%s",
			champion, role, patch)

		matchups, err := db.GetAllMatchups(champion, role, patch)
		if err != nil {
			log.Printf("Error getting all matchups: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Printf("GetAllMatchups returned %d matchups", len(matchups))

		if len(matchups) == 0 {
			log.Printf("No matchups found for %s in %s role", champion, role)
			c.JSON(404, gin.H{"error": "No matchups found", "patch": patch})
			return
		}

		log.Printf("Returning %d matchups for %s in %s role", len(matchups), champion, role)
		c.JSON(200, gin.H{"patch": patch, "matchups": matchups})
	})

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startScraping(db *DB) {
	log.Println("Scraping process started")
	for {
		log.Println("Starting a scraping cycle")
		currentPatch, err := ScrapePatchInfo()
		if err != nil {
			log.Printf("Error scraping patch info: %v", err)
			time.Sleep(1 * time.Hour)
			continue
		}
		log.Printf("Current patch: %s", currentPatch.Version)

		status, err := db.GetScrapingStatus()
		if err != nil {
			log.Printf("Error getting scraping status: %v", err)
			time.Sleep(1 * time.Hour)
			continue
		}

		if currentPatch.Version != status.CurrentPatch || status.LastScrapedPatch == "" {
			log.Printf("New patch detected or first run: %s", currentPatch.Version)

			// Save the new patch first
			if err := db.SavePatch(currentPatch); err != nil {
				log.Printf("Error saving new patch: %v", err)
				time.Sleep(1 * time.Hour)
				continue
			}

			status.CurrentPatch = currentPatch.Version
			status.IsUpdating = true
			if err := db.UpdateScrapingStatus(status); err != nil {
				log.Printf("Error updating scraping status: %v", err)
				time.Sleep(1 * time.Hour)
				continue
			}

			// Start scraping for the new patch
			log.Println("Starting to scrape champions")
			champions, err := ScrapeChampions()
			if err != nil {
				log.Printf("Error scraping champions: %v", err)
			} else {
				log.Printf("Scraped %d champions", len(champions))
				for _, champ := range champions {
					if err := db.SaveChampion(champ); err != nil {
						log.Printf("Error saving champion %s: %v", champ.Name, err)
					}
					log.Printf("Scraping matchups for %s", champ.Name)
					matchups, err := ScrapeMatchups(champ.Name)
					if err != nil {
						log.Printf("Error scraping matchups for %s: %v", champ.Name, err)
						continue
					}
					for role, roleMatchups := range matchups {
						if err := db.SaveMatchups(champ.Name, role, roleMatchups, currentPatch.Version); err != nil {
							log.Printf("Error saving matchups for %s in %s: %v", champ.Name, role, err)
						}
					}
					log.Printf("Finished scraping matchups for %s", champ.Name)
					time.Sleep(scrapingDelay)
				}
			}

			log.Printf("Waiting for %v before serving new data", patchUpdateDelay)
			time.Sleep(patchUpdateDelay)

			status.LastScrapedPatch = currentPatch.Version
			status.IsUpdating = false
			if err := db.UpdateScrapingStatus(status); err != nil {
				log.Printf("Error updating scraping status: %v", err)
			}
			log.Println("Scraping cycle completed")
		}

		log.Println("Sleeping for 6 hours before next scraping cycle")
		time.Sleep(6 * time.Hour)
	}
}
