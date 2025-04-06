package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

var stateLock sync.Mutex

type States struct {
	States map[string]*State `json:"states"`
}

func (states *States) Find(email string) *State {
	state, found := states.States[email]
	if !found {
		state = &State{}
		states.States[email] = state
	}
	return state
}

type State struct {
	BlacklistHash string    `json:"blacklist_hash"`
	SeqNumber     uint32    `json:"seq_number"`
	Date          time.Time `json:"date"`
}

// LoadState loads the state (last processed sequence numbers) from the JSON file
func LoadState() (*States, error) {
	states := &States{States: map[string]*State{}}
	stateLock.Lock()
	defer stateLock.Unlock()

	file, err := os.Open(stateJsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			return states, nil // Return an empty map if the file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(states); err != nil {
		return nil, err
	}
	return states, nil
}

// SaveState saves the state (last processed sequence numbers) to the JSON file
func SaveState(states *States) error {
	stateLock.Lock()
	defer stateLock.Unlock()

	file, err := os.Create(stateJsonPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(states)
}
