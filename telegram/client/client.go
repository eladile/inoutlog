package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

type Client interface {
	SendMessage(id string, text string) error
	GetUpdates(timeout int) ([]Update, error)
}

type client struct {
	lastUpdate int64
	basePath   url.URL
}

func NewClient(apiToken string) (Client, error) {
	c := client{basePath: url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   "bot" + apiToken,
	},
	}
	_, err := http.Get(c.joinPath("getMe"))
	if err != nil {
		return &c, err
	}
	return &c, err
}

func (c *client) SendMessage(id string, message string) error {
	query := url.Values{}
	query.Add("chat_id", id)
	query.Add("text", message)
	res := BaseResp{}
	err := c.getAndParse([]string{"sendMessage"}, query, &res)
	if err != nil {
		return err
	}
	log.Printf("Message sent successfully to %s\n", id)
	return err
}

// GetUpdates Gets the updates from the bot, since is updated accordingly.
// The timeout is for long polling in seconds so notice it is blocking!
func (c *client) GetUpdates(timeout int) ([]Update, error) {
	query := url.Values{}
	query.Add("offset", fmt.Sprintf("%d", c.lastUpdate+1))
	query.Add("timeout", fmt.Sprintf("%d", timeout))
	res := make([]Update, 0, 100)
	err := c.getAndParse([]string{"getUpdates"}, query, &res)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		c.lastUpdate = res[len(res)-1].UpdateId
	}
	return res, nil

}

func (c *client) joinPath(elements ...string) string {
	res := c.basePath
	res.Path = path.Join(res.Path, path.Join(elements...))
	return res.String()
}

func (c *client) getAndParse(pathElements []string, query url.Values, res interface{}) error {
	uri := c.basePath
	for _, element := range pathElements {
		uri.Path = path.Join(uri.Path, element)
	}
	uri.RawQuery = query.Encode()
	log.Println("Call " + uri.String())
	httpRes, err := http.Get(uri.String())
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}
	bres := BaseResp{}
	err = json.Unmarshal(data, &bres)
	if err != nil {
		return err
	}
	if !bres.Ok {
		return errors.New("failed calling" + uri.String())
	}
	return json.Unmarshal(bres.Result, res)
}
