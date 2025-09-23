package main

import (
	"encoding/json"
	"os"
)

// Clue represents a single scavenger hunt clue. Kept minimal for now.
type Clue struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

const CLUES_FILE_PATH = "config/clues.json"

// LoadClues loads clues from the given JSON file path.
func LoadClues() ([]Clue, error) {
	// Attempt to load clues now; fine if missing for hello world.
	f, err := os.Open(CLUES_FILE_PATH)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var clues []Clue
	if err := json.NewDecoder(f).Decode(&clues); err != nil {
		return nil, err
	}
	return clues, nil
}
