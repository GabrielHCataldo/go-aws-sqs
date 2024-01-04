AWS SQS Template
=================
<!--suppress ALL -->
<img align="right" src="gopher-sqs.png" alt="">

[![Project status](https://img.shields.io/badge/version-v1.0.1-vividgreen.svg)](https://github.com/GabrielHCataldo/go-aws-sqs-template/releases/tag/v1.0.1)
[![Go Report Card](https://goreportcard.com/badge/github.com/GabrielHCataldo/go-aws-sqs-template)](https://goreportcard.com/report/github.com/GabrielHCataldo/go-aws-sqs-template)
[![Coverage Status](https://coveralls.io/repos/GabrielHCataldo/go-aws-sqs-template/badge.svg?branch=main&service=github)](https://coveralls.io/github/GabrielHCataldo/go-aws-sqs-template?branch=main)
[![Open Source Helpers](https://www.codetriage.com/gabrielhcataldo/go-aws-sqs-template/badges/users.svg)](https://www.codetriage.com/gabrielhcataldo/go-aws-sqs-template)
[![AWS SDK](https://badgen.net/badge/AWS/SDK/orange)](https://github.com/aws/aws-sdk-go-v2)
[![GoDoc](https://godoc.org/github/GabrielHCataldo/go-aws-sqs-template?status.svg)](https://pkg.go.dev/github.com/GabrielHCataldo/go-aws-sqs-template/sqs)
![License](https://img.shields.io/dub/l/vibe-d.svg)

[//]: # ([![build workflow]&#40;https://github.com/GabrielHCataldo/go-aws-sqs-template/actions/workflows/go.yml/badge.svg&#41;]&#40;https://github.com/GabrielHCataldo/go-aws-sqs-template/actions&#41;)

[//]: # ([![Source graph]&#40;https://sourcegraph.com/github.com/go-aws-sqs-template/sqs/-/badge.svg&#41;]&#40;https://sourcegraph.com/github.com/go-aws-sqs-template/sqs?badge&#41;)

[//]: # ([![TODOs]&#40;https://badgen.net/https/api.tickgit.com/badgen/github.com/GabrielHCataldo/go-aws-sqs-template/sqs&#41;]&#40;https://www.tickgit.com/browse?repo=github.com/GabrielHCataldo/go-aws-sqs-template&#41;)

The go-aws-sqs-template project came to facilitate the use of aws sqs in your go project with incredible flexibility in sending 
messages to any type of body, and a fantastic and simple to use consumer. Below we list some implemented features:

- Simplicity in message production, with auto conversion of body and message attributes
- Powerful customizable consumer, with auto conversion of body and message attributes.
- Automatic removal of the successfully consumed message (Optional).
- More simplistic function calls.
- Log visibility for the consumer.
- Asynchronous message production.
- Asynchronous message consumption.
- Don't worry about unwanted type conversions anymore.

Installation
------------

Use go get.

	go get github.com/GabrielHCataldo/go-aws-sqs-template

Then import the go-aws-sqs-template package into your own code.

```go
import "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
```

Usability and documentation
------------
**IMPORTANT**: Always check the documentation in the structures and functions fields.
For more details on the examples, visit [All examples link](https://github/GabrielHCataldo/go-aws-sqs-template/blob/main/_example)

### Producer
To produce a message, it's simple, in the example below we will send a normal body text:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
)

func main() {
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    body := "body test"
    message, err := sqs.SendMessage(ctx, os.Getenv("SQS_QUEUE_TEST_URL"), body)
    if err != nil {
        logger.Error("error send message:", err)
    } else {
        logger.Info("message sent successfully:", message)
    }
}
```

Output:

    [INFO 2023/12/15 07:22:34] main.go:33: message sent successfully: {"MD5OfMessageAttributes":null,"MD5OfMessageBody":"2e9cc74f6f6aca12eaa2d252df457910","MD5OfMessageSystemAttributes":null,"MessageId":"48276fb7-022d-4114-8282-b407d8fd4dd3","SequenceNumber":null,"ResultMetadata":{}}

You can also pass any type of value, below we will pass the body as a structure:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
    "time"
)

type test struct {
    Name      string    `json:"name,omitempty"`
    BirthDate time.Time `json:"birthDate,omitempty"`
    Emails    []string  `json:"emails,omitempty"`
    Bank      bank      `json:"bank,omitempty"`
    Map       map[string]any
}

type bank struct {
    Account string  `json:"account,omitempty"`
    Digits  string  `json:"digits,omitempty"`
    Balance float64 `json:"balance,omitempty"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    body := initTestStruct()
    message, err := sqs.SendMessage(ctx, os.Getenv("SQS_QUEUE_TEST_URL"), body)
    if err != nil {
        logger.Error("error send message:", err)
    } else {
        logger.Info("message sent successfully:", message)
    }
}
```

Output:
 
    [INFO 2023/12/15 07:30:42] main.go:45: message sent successfully: {"MD5OfMessageAttributes":null,"MD5OfMessageBody":"52e95a6a12e47e7ef6f63ff0ccfb77b6","MD5OfMessageSystemAttributes":null,"MessageId":"b43e19ca-48d9-47e8-8457-183242b86e1e","SequenceNumber":null,"ResultMetadata":{}}

As an optional parameter, we have **DelaySeconds**, **MessageAttributes**, **MessageSystemAttributes**
among others, see below:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
    "time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    body := initTestStruct()
    opt := option.NewProducer().
        // HTTP communication customization options with AWS SQS
        SetHttpClient(option.HttpClient{}).
        // print logs (default: false)
        SetDebugMode(true).
        // delay to delay the availability of message processing (default: 0)
        SetDelaySeconds(5 * time.Second).
        // Message attributes, must be of type Map or Struct, other types are not acceptable.
        SetMessageAttributes(initTestMap()).
        // The message system attribute to send
        SetMessageSystemAttributes(option.MessageSystemAttributes{}).
        // This parameter applies only to FIFO (first-in-first-out) queues. The token used for deduplication of sent messages.
        SetMessageDeduplicationId("").
        // This parameter applies only to FIFO (first-in-first-out) queues. The tag that specifies that a message belongs to a specific message group.
        SetMessageGroupId("")
	
    message, err := sqs.SendMessage(ctx, os.Getenv("SQS_QUEUE_TEST_URL"), body, opt)
    if err != nil {
        logger.Error("error send message:", err)
    } else {
        logger.Info("message sent successfully:", message)
    }
}
```

Output:

    [INFO 2023/12/15 08:00:54] producer.go:40: getting client sqs..
    [INFO 2023/12/15 08:00:54] producer.go:42: preparing message input..
    [INFO 2023/12/15 08:00:54] producer.go:48: sending message..
    [INFO 2023/12/15 08:00:54] producer.go:53: message sent successfully: {"MD5OfMessageAttributes":"2a0b53405b3678d246c0a76b321e1cea","MD5OfMessageBody":"19554d258904ecb0a27b3cf238c2ecf2","MD5OfMessageSystemAttributes":null,"MessageId":"e08368e5-08b5-456b-a4f3-b6fb7862b893","SequenceNumber":null,"ResultMetadata":{}}
    [INFO 2023/12/15 08:00:54] main.go:85: message sent successfully: {"MD5OfMessageAttributes":"2a0b53405b3678d246c0a76b321e1cea","MD5OfMessageBody":"19554d258904ecb0a27b3cf238c2ecf2","MD5OfMessageSystemAttributes":null,"MessageId":"e08368e5-08b5-456b-a4f3-b6fb7862b893","SequenceNumber":null,"ResultMetadata":{}}

We can use all these examples mentioned above asynchronously, calling the function
**SendMessageAsync**, see:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
)

func main() {
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()
    body := "body test"
    sqs.SendMessageAsync(ctx, os.Getenv("SQS_QUEUE_TEST_URL"), body)
}
```

For more producer examples visit: [All examples produce](https://github/GabrielHCataldo/go-aws-sqs-template/blob/main/_example/producer/main.go)

### Consumer

To consume messages from the queue, it is also very simple, you don't need to worry about writing loop
lines, conversion lines, log lines, etc., see a simple example below:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
)

func main() {
    sqs.SimpleReceiveMessage(os.Getenv("SQS_QUEUE_TEST_STRING_URL"), handler)
}

func handler(ctx *sqs.SimpleContext[string]) error {
    logger.Debug("ctx simple to process message:", ctx)
    return nil
}
```

Output:

    [DEBUG 2023/12/15 08:27:49] main.go:43: ctx simple to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test-string","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":"body test","Id":"d2b78dfe-8b68-4b7d-af59-15826c52515a","MD5OfBody":"2e9cc74f6f6aca12eaa2d252df457910","MD5OfMessageAttributes":null,"MessageAttributes":null,"ReceiptHandle":"AQEBG2Ek7PrGMDpnvCbTh3pX8U3I2AjkRyQW1e3QZmgZBmn/x4gB5l1Wd6S/c6NTlsuiIZ6qp3230ZrN2ogjWlTbX65jS3zTZfIsDBCrqf3+Zra31XthuZClEzSMf2qP2tpPFtX0Dm5De4YkoOaXm2mk3z0e201DSF5lZK3HK2ACA+ftc3UhlO7HMaGsTanAq+6QvWVLi49xVn5lUEfh+nEeOpBMywCLneIsBzUH9H/Jt62g3p10tYMFG+TAFiTfdwK2yWhwrVcFCOK8oDd+GA90v38vz0bIK8JcqvkRejMe3TUS5HocP+Adt5MW/ona/eLHIPVUqpuFULS7A1FIVIWSwnC4WSBpDOy8eNg/lESgawY9MmThb53h0jJiMp3IXgGTSTiF7qRZFimNIj8IjdFBNg=="}}
    [DEBUG 2023/12/15 08:27:49] main.go:43: ctx simple to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test-string","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":"body test","Id":"d2b78dfe-8b68-4b7d-af59-15826c52515a","MD5OfBody":"2e9cc74f6f6aca12eaa2d252df457910","MD5OfMessageAttributes":null,"MessageAttributes":null,"ReceiptHandle":"AQEBfQYdy5nmL+BVJFslmID7vW8WyDnL83OGY8i4YfpfqaAM4xZ2+cwMIPsKoyOVZYII565b6I8XvoJmsZBK7adaGZtUk9r4+w9EZVfe3DCxB4Bt5gYbdb15Kjes8iQAQowuwCZ8e8FSnj4wXEBacPbaOLN0qo1afU3hU3kZ7ZION+0jfiup2AE3yK5/xR7vQtNqTZulNpzda5o1XICsdz/67iLQExhEJqbO3iSrclbzGVaggjPXYOlzx7iFeQrS2hcbQe5boZHCl/W8gC12TRPf/WX2OxNiU9UddY/nBPpy6FJG6I1HPmQla6OMznllE/f4yhjk1a5DzbSEDUjXaaalDFzUBYRbhpEFmb5NbNrsto/myP+QOhoU+DtgdL9iUQen83+MKJDpbCglmdpBYvoIKA=="}}


If you want to obtain the message attributes serialized in Map or Struct, you can call the function
**ReceiveMessage** passing in the second parameter of the handler context indicating the type you want to serialize
message attributes, see:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
)

type messageAttTest struct {
    Name        string  `json:"account,omitempty"`
    Text        string  `json:"text,omitempty"`
    Balance     float64 `json:"balance,omitempty"`
    Bool        bool    `json:"bool"`
    Int         int     `json:"int"`
    SubStruct   test    `json:"subStruct,omitempty"`
    PointerTest *test   `json:"pointerBank,omitempty"`
    Map         map[string]any
    Any         any
    EmptyString string `json:"emptyString,omitempty"`
    HideString  string `json:"-"`
}

func main() {
    sqs.ReceiveMessage(os.Getenv("SQS_QUEUE_TEST_STRING_URL"), handler)
}

func handler(ctx *sqs.Context[test, messageAttTest]) error {
    logger.Debug("ctx to process message:", ctx)
    return nil
}
```

Output:

    [DEBUG 2023/12/15 08:19:08] main.go:43: ctx to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test-string","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":"body test","Id":"0fe8cca6-92ac-4a6d-98de-ec373f433266","MD5OfBody":"2e9cc74f6f6aca12eaa2d252df457910","MD5OfMessageAttributes":null,"MessageAttributes":{"Any":null,"Map":null,"account":"","balance":0,"bool":false,"emptyString":"","int":0,"subStruct":{"Map":null,"bank":{"account":"","balance":0,"digits":""},"birthDate":"0001-01-01T00:00:00Z","name":""},"text":""},"ReceiptHandle":"AQEBFaWjTuCHZcWH9JNlAJt0/cyJBwnK+K5yWVYpCvT+/EBsnSRkEuOp1G6YqSA5szJKc7EXULmDlUxWzQ4pYvCT9BhP6nk82blgcg3Kw5zQy/fym6vg46CdoqLwT0d9vAlnKOsvtkN/dGC6ze0lcSsGjSbUIujo3NC4uXzmflbuanaUIfPo7eju2LGYtw/eAgeVHIdMeH1kkrYfNfTl5iRQejSuSoaozq0V/E3+hfyNhBpmm+J2Bsays+AHvXEZpMNwouycwTpiBWpQzoybUJbXVwHSiQj539hAOGsrmJ2vG+ZQfYb19AClh0gdvoyKsuR6F4Qny0mGqOaGOqD2k/CiuTHkz+W/6nclTKTmaUlUogc+EHCC5qirnKq3i/pCd8UFWGV8PZ7vy48Deh/KYRvLJw=="}}

We can consume messages with body of all types in all message receiving functions,
See below an example consuming a body with structure:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
)

type test struct {
    Name      string    `json:"name,omitempty"`
    BirthDate time.Time `json:"birthDate,omitempty"`
    Emails    []string  `json:"emails,omitempty"`
    Bank      bank      `json:"bank,omitempty"`
    Map       map[string]any
}

type bank struct {
    Account string  `json:"account,omitempty"`
    Digits  string  `json:"digits,omitempty"`
    Balance float64 `json:"balance,omitempty"`
}

func main() {
    sqs.SimpleReceiveMessage(os.Getenv("SQS_QUEUE_TEST_URL"), handler)
}

func handler(ctx *sqs.SimpleContext[test]) error {
    logger.Debug("ctx simple body struct to process message:", ctx)
    return nil
}
```

Output:

    [DEBUG 2023/12/15 09:08:39] main.go:47: ctx simple body struct to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":{"Map":{"bool":true,"float":1.23,"int":1,"string":"text test"},"bank":{"account":"123456","balance":200.12,"digits":"2"},"birthDate":"2023-12-15T09:08:25-03:00","emails":["test@gmail.com","gabriel@gmail.com","gabriel.test@gmail.com"],"name":"Test Name"},"Id":"ac141711-fcf0-4073-a206-78a8a8914756","MD5OfBody":"9d70de56f2ce5bb257a0e88a4ae5bb64","MD5OfMessageAttributes":null,"MessageAttributes":null,"ReceiptHandle":"AQEBOSU25xkDe6sED3CxRGKFrgwTcLpaHpwgsBUSUkG2VXzVmcKLTnXp1JtnbB9fS1S/C1s3R5XWyFfgAasNQyxuFpwYJ07O0ZrTzixElJxUPCrcOxt0glM7Oa7wHeiU4yzT4xfVy7nv4r+oOPcaLUgiWvk0pp6q72Uu/Y51zqIRbhlpRnG189FPbhExfjqvTWjqzHNUzWop5cCMVhqF7F8BKp1LXVADOfbtaBTaJX17lv5x+tHOjjMg8hecJNWAo7k9PZg3+b23Fxe3/k2scjayKFXdsjkpgWrow+VvzJP5FrZ6b/jTgFTTkU4hz6zLmrntRCSqZpwmdU1rb+tCDGH3cHLoZ61vl0KOcgmCLnrkrG3pU0jfP1zOeHzYm+EZV+75LAnbPtjqdN/iZGLqseotvg=="}}

If you want to customize consumer roles, see example below:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
)

func main() {
    opt := option.NewConsumer().
        // HTTP communication customization options with AWS SQS
        SetHttpClient(option.HttpClient{}).
        // print logs (default: false)
        SetDebugMode(true).
        // If true remove the message from the queue after successfully processed (handler error return is null) (default: false)
        SetDeleteMessageProcessedSuccess(true).
        // Duration time to process the message, timeout applied in the past context. (default: 5 seconds)
        SetConsumerMessageTimeout(5 * time.Second).
        // Delay to run the next search for messages in the queue (default: 0)
        SetDelayQueryLoop(5 * time.Second).
        // The maximum number of messages to return. 1 a 10 (default: 10)
        SetMaxNumberOfMessages(10).
        // The maximum number of messages to return. 1 a 10 (default: 0)
        SetVisibilityTimeout(5 * time.Second).
        // The duration that the received messages are hidden from subsequent
        // retrieve requests after being retrieved by a ReceiveMessage request.
        SetReceiveRequestAttemptId("").
        // The duration for which the call waits for a message to arrive in
        // the queue before returning (default: 0)
        SetWaitTimeSeconds(1 * time.Second)
	
    sqs.SimpleReceiveMessage(os.Getenv("SQS_QUEUE_TEST_URL"), handler, opt)
}

func handler(ctx *sqs.SimpleContext[test]) error {
    logger.Debug("ctx simple body struct to process message:", ctx)
    return nil
}
```

Output:

    [INFO 2023/12/15 10:08:44] consumer.go:265: Run start find messages with options: {"Default":{"debugMode":true,"httpClient":{"APIOptions":null,"AppID":"","AuthSchemeResolver":null,"AuthSchemes":null,"BaseEndpoint":null,"ClientLogMode":0,"Credentials":null,"DefaultsMode":"","DisableMessageChecksumValidation":false,"EndpointOptions":{"DisableHTTPS":false,"LogDeprecated":false,"Logger":null,"ResolvedRegion":"","UseDualStackEndpoint":0,"UseFIPSEndpoint":0},"EndpointResolverV2":null,"HTTPClient":null,"HTTPSignerV4":null,"Logger":null,"Region":"","RetryMaxAttempts":0,"RetryMode":"","Retryer":null,"RuntimeEnvironment":{"EC2InstanceMetadataRegion":"","EnvironmentIdentifier":"","Region":""}}},"DeleteMessageProcessedSuccess":true,"ConsumerMessageTimeout":5000000000,"DelayQueryLoop":5000000000,"MaxNumberOfMessages":10,"VisibilityTimeout":5000000000,"ReceiveRequestAttemptId":null,"WaitTimeSeconds":1000000000}
    [INFO 2023/12/15 10:08:44] consumer.go:244: Start process received messages size: 3
    [DEBUG 2023/12/15 10:08:44] main.go:77: ctx simple body struct to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":{"Map":{"balance":10.23,"bool":true,"emptyString":"","int":3,"name":"Name test producer","nil":null,"text":"Text field"},"bank":{"account":"123456","balance":200.12,"digits":"2"},"birthDate":"2023-12-15T09:51:48-03:00","emails":["test@gmail.com","gabriel@gmail.com","gabriel.test@gmail.com"],"name":"Test Name"},"Id":"8a41c561-173e-4c8d-b1a9-645e61cd1e85","MD5OfBody":"f2deace7b685238892eab686ce6d4932","MD5OfMessageAttributes":"7aea2eb88bd4799daf004990a6397c00","MessageAttributes":{"Map":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"}"},"account":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"Name test producer"},"balance":{"BinaryListValues":null,"BinaryValue":null,"DataType":"Number","StringListValues":null,"StringValue":"10.23"},"bool":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"true"},"int":{"BinaryListValues":null,"BinaryValue":null,"DataType":"Number","StringListValues":null,"StringValue":"3"},"pointerBank":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"name\":\"Test Name\",\"birthDate\":\"2023-12-15T09:51:48.015621-03:00\",\"emails\":[\"test@gmail.com\",\"gabriel@gmail.com\",\"gabriel.test@gmail.com\"],\"bank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Map\":{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"},\"emptyString\":\"\",\"pointerBank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Any\":null}"},"subStruct":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"name\":\"Test Name\",\"birthDate\":\"2023-12-15T09:51:48.015621-03:00\",\"emails\":[\"test@gmail.com\",\"gabriel@gmail.com\",\"gabriel.test@gmail.com\"],\"bank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Map\":{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"},\"emptyString\":\"\",\"pointerBank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Any\":null}"},"text":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"Text field"}},"ReceiptHandle":"AQEB8Cs8DmT2np6msZsgZq0BrL7HR5Q5imP4sK89UnwfvIKPLTKn5Ca+j0Cy1xDdSmFONpKtTvyhMDTIX5hIJYMXJmf45rZlf0msHYZ6GmBQoKRJbbDseAIFTmqhJTtsD8oXweBvjSAeVh5hw7kuZCW+bMQf2x6cQyk0WPw7oklHbGos+yVuyzn/OWeVx+7wgkF9E9ozfJfuf3lxTrygXOne7bB3q2unnM4luNIVgcD/Hxni6lIoe0RdOK8rpo7XrxdHWgUh64JcR2bQNZbnY6dPU0Tch6TEkb+sB7Yenvdkr2ZTdMClVhDvFqYjeReBawuLXamtVpiOp+dfTvCy8tEVvkZj0cIV+nCrj3nv2Nsz/jqkv/bM9FF4+O8lMYuHvskTz49SYxDmlwxp6vyMZTDROg=="}}
    [DEBUG 2023/12/15 10:08:44] main.go:77: ctx simple body struct to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":{"Map":{"balance":10.23,"bool":true,"emptyString":"","int":3,"name":"Name test producer","nil":null,"text":"Text field"},"bank":{"account":"123456","balance":200.12,"digits":"2"},"birthDate":"2023-12-15T09:54:27-03:00","emails":["test@gmail.com","gabriel@gmail.com","gabriel.test@gmail.com"],"name":"Test Name"},"Id":"c3e1360c-2eb2-4326-b80e-85d69c8f32f3","MD5OfBody":"a599c012497c800aed1679812981d6da","MD5OfMessageAttributes":"e1fc8c890c4720384ffca66c432cc0ac","MessageAttributes":{"Map":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"}"},"account":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"Name test producer"},"balance":{"BinaryListValues":null,"BinaryValue":null,"DataType":"Number","StringListValues":null,"StringValue":"10.23"},"bool":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"true"},"int":{"BinaryListValues":null,"BinaryValue":null,"DataType":"Number","StringListValues":null,"StringValue":"3"},"pointerBank":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"name\":\"Test Name\",\"birthDate\":\"2023-12-15T09:54:27.558714-03:00\",\"emails\":[\"test@gmail.com\",\"gabriel@gmail.com\",\"gabriel.test@gmail.com\"],\"bank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Map\":{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"},\"emptyString\":\"\",\"pointerBank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Any\":null}"},"subStruct":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"name\":\"Test Name\",\"birthDate\":\"2023-12-15T09:54:27.558714-03:00\",\"emails\":[\"test@gmail.com\",\"gabriel@gmail.com\",\"gabriel.test@gmail.com\"],\"bank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Map\":{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"},\"emptyString\":\"\",\"pointerBank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Any\":null}"},"text":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"Text field"}},"ReceiptHandle":"AQEBiTWyaLxTUXk+sQbCogksDaphcYmOrNhlaUEk6DuREqbRL6/5Lu8tq6nMXAYCWDytUAH7CnNtuh8bXeeB89S4kl24pUSbKS3/OmOADichHPhgNteUANGxDvutZqonUjQW9DqwoqnEXJAapFSSW3vRBmcSOTOQBroDhunIz0BQzbr1ylkw6J11FU0iIzHMWA90vWQj5hIIhK5tZGTbhTIys+DIJLYqdvY4VQR2QRScVbKj4rSIrz7j0RE94+vkP1h8dB+XR1Guq7+E4N68batQq1oAA2SjDa/q3mYkYTtputb1sw55TqvW9boLetwqb98dQlcE6OoMN3vV3Rudl3ixLESNHXn7Od+mDKmHRKPm0irWhP05m0dHfTlVnPGU5dreFO/9DCJ4vEGZjTONbNpxTg=="}}
    [DEBUG 2023/12/15 10:08:44] main.go:77: ctx simple body struct to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":{"Map":{"balance":10.23,"bool":true,"emptyString":"","int":3,"name":"Name test producer","nil":null,"text":"Text field"},"bank":{"account":"123456","balance":200.12,"digits":"2"},"birthDate":"2023-12-15T09:55:12-03:00","emails":["test@gmail.com","gabriel@gmail.com","gabriel.test@gmail.com"],"name":"Test Name"},"Id":"807ccb3f-7bfd-4e8a-83fd-18f461ce7f6b","MD5OfBody":"edcf363ca1c42c81c34fbb3586c4fbb6","MD5OfMessageAttributes":"e7b35866b4dbfa2075553540b18433b5","MessageAttributes":{"Map":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"}"},"account":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"Name test producer"},"balance":{"BinaryListValues":null,"BinaryValue":null,"DataType":"Number","StringListValues":null,"StringValue":"10.23"},"bool":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"true"},"int":{"BinaryListValues":null,"BinaryValue":null,"DataType":"Number","StringListValues":null,"StringValue":"3"},"pointerBank":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"name\":\"Test Name\",\"birthDate\":\"2023-12-15T09:55:12.111056-03:00\",\"emails\":[\"test@gmail.com\",\"gabriel@gmail.com\",\"gabriel.test@gmail.com\"],\"bank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Map\":{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"},\"emptyString\":\"\",\"pointerBank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Any\":null}"},"subStruct":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"{\"name\":\"Test Name\",\"birthDate\":\"2023-12-15T09:55:12.111056-03:00\",\"emails\":[\"test@gmail.com\",\"gabriel@gmail.com\",\"gabriel.test@gmail.com\"],\"bank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Map\":{\"balance\":10.23,\"bool\":true,\"emptyString\":\"\",\"int\":3,\"name\":\"Name test producer\",\"nil\":null,\"text\":\"Text field\"},\"emptyString\":\"\",\"pointerBank\":{\"account\":\"123456\",\"digits\":\"2\",\"balance\":200.12},\"Any\":null}"},"text":{"BinaryListValues":null,"BinaryValue":null,"DataType":"String","StringListValues":null,"StringValue":"Text field"}},"ReceiptHandle":"AQEBw7GpJmzuQjetJgl7VLNVy7ESboXzdTcprQoRxH/xV1aQewH4P5cXmgB8sY/Zoa/QoGSbzq42TbRaBOr3y+htzzv399ysn7vW9FlImpsXWQHgf/TOHr3L4Uv6hwr3KcWZzfSe0nRIzf+YQQzDBf0MayU+NWSfkbcyM1g2FCwByho+SFvLIi5RkYIGpGVh2IpUEP2Bfym5QKHBddV9bXR7Uf32xTQVhHDb+7LxJmhgCyBFawGiRh9ZVVMBP0maZMyHh1keOk7N5Z8v9zVvvzsOMagciMB2ysT9wmbO9kinuPHE683a+mKDdibmnu/Bo0tGTXp2Iv6ksT31nUdM05JwYdCtkoyH7ZU556sY2D/XChc0a+aOJ1PvfHAxe756MmglldU5fdmn30VVsfAzj8/VPQ=="}}
    [INFO 2023/12/15 10:08:44] consumer.go:290: Finish process messages! processed: 3 success: ["807ccb3f-7bfd-4e8a-83fd-18f461ce7f6b"] failed: null


You can also consume messages asynchronously, see:

```go
import (
    "context"
    "github.com/GabrielHCataldo/go-aws-sqs-template/sqs"
    "github.com/GabrielHCataldo/go-logger/logger"
    "os"
)

func main() {
    sqs.SimpleReceiveMessageAsync(os.Getenv("SQS_QUEUE_TEST_URL"), handler, option.NewConsumer().SetDebugMode(true))
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    select {
    case <-c:
    logger.Info("Stopped application!")
    }
}

func handler(ctx *sqs.SimpleContext[test]) error {
    logger.Debug("ctx simple body struct to process message:", ctx)
    return nil
}
```

Output:

    [INFO 2023/12/15 10:12:52] consumer.go:265: Run start find messages with options: {"Default":{"debugMode":true},"DeleteMessageProcessedSuccess":false,"ConsumerMessageTimeout":5000000000,"DelayQueryLoop":5000000000,"MaxNumberOfMessages":10,"VisibilityTimeout":0,"ReceiveRequestAttemptId":null,"WaitTimeSeconds":0}
    [INFO 2023/12/15 10:12:52] consumer.go:240: No msg available to be processed, searching again in 5s
    [INFO 2023/12/15 10:12:57] consumer.go:244: Start process received messages size: 1
    [DEBUG 2023/12/15 10:12:57] main.go:88: ctx simple body struct to process message: {"QueueUrl":"https://sqs.sa-east-1.amazonaws.com/430896945629/go-aws-sqs-template-test","Message":{"Attributes":{"ApproximateFirstReceiveTimestamp":"0001-01-01T00:00:00Z","ApproximateReceiveCount":0,"MessageDeduplicationId":"","MessageGroupId":"","SenderId":"","SentTimestamp":"0001-01-01T00:00:00Z","SequenceNumber":0},"Body":{"Map":{"bool":true,"float":1.23,"int":1,"string":"text test"},"bank":{"account":"123456","balance":200.12,"digits":"2"},"birthDate":"2023-12-15T10:13:28-03:00","emails":["test@gmail.com","gabriel@gmail.com","gabriel.test@gmail.com"],"name":"Test Name"},"Id":"4ce4ed91-9b4b-4820-897b-9da5a8c6c630","MD5OfBody":"ad4c4de0c67e5eb5fb373b308fd38f53","MD5OfMessageAttributes":null,"MessageAttributes":null,"ReceiptHandle":"AQEB96WZYV/s91wxR1n8MUr39NllE4h+oVRZUN2WdCRXQoZEJ4DrV09kqyRRxnRxkyh9J/jTCX6sbyjKj6WzJ7YYkhrwr3ruQBW4BR/S2zz8b94waclJpfTaoq2fAa3SBl2/5nQhrCfsuFGfdDhIANlQNv9fyRMv4Vxsil9cjWauMXM2ilgsrcSnaDX2mRFzxDPzGhTnLJEIXOsJ+5nWb4ex5IVsYc81V8TEfs1c4dYmGc4PNs7s2MXtlXtDy2mOg+YnfgCT5evflCy3oN8qEXNglKqCDEuqqeQU2JXslN76zsObKG8ReZyB9PZIimfnhbM6AhIif5YbSfMX5ZylcyKZGF/9K8qwhAh51Lq3TBMxnHT+tOhslnov8MlECcWPjfQW2HdaXGnBGd/phV83MwlhVw=="}}
    [INFO 2023/12/15 10:14:18] consumer.go:290: Finish process messages! processed: 1 success: ["4ce4ed91-9b4b-4820-897b-9da5a8c6c630"] failed: null

For more consumer examples visit: [All examples consumer](https://github/GabrielHCataldo/go-aws-sqs-template/blob/main/_example/consumer/main.go)

### For more examples

- [Producer](https://github/GabrielHCataldo/go-aws-sqs-template/blob/main/_example/consumer/main.go)
- [Consumer](https://github/GabrielHCataldo/go-aws-sqs-template/blob/main/_example/consumer/main.go)
- [Queue](https://github/GabrielHCataldo/go-aws-sqs-template/blob/main/_example/queue/main.go)
- [Message](https://github/GabrielHCataldo/go-aws-sqs-template/blob/main/_example/message/main.go)

How to contribute
------
Make a pull request, or if you find a bug, open it
an Issues.

License
-------
Distributed under MIT license, see the license file within the code for more details.