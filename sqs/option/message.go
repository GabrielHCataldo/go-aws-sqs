package option

type ListMessageMoveTasks struct {
	Default
	// The maximum number of results to include in the response. The default is 1,
	// which provides the most recent message movement task. The upper limit is 10.
	MaxResults int32
}

func NewListMessageMoveTasks() ListMessageMoveTasks {
	return ListMessageMoveTasks{}
}

func (o ListMessageMoveTasks) SetDebugMode(b bool) ListMessageMoveTasks {
	o.DebugMode = b
	return o
}

func (o ListMessageMoveTasks) SetHttpClient(opt HttpClient) ListMessageMoveTasks {
	o.HttpClient = &opt
	return o
}

func (o ListMessageMoveTasks) SetMaxResults(i int32) ListMessageMoveTasks {
	o.MaxResults = i
	return o
}

func GetListMessageMoveTaskByParams(opts []ListMessageMoveTasks) ListMessageMoveTasks {
	var result ListMessageMoveTasks
	for _, opt := range opts {
		fillDefaultFields(opt.Default, &result.Default)
		if opt.MaxResults > 0 {
			result.MaxResults = opt.MaxResults
		}
	}
	if result.MaxResults == 0 {
		result.MaxResults = 1
	}
	return result
}
