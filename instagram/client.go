package instagram

import (
  "fmt"
	"github.com/levigross/grequests"
)

type Client struct {
  Token string
}

func (c *Client) options() *grequests.RequestOptions {
  return &grequests.RequestOptions{
    Params: map[string]string{"access_token": c.Token},
  }
}

func (c *Client) GetTagRecent(tag string) (*grequests.Response, error) {
  url := fmt.Sprintf("https://api.instagram.com/v1/tags/%v/media/recent", tag)
  options := c.options()
  return grequests.Get(url, options)
}