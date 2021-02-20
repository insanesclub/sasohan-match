###################
### Build stage ###
###################

FROM    golang:1.15.8 AS builder
RUN     go get github.com/insanesclub/sasohan-match
WORKDIR /go/src/github.com/insanesclub/sasohan-match
RUN     make build

###

FROM    fedora:33
WORKDIR /bin/
COPY    --from=builder /go/src/github.com/insanesclub/sasohan-match/bin/match .
EXPOSE  1324
CMD     ["./match"]