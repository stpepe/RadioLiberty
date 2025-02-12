package queue_service_client

import (
	"RadioLiberty/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type QueueServiceClient struct {
	url *url.URL
}

func NewQueueServiceClient(host string, port string) *QueueServiceClient {
	return &QueueServiceClient{
		url: &url.URL{
			Scheme: "http",
			Host:   host + ":" + port,
		},
	}
}

func (qs *QueueServiceClient) GetNextAudioInfo(ctx context.Context) (*models.AudioInfo, error) {
	const errorMsg = "from queueService GetNextAudioInfo error: %w"
	req, err := http.NewRequest("GET", qs.url.String()+"/next", nil)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		}
		slog.Error("from queueService GetNextAudioInfo error", "response_status_code", resp.StatusCode)
		return nil, fmt.Errorf(errorMsg, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	var audioInfo models.AudioInfo
	err = json.Unmarshal(body, &audioInfo)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	return &audioInfo, nil
}
