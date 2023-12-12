package option

import (
	"time"
)

type MessageAttribute map[any]any

type OptionsProducer struct {
	baseOptions
	DelaySeconds            time.Duration    `json:"delaySeconds,omitempty"`
	MessageAttributes       MessageAttribute `json:"messageAttributes,omitempty"`
	MessageSystemAttributes MessageAttribute `json:"messageSystemAttributes,omitempty"`
	MessageDeduplicationId  string           `json:"messageDeduplicationId,omitempty"`
	MessageGroupId          string           `json:"messageGroupId,omitempty"`
}

func Producer() *OptionsProducer {
	return &OptionsProducer{}
}

func (o *OptionsProducer) SetDelaySeconds(t time.Duration) {
	o.DelaySeconds = t
}

func (o *OptionsProducer) SetMessageAttributes(m MessageAttribute) *OptionsProducer {
	o.MessageAttributes = m
	return o
}

func (o *OptionsProducer) SetMessageSystemAttributes(m MessageAttribute) *OptionsProducer {
	o.MessageSystemAttributes = m
	return o
}

func (o *OptionsProducer) SetMessageDeduplicationId(s string) *OptionsProducer {
	o.MessageDeduplicationId = s
	return o
}

func (o *OptionsProducer) SetMessageGroupId(s string) *OptionsProducer {
	o.MessageGroupId = s
	return o
}

func (o *OptionsProducer) SetDebugMode(b bool) *OptionsProducer {
	o.DebugMode = b
	return o
}

func (o *OptionsProducer) SetOptionsHttp(opt OptionsHttp) *OptionsProducer {
	o.OptionsHttp = &opt
	return o
}
