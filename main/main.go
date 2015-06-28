package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"frequentbio/instagram"
	"log"
	"os"
)

func main() {

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

	data := dat["data"].([]interface{})

	for _, v := range data {
		datum := v.(map[string]interface{})
		user := datum["user"].(map[string]interface{})
		fmt.Println(user["username"])
	}
}
