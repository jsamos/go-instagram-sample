package main

import (
	"os"
	"fmt"
	"log"
	"github.com/levigross/grequests"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
  
  if err != nil {
    log.Fatalln("Error loading .env file")
  }

	ro := &grequests.RequestOptions{
		Params: map[string]string{"access_token": os.Getenv("IG_ACCESS_TOKEN")},
	}

	resp, err := grequests.Get("https://api.instagram.com/v1/tags/shaving/media/recent", ro)

	if err != nil {
		log.Fatalln("Unable to make request: ", err)
	} 

	var dat map[string]interface{}

	if err := resp.JSON(&dat); err != nil {
		log.Fatalln("Unable to parse JSON: ", err)
	}

	data := dat["data"].([]interface{})

	for _, v := range data {
		datum := v.(map[string]interface{})
		user := datum["user"].(map[string]interface{})
		fmt.Println(user["username"])
	}
}
