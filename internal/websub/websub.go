package websub

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/yseto/gion-go/internal/config"
)

const subscribeEndpoint = "/api/websub/subscribe"

func Subscribe(ctx context.Context, topicURL string) error {
	if topicURLs.insert(topicURL) != true {
		return fmt.Errorf("failed topicURLs.insert")
	}

	return subscribe(ctx, topicURL, "subscribe")
}

func UnSubscribe(ctx context.Context, topicURL string) error {
	if topicURLs.insert(topicURL) != true {
		return fmt.Errorf("failed topicURLs.insert")
	}

	return subscribe(ctx, topicURL, "unsubscribe")
}

func subscribe(ctx context.Context, topicURL, mode string) error {
	cfg, _ := config.ReadConfig()

	var data = url.Values{}
	data.Add("hub.callback", fmt.Sprintf("https://%s%s", cfg.WebsubHostName, subscribeEndpoint))
	data.Add("hub.lease_seconds", "864000")
	data.Add("hub.mode", mode)
	data.Add("hub.topic", topicURL)
	data.Add("hub.verify", "sync")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://pubsubhubbub.appspot.com/subscribe", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 202 || resp.StatusCode == 204 {
		return nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("subscribe: statusCode:%d, response:%s", resp.StatusCode, string(b))
}
