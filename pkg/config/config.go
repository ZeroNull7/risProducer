package config

// GrpcOptions ...
type Kafka struct {
	Host      string
	Port      int
	Cert      string
	Key       string
	CA        string
	VerifySSL bool
}

type RIS struct {
	URL          string
	ClientString string
	LogUnknowns  bool
}

type Metrics struct {
	Enable bool
	Port   int
	Path   string
}

type Service struct {
	Kafka   Kafka
	Ris     RIS
	Metrics Metrics
}
