package api

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/sailhouse/sailhouse/models"
	"github.com/spf13/viper"
)

type SailhouseClient struct {
	token string
	team  string
}

func NewSailhouseClient(token string) *SailhouseClient {
	team := viper.GetString("team")
	return &SailhouseClient{token, team}
}

func (c *SailhouseClient) req() *requests.Builder {
	return requests.
		URL("https://api.sailhouse.dev").
		Header("Authorization", c.token)
}

func (c *SailhouseClient) GetTeams(ctx context.Context) ([]models.Team, error) {
	teams := []models.Team{}

	err := c.req().
		Path("teams").
		ToJSON(&teams).
		Fetch(ctx)

	return teams, err
}

func (c *SailhouseClient) CreateApp(ctx context.Context, name string) error {
	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps", c.team)

	return requests.
		URL(url).
		Header("Authorization", c.token).
		BodyJSON(map[string]string{"name": name, "slug": name}).
		Fetch(ctx)
}

func (c *SailhouseClient) GetApps(ctx context.Context) ([]models.App, error) {
	apps := []models.App{}

	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps", c.team)

	err := requests.
		URL(url).
		Header("Authorization", c.token).
		ToJSON(&apps).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (c *SailhouseClient) GetTopics(ctx context.Context, appID string) ([]models.Topic, error) {
	topics := []models.Topic{}

	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/topics", c.team, appID)

	err := requests.
		URL(url).
		Header("Authorization", c.token).
		ToJSON(&topics).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return topics, nil
}

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (c *SailhouseClient) CreateToken(ctx context.Context, appID string) (string, error) {
	var resp CreateTokenResponse

	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/tokens", c.team, appID)

	err := requests.
		URL(url).
		Header("Authorization", c.token).
		Method("POST").
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (c *SailhouseClient) GetTokens(ctx context.Context, appID string) ([]models.TokenPreview, error) {
	var resp []models.TokenPreview

	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/tokens", c.team, appID)

	err := requests.
		URL(url).
		Header("Authorization", c.token).
		ToJSON(&resp).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *SailhouseClient) CreateTopic(ctx context.Context, appID, slug string) error {
	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/topics", c.team, appID)

	return requests.
		URL(url).
		Header("Authorization", c.token).
		BodyJSON(map[string]any{"slug": slug, "subscriptions": []string{}}).
		Fetch(ctx)
}

func (c *SailhouseClient) DeleteTopic(ctx context.Context, appID, slug string) error {
	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/topics/%s", c.team, appID, slug)

	return requests.
		URL(url).
		Header("Authorization", c.token).
		Method("DELETE").
		Fetch(ctx)
}

func (c *SailhouseClient) CreateTopicWithKey(ctx context.Context, appID, slug, key string) error {
	return requests.
		URL("https://api.sailhouse.dev/topics").
		Header("Authorization", c.token).
		BodyJSON(map[string]any{"slug": slug, "subscriptions": []string{}, "schema_key": key}).
		Fetch(ctx)
}

func (c *SailhouseClient) GetSubscriptions(ctx context.Context, appID, topicSlug string) ([]models.Subscription, error) {
	subscriptions := []models.Subscription{}

	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/topics/%s/subscriptions", c.team, appID, topicSlug)

	err := requests.
		URL(url).
		Header("Authorization", c.token).
		ToJSON(&subscriptions).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (c *SailhouseClient) GetSubscription(ctx context.Context, appID, topicSlug, subscriptionSlug string) (*models.Subscription, error) {
	subscription := models.Subscription{}
	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/topics/%s/subscriptions/%s", c.team, appID, topicSlug, subscriptionSlug)

	err := requests.
		URL(url).
		Header("Authorization", c.token).
		ToJSON(&subscription).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

type CreateSubscription struct {
	Slug        string
	TopicSlug   string
	Type        string
	Endpoint    string
	SchemaKey   string
	FilterPath  string
	FilterValue string
}

func (c *SailhouseClient) CreateSubscription(ctx context.Context, appID string, newSub CreateSubscription) (models.Subscription, error) {
	body := map[string]string{
		"slug": newSub.Slug,
		"type": newSub.Type,
	}

	if newSub.Endpoint != "" {
		body["endpoint"] = newSub.Endpoint
	}

	if newSub.SchemaKey != "" {
		body["schema_key"] = newSub.SchemaKey
	}

	if newSub.FilterPath != "" {
		body["filter_path"] = newSub.FilterPath
		fmt.Printf("filter path: %s\n", newSub.FilterPath)
	}

	if newSub.FilterValue != "" {
		body["filter_value"] = newSub.FilterValue
		fmt.Printf("filter value: %s\n", newSub.FilterValue)
	}

	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/topics/%s/subscriptions", c.team, appID, newSub.TopicSlug)

	var sub models.Subscription
	err := requests.
		URL(url).
		Header("Authorization", c.token).
		BodyJSON(body).
		ToJSON(&sub).
		Fetch(ctx)

	return sub, err
}

func (c *SailhouseClient) DeleteSubscription(ctx context.Context, appID, topicSlug, subscriptionSlug string) error {
	url := fmt.Sprintf("https://api.sailhouse.dev/teams/%s/apps/%s/topics/%s/subscriptions/%s", c.team, appID, topicSlug, subscriptionSlug)

	return requests.
		URL(url).
		Header("Authorization", c.token).
		Method("DELETE").
		Fetch(ctx)
}

type GetSchema struct {
	Topics        []models.Topic        `json:"topics"`
	Subscriptions []models.Subscription `json:"subscriptions"`
}
