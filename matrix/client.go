package matrix

import (
	"github.com/matrix-org/gomatrix"

	"github.com/sirupsen/logrus"
)

type Client struct {
	mxClient *gomatrix.Client
	UserID   string
}

func NewClient(csUrl string, accessToken string) (*Client, error) {
	client := &Client{}
	mxClient, err := gomatrix.NewClient(csUrl, "", accessToken)
	if err != nil {
		return nil, err
	}

	client.mxClient = mxClient

	logrus.Info("Querying for user ID")
	resp := &WhoAmIResponse{}
	url := mxClient.BuildURL("/account/whoami")
	_, err = mxClient.MakeRequest("GET", url, nil, resp)
	if err != nil {
		return nil, err
	}
	client.UserID = resp.UserId
	mxClient.UserID = resp.UserId

	return client, nil
}

func (c *Client) JoinRoom(roomIdOrAlias string) (error) {
	_, err := c.mxClient.JoinRoom(roomIdOrAlias, "", nil)
	return err
}

func (c *Client) SendNotice(roomId string, status string, message string) (error) {
	_, err := c.mxClient.SendMessageEvent(roomId, "m.room.message", struct {
		Body    string `json:"body"`
		Msgtype string `json:"msgtype"`
		Status  string `json:"status"`
	}{
		Body:    status + " | " + message,
		Msgtype: "m.notice",
		Status:  status,
	})
	return err
}

func (c *Client) SendMessage(roomId string, message string) (error) {
	_, err := c.mxClient.SendText(roomId, message)
	return err
}
