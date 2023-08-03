package models

type App struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

type AppUsage struct {
	AppID string `json:"app_id"`
	Count int    `json:"count"`
}

type TokenPreview struct {
	ID      string `json:"id"`
	Preview string `json:"preview"`
}
