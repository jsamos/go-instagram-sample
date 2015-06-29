package main

import (
	"fmt"
	"github.com/levigross/grequests"
	"frequentbio/instagram"
	"frequentbio/models"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"runtime"
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
	success := make(chan *grequests.Response)
	fail := make(chan bool)
	
	go func() {
		client := &instagram.Client{os.Getenv("IG_ACCESS_TOKEN")}
		resp, err := client.GetUser(userId) 

		if err != nil {
			fail <- true
		} else {
			success <- resp
		}
	} ()

  var value string
  var failure bool
  for i := 0; i < 1; i++ {
 		select {
    case <- fail:
    	value = ""
    	failure = true
    case <- time.After(500 * time.Millisecond):
      fmt.Println("timed out")
    	value = ""
    	failure = true
    case resp := <- success:
   		fmt.Printf("User: %v Received\n", userId)
			var udat map[string]interface{}
			resp.JSON(&udat)
			user := udat["data"].(map[string]interface{})
			value = user["bio"].(string)
			failure = false
    }
  }

  return value, failure
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

func manageWords(c chan string, model *models.FreqWords) {
	for {
		data := <-c
		model.AddWord(data)
	}
}

func main() {
	runtime.GOMAXPROCS(4)
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
	data := dat["data"].([]interface{})
	
	for _, v := range data {
		userId := extractUserId(v)
		
		if userWasProcessed(userIds, userId) != true {
			userIds = append(userIds, userId)
		}
	}

	freqWords := &models.FreqWords{map[string]int{}}
	wordChannel := make(chan string)
	go manageWords(wordChannel, freqWords)

	for _, userId := range userIds {
		words := processBioWords(userId)

		for _, w := range words {
			wordChannel <- w
		}
	}

	for word, count := range freqWords.Words {
		if count > 1 {
			fmt.Printf("%v : %v\n", word, count)
		}
	}

	fmt.Println("Time: ", time.Since(start))
}
