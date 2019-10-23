FROM golang:1.8

WORKDIR /go/src/github.com/d-sparks/ace-of-trades/scrapers/alphavantage
COPY . .

RUN /usr/local/go/bin/go get -d -v ./...
RUN /usr/local/go/bin/go install -v ./...

CMD ["./scrape_from_docker.sh"]
