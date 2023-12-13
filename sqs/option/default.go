package option

type Default struct {
	OptionHttp *Http `json:"optionsHttp,omitempty"`
	DebugMode  bool  `json:"debugMode,omitempty"`
}

func NewDefault() Default {
	return Default{}
}

func (o Default) SetDebugMode(b bool) Default {
	o.DebugMode = b
	return o
}

func (o Default) SetOptionHttp(opt Http) Default {
	o.OptionHttp = &opt
	return o
}

func GetDefaultByParams(opts []Default) Default {
	var result Default
	for _, opt := range opts {
		fillDefaultFields(opt, &result)
	}
	return result
}

func fillDefaultFields(opt Default, dest *Default) {
	if opt.DebugMode {
		dest.DebugMode = true
	}
	if opt.OptionHttp != nil {
		dest.OptionHttp = opt.OptionHttp
	}
}
