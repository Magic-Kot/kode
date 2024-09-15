package speller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Misspell struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Row         int      `json:"row"`
	Col         int      `json:"col"`
	Len         int      `json:"len"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

type Response [][]Misspell

type Speller struct {
	url string
}

func NewSpeller(url string) *Speller {
	return &Speller{
		url: url,
	}
}

func (s *Speller) CheckTexts(ctx context.Context, postData []string) (Response, error) {
	resp, err := http.PostForm(s.url, url.Values{
		"text": postData,
	})
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var yr Response

	if err = json.Unmarshal(body, &yr); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return yr, errors.New(fmt.Sprint("Response status: ", resp.Status))
	}

	return yr, err
}
