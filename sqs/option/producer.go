package option

import (
	"go-aws-sqs/internal/util"
	"time"
)

type Producer struct {
	Default
	DelaySeconds            time.Duration `json:"delaySeconds,omitempty"`
	MessageAttributes       any           `json:"messageAttributes,omitempty"`
	MessageSystemAttributes any           `json:"messageSystemAttributes,omitempty"`
	MessageDeduplicationId  string        `json:"messageDeduplicationId,omitempty"`
	MessageGroupId          string        `json:"messageGroupId,omitempty"`
}

func NewProducer() Producer {
	return Producer{}
}

func (o Producer) SetDelaySeconds(t time.Duration) Producer {
	o.DelaySeconds = t
	return o
}

func (o Producer) SetMessageAttributes(m any) Producer {
	o.MessageAttributes = m
	return o
}

func (o Producer) SetMessageSystemAttributes(m any) Producer {
	o.MessageSystemAttributes = m
	return o
}

func (o Producer) SetMessageDeduplicationId(s string) Producer {
	o.MessageDeduplicationId = s
	return o
}

func (o Producer) SetMessageGroupId(s string) Producer {
	o.MessageGroupId = s
	return o
}

func (o Producer) SetDebugMode(b bool) Producer {
	o.DebugMode = b
	return o
}

func (o Producer) SetOptionHttp(opt Http) Producer {
	o.OptionHttp = &opt
	return o
}

func GetProducerByParams(opts []Producer) Producer {
	var result Producer
	for _, opt := range opts {
		fillDefaultFields(opt.Default, &result.Default)
		if opt.DelaySeconds > 0 {
			result.DelaySeconds = opt.DelaySeconds
		}
		if util.IsValidType(opt.MessageAttributes) {
			result.MessageAttributes = opt.MessageAttributes
		}
		if util.IsValidType(opt.MessageSystemAttributes) {
			result.MessageSystemAttributes = opt.MessageSystemAttributes
		}
		if len(opt.MessageDeduplicationId) != 0 {
			result.MessageDeduplicationId = opt.MessageDeduplicationId
		}
		if len(opt.MessageGroupId) != 0 {
			result.MessageGroupId = opt.MessageGroupId
		}
	}
	return result
}
