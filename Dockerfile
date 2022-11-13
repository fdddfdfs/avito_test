FROM golang
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY go.sum ./

COPY *.go ./

RUN go build -o /AvitoTest

EXPOSE 8080

CMD [ "/AvitoTest" ]