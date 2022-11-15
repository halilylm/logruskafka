
# Kafka hook for logrus

Kafka hook for logrus with sasl support. It only supports one topic for now.

## Installation

```bash 
go get github.com/halilylm/logruskafka@v0.1.0
```
    
## Usage

```go
hook, err := logruskafka.NewKafkaHook(
		"id", // id of the hook
		[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}, // levels to be logged
		&logrus.JSONFormatter{}, // log formatter
		"log-topic", // kafka topic
		[]string{
			"localhost:9093",
			"localhost:9094",
			"localhost:9095",
		}, // brokers
	)
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.Hooks.Add(hook)
	logger.Error("This is an Error msg")
```

With SASL Authentcation

```go
	hook, err := logruskafka.NewKafkaHookWithSaslAuth(
		"log", // topic
		[]logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}, // log levels
		&logrus.JSONFormatter{}, // formatter
		"log-topic", // topic
		[]string{
			"localhost:9093",
			"localhost:9094",
			"localhost:9095",
		}, // brokers
		&logruskafka.KafkaAuth{
			Username:  "username",
			Password:  "password",
			Algorithm: scram.SHA512, // algorithm. if nil, it will use plain mechanism
			TLSConfig: tls.Config{},
		},
	)
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.Hooks.Add(hook)
	logger.Error("This is an Error msg")
```

  
