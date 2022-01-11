## ./dev-ops/prod.Dockerfile

FROM golang:1.15-alpine3.12

RUN apk add git

WORKDIR /go/src/gitlab.com/dashers/image-delivery-service

## same as in dev... add files with dependencies requierments
ENV GO111MODULE=on
ADD ./go.mod .
ADD ./go.sum .

## pull in any dependencies
RUN go mod download

## add source files
ADD . .

## our app will now successfully build with the necessary go dependencies included
## creates ./main binary executable file
RUN go build -o main .

## runs our newly created binary executable
CMD ["./main"]