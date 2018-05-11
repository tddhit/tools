package option

type NsqProducer struct {
	Enable        bool   `yaml:"enable"`
	Registry      string `yaml:"registry"`
	RetryTimes    int    `yaml:"retryTimes"`
	RetryInterval int    `yaml:"retryInterval"`
}

type NsqConsumer struct {
	Enable   bool     `yaml:"enable"`
	Registry string   `yaml:"registry"`
	Channel  string   `yaml:"channel"`
	Topics   []string `yaml:"topics"`
}
