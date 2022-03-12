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