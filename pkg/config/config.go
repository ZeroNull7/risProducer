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
