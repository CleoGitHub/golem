package httpclient

type HttpClientConfig struct {
	Host          string `yaml:"url"`
	Port          string `yaml:"port"`
	SA            string `yaml:"sa"`
	Password      string `yaml:"password"`
	NbRetry       int    `yaml:"nbRetry"`
	IntervalRetry int    `yaml:"intervalRetry"`
}
