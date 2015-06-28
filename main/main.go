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

func extractUserId(data interface{}) string {
	datum := data.(map[string]interface{})
	user := datum["user"].(map[string]interface{})
	return user["id"].(string)
}

func getUserBio(userId string) (string, bool) {
	fmt.Println("Fetching user: ", userId)
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

func processBioWords(userId string) []string {
	bio, err := getUserBio(userId)

	if err {
		fmt.Println("Could Not Fetch user: ", userId)
		return []string{}
	} else {
		return normalizeWords(bio)
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

	userIds := []string{}
	freqWords := &models.FreqWords{map[string]int{}}
	data := dat["data"].([]interface{})
	
	for _, v := range data {
		userId := extractUserId(v)
		
		if userWasProcessed(userIds, userId) != true {
			userIds = append(userIds, userId)
		}
	}

	for _, userId := range userIds {
		words := processBioWords(userId)

		for _, w := range words {
			freqWords.AddWord(w)
		}
	}

	time.Sleep(time.Second * 10)

	for word, count := range freqWords.Words {
		if count > 1 {
			fmt.Printf("%v : %v\n", word, count)
		}
	}

	fmt.Println("Time: ", time.Since(start))
}
