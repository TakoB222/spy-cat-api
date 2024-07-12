package utils

import (
	"encoding/json"
	"net/http"
)

const CatApiEndpoint = "https://api.thecatapi.com/v1/breeds"

type Breed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ValidateCatBreed(breed string) bool {
	resp, err := http.Get(CatApiEndpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false
	}

	for _, b := range breeds {
		if b.Name == breed {
			return true
		}
	}
	return false
}
