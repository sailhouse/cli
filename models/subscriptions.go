package models

type Subscription struct {
	ID          string `json:"id"`
	TopicID     string `json:"topic_id"`
	Slug        string `json:"slug"`
	Type        string `json:"type"`
	FilterPath  string `json:"filter_path"`
	FilterValue string `json:"filter_value"`
	Endpoint    string `json:"endpoint"`
}
