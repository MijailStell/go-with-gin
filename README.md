# go-with-gin

# docker
docker rm -f go-server
docker rmi go-image
docker build --tag go-image .
docker run -d -p 5000:5000 --name go-server go-image

# local test unit
go test -v ./service

# azuredevops test unit
go test -v ./service 2>&1 | go-junit-report > report.xml

# azuredevops coverage
gocov test company/system/microservices/service | gocov-xml > coverage.xml

# local coverage
go test -v -coverprofile=coverage.txt -covermode count ./service 2>&1 | go-junit-report > report.xml
gocov convert coverage.txt > coverage.json
gocov-xml < coverage.json > coverage.xml
mkdir coverage
gocov-html < coverage.json > coverage/index.html


# local kafka
```bash
$ bin/zookeeper-server-start.sh config/zookeeper.properties
```
```bash
$ bin/kafka-server-start.sh config/server.properties
```
```bash
$ bin/kafka-topics.sh --create --topic searchDocumentEvented --bootstrap-server localhost:9092
```
```bash
$ bin/kafka-topics.sh --create --topic foundDocumentEvented --bootstrap-server localhost:9092
```

```bash
$ bin/kafka-console-consumer.sh --topic foundDocumentEvented --bootstrap-server localhost:9092
```
```bash
$ bin/kafka-console-consumer.sh --topic searchDocumentEvented --bootstrap-server localhost:9092
```

```bash
$ go run .\server.go
```

```bash
$ go run .\consumer\main.go
```


Ahora podemos realizar una petición al api del producer:

```bash
POST http://localhost:5001/api/v1/account/searchDocumentEvent
Content-Type: application/json
POST data:
{
    "document": "72579090"
}
```

```bash
$ bin/kafka-topics.sh --zookeeper zookeeper:2181 --delete --topic foundDocumentEvented
```

```bash
$ bin/kafka-topics.sh --zookeeper zookeeper:2181 --delete --topic searchDocumentEvented
```