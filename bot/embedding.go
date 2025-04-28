package bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

// EmbedText sends prompt text to OpenAI and gets 1536-dim vector
func EmbedText(text string) ([]float32, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	endpoint := "https://api.openai.com/v1/embeddings"

	body := map[string]interface{}{
		"input": text,
		"model": "text-embedding-ada-002",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if len(res.Data) == 0 {
		return nil, err
	}

	return res.Data[0].Embedding, nil
}
