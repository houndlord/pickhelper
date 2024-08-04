package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS patches (
			version TEXT PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS champions (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			avatar_url TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS matchups (
			id SERIAL PRIMARY KEY,
			champion_id INT REFERENCES champions(id),
			opponent_id INT REFERENCES champions(id),
			role TEXT NOT NULL,
			win_rate FLOAT NOT NULL,
			sample_size INT NOT NULL,
			patch TEXT REFERENCES patches(version),
			UNIQUE(champion_id, opponent_id, role, patch)
		)`,
		`CREATE TABLE IF NOT EXISTS scraping_status (
			id INT PRIMARY KEY DEFAULT 1,
			current_patch TEXT REFERENCES patches(version),
			last_scraped_patch TEXT REFERENCES patches(version),
			is_updating BOOLEAN NOT NULL DEFAULT false,
			CHECK (id = 1)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}

func (db *DB) SavePatch(patch PatchInfo) error {
	_, err := db.Exec(`
		INSERT INTO patches (version)
		VALUES ($1)
		ON CONFLICT (version) DO NOTHING
	`, patch.Version)
	return err
}

func (db *DB) UpdateScrapingStatus(status ScrapingStatus) error {
	if status.CurrentPatch == "" {
		return fmt.Errorf("current_patch cannot be empty")
	}

	_, err := db.Exec(`
        INSERT INTO scraping_status (id, current_patch, last_scraped_patch, is_updating)
        VALUES (1, $1, $2, $3)
        ON CONFLICT (id) DO UPDATE 
        SET current_patch = $1, last_scraped_patch = $2, is_updating = $3
    `, status.CurrentPatch, status.LastScrapedPatch, status.IsUpdating)
	return err
}

func (db *DB) GetCurrentPatch() (PatchInfo, error) {
	var patch PatchInfo
	err := db.QueryRow(`
		SELECT version
		FROM patches
		ORDER BY version DESC
		LIMIT 1
	`).Scan(&patch.Version)
	return patch, err
}

func (db *DB) SaveChampion(champ Champion) error {
	_, err := db.Exec(`
		INSERT INTO champions (name, avatar_url)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET avatar_url = $2
	`, champ.Name, champ.AvatarURL)
	return err
}

func (db *DB) SaveMatchups(champName string, role string, matchups []Matchup, patch string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, m := range matchups {
		winRate, err := strconv.ParseFloat(m.WinRate, 64)
		if err != nil {
			log.Printf("Error parsing win rate for %s vs %s: %v", champName, m.Champion, err)
			continue
		}

		sampleSize, err := strconv.Atoi(strings.ReplaceAll(m.SampleSize, ",", ""))
		if err != nil {
			log.Printf("Error parsing sample size for %s vs %s: %v", champName, m.Champion, err)
			continue
		}

		_, err = tx.Exec(`
			WITH champ AS (
				SELECT id FROM champions WHERE name = $1
			), opp AS (
				SELECT id FROM champions WHERE name = $2
			)
			INSERT INTO matchups (champion_id, opponent_id, role, win_rate, sample_size, patch)
			SELECT champ.id, opp.id, $3, $4, $5, $6
			FROM champ, opp
			ON CONFLICT (champion_id, opponent_id, role, patch) 
			DO UPDATE SET win_rate = $4, sample_size = $5
		`, champName, m.Champion, role, winRate, sampleSize, patch)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db *DB) GetScrapingStatus() (ScrapingStatus, error) {
	var status ScrapingStatus
	err := db.QueryRow(`
        SELECT current_patch, last_scraped_patch, is_updating
        FROM scraping_status
        WHERE id = 1
    `).Scan(&status.CurrentPatch, &status.LastScrapedPatch, &status.IsUpdating)

	if err == sql.ErrNoRows {
		// If no row exists, return an empty status without an error
		return ScrapingStatus{}, nil
	}

	return status, err
}

func (db *DB) InitializeScrapingStatus() (ScrapingStatus, error) {
	status := ScrapingStatus{
		CurrentPatch:     "",
		LastScrapedPatch: "",
		IsUpdating:       false,
	}

	_, err := db.Exec(`
		INSERT INTO scraping_status (id, current_patch, last_scraped_patch, is_updating)
		VALUES (1, $1, $2, $3)
		ON CONFLICT (id) DO UPDATE 
		SET current_patch = $1, last_scraped_patch = $2, is_updating = $3
	`, status.CurrentPatch, status.LastScrapedPatch, status.IsUpdating)

	return status, err
}

func (db *DB) GetTopMatchups(champName string, role string, limit int, patch string) ([]Matchup, error) {
	rows, err := db.Query(`
		SELECT c.name, m.win_rate, m.sample_size
		FROM matchups m
		JOIN champions c ON m.opponent_id = c.id
		JOIN champions champ ON m.champion_id = champ.id
		WHERE LOWER(champ.name) = LOWER($1) AND LOWER(m.role) = LOWER($2) AND m.patch = $3
		ORDER BY m.win_rate DESC
		LIMIT $4
	`, champName, role, patch, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchups []Matchup
	for rows.Next() {
		var m Matchup
		var winRate float64
		var sampleSize int
		if err := rows.Scan(&m.Champion, &winRate, &sampleSize); err != nil {
			return nil, err
		}
		m.WinRate = fmt.Sprintf("%.2f", winRate)
		m.SampleSize = strconv.Itoa(sampleSize)
		matchups = append(matchups, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matchups, nil
}

func (db *DB) GetAllMatchups(champName string, role string, patch string) ([]Matchup, error) {
	log.Printf("GetAllMatchups called with champName: %s, role: %s, patch: %s", champName, role, patch)

	query := `
		SELECT c.name, m.win_rate, m.sample_size
		FROM matchups m
		JOIN champions c ON m.opponent_id = c.id
		JOIN champions champ ON m.champion_id = champ.id
		WHERE LOWER(champ.name) = LOWER($1) AND LOWER(m.role) = LOWER($2) AND m.patch = $3
		ORDER BY m.win_rate DESC
	`

	rows, err := db.Query(query, champName, role, patch)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var matchups []Matchup
	for rows.Next() {
		var m Matchup
		var winRate float64
		var sampleSize int
		if err := rows.Scan(&m.Champion, &winRate, &sampleSize); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		m.WinRate = fmt.Sprintf("%.2f", winRate)
		m.SampleSize = strconv.Itoa(sampleSize)
		matchups = append(matchups, m)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d matchups", len(matchups))

	return matchups, nil
}
