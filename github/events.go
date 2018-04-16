package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EventPayload discribes the different event payload in a
// generic way.
type EventPayload map[string]interface{}

// EventRepo describes the event repository.
type EventRepo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// EventUser describes the user behind organization and
// event triggers.
type EventUser struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	GravatarID string `json:"gravatar_id"`
	AvatarURL  string `json:"avatar_url"`
	URL        string `json:"url"`
}

// Event describes one event returned by GitHub.
type Event struct {
	Type      string       `json:"type"`
	Public    bool         `json:"public"`
	Payload   EventPayload `json:"payload"`
	Repo      EventRepo    `json:"repo"`
	Actor     EventUser    `json:"actor"`
	Org       EventUser    `json:"org"`
	CreatedAt time.Time    `json:"created_at"`
	ID        string       `json:"id"`
}

// Events contains a number of events.
type Events []Event

// RepoEventor retrieves the events of one GitHub repository.
type RepoEventor struct {
	owner string
	repo  string
}

// NewRepoEventor creates the retriever for the events of one
// repository.
func NewRepoEventor(o, r string) *RepoEventor {
	return &RepoEventor{
		owner: o,
		repo:  r,
	}
}

// Get retrieves the newest events for the configured repository.
func (e *RepoEventor) Get() (Events, error) {
	// Prepare request.
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/events", e.owner, e.repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Perform request.
	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code %d", resp.StatusCode)
	}
	// Unmarshall events.
	var buf []byte
	var events Events
	_, err = resp.Body.Read(buf)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
