package config

// GrpcOptions ...
type Kafka struct {
	Host string
	Port int
	Cert string
	Key  string
	CA   string
}

type RIS struct {
	URL          string
	ClientString string
	LogUnknowns  bool
}

type Service struct {
	Kafka Kafka
	Ris   RIS
}
