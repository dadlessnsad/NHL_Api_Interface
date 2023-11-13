package services

import (
	"database/sql"
	"log"
	"nhl_interface/models"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Check if the database is reachable
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Connected to the database successfully")
}

func CloseDB() {
	DB.Close()
}

// In QueryTeam, fix the database query logic
func QueryTeam(teamAbbr string) (models.Team, error) {
	var team models.Team

	// Correctly retrieve the team name
	err := DB.QueryRow("SELECT name FROM teams WHERE abbreviation = $1", teamAbbr).Scan(&team.Name)
	if err != nil {
		log.Printf("error getting team name: %v", err)
		return team, err
	}

	// Initialize the RosterResponse struct
	team.Roster = models.RosterResponse{
		Forwards:   []models.Player{},
		Defensemen: []models.Player{},
		Goalies:    []models.Player{},
	}

	// Query for the active roster
	rows, err := DB.Query(`SELECT * FROM players p
		JOIN rosters r ON p.id = r.player_id
		WHERE r.team_name = $1 AND r.active = TRUE`, team.Name)

	if err != nil {
		log.Printf("error getting roster: %v", err)
		return team, err
	}
	defer rows.Close()

	for rows.Next() {
		var player models.Player
		// Adjust Scan to correctly match the player fields
		_, err := DB.Query(`SELECT p.id, p.headshot, p.first_name, p.last_name, p.full_name, p.sweater_number, p.position_code, p.shoots_catches, p.height_in_inches, p.weight_in_pounds, p.height_in_centimeters, p.weight_in_kilograms, p.birth_date, p.birth_city, p.birth_country
		FROM players p
		JOIN rosters r ON p.id = r.player_id
		WHERE r.team_name = $1 AND r.active`, team.Name)

		if err != nil {
			log.Printf("error scanning player: %v", err)
			return team, err
		}

		// Append the player to the appropriate slice based on position
		switch player.PositionCode {
		case "C", "LW", "RW":
			team.Roster.Forwards = append(team.Roster.Forwards, player)
		case "D":
			team.Roster.Defensemen = append(team.Roster.Defensemen, player)
		case "G":
			team.Roster.Goalies = append(team.Roster.Goalies, player)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("error iterating rows: %v", err)
		return team, err
	}

	return team, nil
}

func InsertTeam(team models.Team) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO teams (name, abbreviation) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET abbreviation = $2", team.Name, team.Abbrev)
	if err != nil {
		return err
	}

	for _, player := range team.Roster.Forwards {
		playerID, err := insertPlayer(tx, player)
		if err != nil {
			return err
		}
		err = updateRoster(tx, team.Name, playerID)
		if err != nil {
			return err
		}
	}

	for _, player := range team.Roster.Defensemen {
		playerID, err := insertPlayer(tx, player)
		if err != nil {
			return err
		}
		err = updateRoster(tx, team.Name, playerID)
		if err != nil {
			return err
		}
	}

	for _, player := range team.Roster.Goalies {
		playerID, err := insertPlayer(tx, player)
		if err != nil {
			return err
		}
		err = updateRoster(tx, team.Name, playerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func insertPlayer(tx *sql.Tx, player models.Player) (int, error) {
	var playerID int
	err := tx.QueryRow(`INSERT INTO players
		(id, headshot, first_name, last_name, full_name, sweater_number, position_code, shoots_catches, height_in_inches, weight_in_pounds, height_in_centimeters, weight_in_kilograms, birth_date, birth_city, birth_country)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, to_date($13, 'YYYY-MM-DD'), $14, $15)
		ON CONFLICT (id) DO UPDATE SET
		headshot = $2, first_name = $3, last_name = $4, full_name = $5, sweater_number = $6, position_code = $7, shoots_catches = $8, height_in_inches = $9, weight_in_pounds = $10, height_in_centimeters = $11, weight_in_kilograms = $12, birth_date = to_date($13, 'YYYY-MM-DD'), birth_city = $14, birth_country = $15
		RETURNING id`, player.ID, player.Headshot, player.FirstName.Default, player.LastName.Default, player.FullName, player.SweaterNumber, player.PositionCode, player.ShootsCatches, player.HeightInInches, player.WeightInPounds, player.HeightInCentimeters, player.WeightInKilograms, player.BirthDate, player.BirthCity.Default, player.BirthCountry).Scan(&playerID)
	if err != nil {
		return 0, err
	}

	return playerID, nil
}

func updateRoster(tx *sql.Tx, teamName string, playerID int) error {
	// Insert or update the rosters table with the team and player relationship
	_, err := tx.Exec(`INSERT INTO rosters
		(team_name, player_id)
		VALUES ($1, $2)
		ON CONFLICT (team_name, player_id) DO UPDATE SET active = TRUE`, teamName, playerID)

	if err != nil {
		return err
	}
	return nil
}

/*
Table Schema

CREATE TABLE IF NOT EXISTS teams (
    name TEXT PRIMARY KEY,
    abbreviation TEXT NOT NULL
);

-- Create Players Table
CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    headshot TEXT,
    first_name TEXT,
    last_name TEXT,
    full_name TEXT,
    sweater_number INT,
    position_code TEXT,
    shoots_catches TEXT,
    height_in_inches INT,
    weight_in_pounds INT,
    height_in_centimeters INT,
    weight_in_kilograms INT,
    birth_date DATE,
    birth_city TEXT,
    birth_country TEXT
);

-- Create Rosters Table
CREATE TABLE IF NOT EXISTS rosters (
    id SERIAL PRIMARY KEY,
    team_name TEXT REFERENCES teams(name),
    player_id INT REFERENCES players(id),
    active BOOLEAN DEFAULT TRUE
);

*/
