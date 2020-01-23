FROM golang:1.13.5
RUN echo $GOPATH
RUN mkdir -p /go/src/Filebox-Micro/auth-service
COPY . /go/src/Filebox-Micro/auth-service
WORKDIR /go/src/Filebox-Micro/auth-service
RUN go get -d -v ./...
RUN go build -o main .
EXPOSE 8081
CMD [ "./main" ]