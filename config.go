package transit

type Config struct {
	URL string
	// parallel?
	Upload ConfigUpload
}

type ConfigUpload struct {
	Bucket    string // R2 Bucket
	FilePath  string
	Endpoint  string // R2 Endpoint (region?)
	AccountID string // R2 account-id
	AccessKey string
	SecretKey string
}
