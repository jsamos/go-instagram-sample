package main

import (
	"fmt"
	"frequentbio/instagram"
	"frequentbio/models"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func userWasProcessed(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func normalizeWords(s string) []string {
	words := []string{}
	r, _ := regexp.Compile("[^A-Za-z0-9 ]")
	clean := r.ReplaceAllString(s, " ")
	rr, _ := regexp.Compile(" +")
	moreClean := rr.ReplaceAllString(clean, " ")
	candidates := strings.Split(moreClean, " ")

	for _, v := range candidates {
		if len(v) > 3 {
			words = append(words, strings.ToLower(v))
		}
	}

	return words
}

func getUserBio(userId string) (string, bool) {
	client := &instagram.Client{os.Getenv("IG_ACCESS_TOKEN")}
	resp, err := client.GetUser(userId)

	if err != nil {
		return "", true
	} else {
		var udat map[string]interface{}
		resp.JSON(&udat)
		user := udat["data"].(map[string]interface{})
		return user["bio"].(string), false
	}
}

func processBio(userId string, freqWords *models.FreqWords) {
	bio, err := getUserBio(userId)

	if err {
		fmt.Println("Could Not Fetch user: ", userId)
	} else {
		words := normalizeWords(bio)
		for _, w := range words {
			freqWords.AddWord(w)
		}
	}
}

func main() {
	start := time.Now()
	err := godotenv.Load()

	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	client := &instagram.Client{os.Getenv("IG_ACCESS_TOKEN")}
	resp, err := client.GetTagRecent("shaving")

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	}

	var dat map[string]interface{}

	if err := resp.JSON(&dat); err != nil {
		log.Fatalln("Unable to parse JSON: ", err)
	}

	processedUsers := []string{}
	freqWords := &models.FreqWords{map[string]int{}}
	data := dat["data"].([]interface{})

	for _, v := range data {
		datum := v.(map[string]interface{})
		user := datum["user"].(map[string]interface{})
		userId := user["id"].(string)

		if userWasProcessed(processedUsers, userId) {
			fmt.Println("Skipping: ", user["username"])
		} else {
			fmt.Println("Fetching user: ", user["username"])
			processBio(userId, freqWords)
			processedUsers = append(processedUsers, userId)
		}
	}

	for word, count := range freqWords.Words {
		if count > 1 {
			fmt.Printf("%v : %v\n", word, count)
		}
	}

	fmt.Println("Time: ", time.Since(start))
}
