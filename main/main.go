package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/levigross/grequests"
	"log"
	"os"
)

type Client struct {
	Token string
}

func (c *Client) GetTagRecent(tag string) (*grequests.Response, error) {
	ro := &grequests.RequestOptions{
		Params: map[string]string{"access_token": c.Token},
	}
	url := fmt.Sprintf("https://api.instagram.com/v1/tags/%v/media/recent", tag)
	return grequests.Get(url, ro)
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	client := &Client{os.Getenv("IG_ACCESS_TOKEN")}
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
