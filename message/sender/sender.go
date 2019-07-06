package sender

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func SendMessage(channelName string, message string) {
	values := url.Values{}
	token := os.Getenv("SLACK_OAUTH_TOKEN")
	values.Add("token", token)
	values.Add("channel", channelName)
	values.Add("text", message)
	values.Add("mrkdwn", "true")
	resp, err := http.PostForm("https://slack.com/api/chat.postMessage", values)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	byteString, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteString))
}
