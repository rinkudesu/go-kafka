module github.com/rinkudesu/go-kafka

go 1.20

require (
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/rinkudesu/go-kafka/configuration => ./configuration
	github.com/rinkudesu/go-kafka/producer => ./producer
	github.com/rinkudesu/go-kafka/subscriber => ./subscriber
)
