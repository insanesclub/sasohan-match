###################
### Build stage ###
###################

FROM golang:latest AS build

RUN go get github.com/insanesclub/sasohan-match

WORKDIR /go/src/github.com/insanesclub/sasohan-match
RUN make build

###

FROM alpine:3.13.1

COPY --from=build /go/src/github.com/insanesclub/sasohan-match/bin/match /bin

EXPOSE 1324

CMD ./bin/match