package instagram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"backend-service/internal/token"
)

type Service struct {
	Client     *http.Client
	IgUserID   string
	TokenStore *token.TokenRuntime
}

func (s *Service) FetchMedia() ([]Media, error) {
	return s.FetchMediaWithLimit(0) // 0 means fetch all
}

func (s *Service) FetchMediaWithLimit(limit int) ([]Media, error) {
	token := s.TokenStore.Get()
	var allMedia []Media

	url := fmt.Sprintf(
		os.Getenv("FB_API_BASE_URL")+"/%s/media?fields=id,caption,media_type,media_url,permalink,timestamp&access_token=%s",
		s.IgUserID, token,
	)

	if limit > 0 {
		url = fmt.Sprintf("%s&limit=%d", url, limit)
	}

	for url != "" {
		res, err := s.Client.Get(url)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("[GRAPH API] through error with status %s", res.Status)
			return nil, fmt.Errorf("[GRAPH API] through error with status %s", res.Status)
		}

		var result struct {
			Data   []Media `json:"data"`
			Paging struct {
				Next string `json:"next"`
			} `json:"paging"`
		}

		err = json.NewDecoder(res.Body).Decode(&result)
		res.Body.Close()

		if err != nil {
			return nil, err
		}

		allMedia = append(allMedia, result.Data...)

		if limit > 0 && len(allMedia) >= limit {
			return allMedia[:limit], nil
		}

		url = result.Paging.Next

		if url == "" {
			break
		}
	}

	return allMedia, nil
}

func RefreshAccessToken(client *http.Client, current string) (token.Token, error) {
	url := fmt.Sprintf(
		os.Getenv("FB_API_BASE_URL")+"/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), current,
	)

	log.Printf("Refreshing token with URL: %s", url)

	res, err := client.Get(url)

	if err != nil {
		return token.Token{}, err
	}

	if res.StatusCode != http.StatusOK {
		return token.Token{}, fmt.Errorf("[GRAPH API] through error with status %s", res.Status)
	}

	defer res.Body.Close() // closing body to prevent memory leaks

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return token.Token{}, err
	}

	return token.Token{
		AccessToken: result.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
	}, nil
}
