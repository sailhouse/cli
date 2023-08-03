package models

type SchemaSubscriptionFilter struct {
	Path  string `yaml:"path"`
	Value string `yaml:"value"`
}

type SchemaSubscription struct {
	Slug      string                   `yaml:"slug"`
	TopicSlug string                   `yaml:"topic"`
	Type      string                   `yaml:"type"`
	Endpoint  string                   `yaml:"endpoint"`
	Filter    SchemaSubscriptionFilter `yaml:"filter"`
}

type SchemaTopic struct {
	Slug string `yaml:"slug"`
}

type Schema struct {
	Key string `yaml:"key" validate:"max=12"`

	Topics        []SchemaTopic        `yaml:"topics"`
	Subscriptions []SchemaSubscription `yaml:"subscriptions"`
}
