package option

type Default struct {
	// HTTP communication customization options with AWS SQS
	HttpClient *HttpClient `json:"httpClient,omitempty"`
	// if true and print all information and error logs
	DebugMode bool `json:"debugMode,omitempty"`
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) SetDebugMode(b bool) *Default {
	d.DebugMode = b
	return d
}

func (d *Default) SetHttpClient(httpClient HttpClient) *Default {
	d.HttpClient = &httpClient
	return d
}

func GetDefaultByParams(opts []*Default) *Default {
	var result Default
	for _, opt := range opts {
		if opt != nil {
			fillDefaultFields(*opt, &result)
		}
	}
	return &result
}

func fillDefaultFields(opt Default, dest *Default) {
	if opt.DebugMode {
		dest.DebugMode = true
	}
	if opt.HttpClient != nil {
		dest.HttpClient = opt.HttpClient
	}
}
