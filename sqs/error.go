package sqs

import "errors"

var ErrMessageEmpty = errors.New("sqs: no message body passed")
