package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)
type PostData struct {
	Kind string `json:"kind"`
	Object string `json:"object"`
	Namespace string `json:"namespace"`
	Message string `json:"message"`
	Cluster string `json:"cluster"`
}

type Webhook struct {
	URL string
	Post PostData
}
func (w Webhook) WebhookPost() error {
	d ,err := json.Marshal(w.Post)
	dbytes := bytes.NewBuffer(d)
	if err != nil {
		return err
	}
	resp, err := http.Post(w.URL,"application/json", dbytes,)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Invalid Response Code %s", fmt.Sprint(resp.StatusCode))
	}

	defer resp.Body.Close()
	return nil
}