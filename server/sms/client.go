package sms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client - a SMS Client
type Client struct {
	accountSid  string
	authToken   string
	messageFrom string
}

// New creates a new SMS Client
func New(sid string, token string) *Client {
	return &Client{
		accountSid:  sid,
		authToken:   token,
		messageFrom: "Genesis",
	}
}

// SendMessage sends an SMS message to mobileNumber
func (c *Client) SendMessage(mobileNumber string, content string) error {
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", c.accountSid)

	msgData := url.Values{}
	msgData.Set("To", mobileNumber)
	msgData.Set("From", c.messageFrom)
	msgData.Set("Body", content)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(c.accountSid, c.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Twilio send to %s: Success\n", mobileNumber)
	} else {
		fmt.Printf("Twilio send to %s: Fail error: %s \n", mobileNumber, resp.Status)
		return fmt.Errorf("failed to send SMS with Twilio error: %s", resp.Status)
	}

	return nil
}

// SendToken sends a verification token via SMS
func (c *Client) SendToken(mobileNumber string, token string) error {
	content := fmt.Sprintf("%s is your Genesis verification code.", token)
	return c.SendMessage(mobileNumber, content)
}
