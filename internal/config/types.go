package config

type Config struct {
	HostName string    `json:"hostName"`
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	User          string `json:"user"`
	FileNameRegEx string `json:"fileNameRegEx"`
	FileEncoding  string `json:"fileEncoding"`
	WebhookUrl    string `json:"webhookUrl"`
}
