package sqs

import "errors"

var ErrMessageBodyEmpty = errors.New("sqs: no message body passed")
var ErrParseBody = errors.New("sqs: message parse body failed")
var ErrInvalidBodyType = errors.New("sqs: invalid type body")
