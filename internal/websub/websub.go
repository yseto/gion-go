package websub

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yseto/gion-go/internal/config"
)

const subscribeEndpoint = "/api/websub/subscribe"

func Subscribe(ctx context.Context, topicURL string) error {
	return subscribe(ctx, topicURL, "subscribe")
}

func UnSubscribe(ctx context.Context, topicURL string) error {
	return subscribe(ctx, topicURL, "unsubscribe")
}

func subscribe(ctx context.Context, topicURL, mode string) error {
	cfg, _ := config.ReadConfig()

	var data = url.Values{}
	data.Add("hub.callback", fmt.Sprintf("https://%s%s", cfg.WebsubHostName, subscribeEndpoint))
	data.Add("hub.lease_seconds", fmt.Sprintf("%d", int(randomLeaseSeconds(10*24*time.Hour).Seconds())))
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

func randomLeaseSeconds(currentInterval time.Duration) time.Duration {
	var minInterval = currentInterval - time.Duration(12*time.Hour)
	var maxInterval = currentInterval + time.Duration(12*time.Hour)

	return time.Duration(minInterval + time.Duration(rand.Float64()*float64(maxInterval-minInterval+1))).Truncate(time.Minute)
}
