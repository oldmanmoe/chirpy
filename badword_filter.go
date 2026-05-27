package main

import (
	"strings"
)

func cleanBadWords(chirp string) string{
	var cleanWords []string
	badwords := [...]string{"kerfuffle", "sharbert", "fornax"}
	
	chirpWords := strings.Split(chirp, " ")

	

	for _, word := range chirpWords {
		found := false
		for _, target := range badwords {
			if strings.ToLower(word) == target {
				cleanWords = append(cleanWords, "****")
				found = true
			}
		}
		if found == true {
			continue
		}
		cleanWords = append(cleanWords, word)
	}

	return strings.Join(cleanWords, " ")

}