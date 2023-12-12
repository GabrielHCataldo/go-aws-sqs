package option

type baseOptions struct {
	OptionsHttp *OptionsHttp `json:"optionsHttp,omitempty"`
	DebugMode   bool         `json:"debugMode,omitempty"`
}

type OptionsDefault struct {
	baseOptions
}

func Default() *OptionsDefault {
	return &OptionsDefault{}
}

func (o *OptionsDefault) SetDebugMode(b bool) *OptionsDefault {
	o.DebugMode = b
	return o
}

func (o *OptionsDefault) SetOptionsHttp(opt OptionsHttp) *OptionsDefault {
	o.OptionsHttp = &opt
	return o
}
