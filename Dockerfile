FROM golang:1.20 as builder

ENV CGO_ENABLED=0

ADD . /FinUslugi
WORKDIR /FinUslugi

RUN go build -o matertial_school cmd/main.go

FROM  alpine:3.19.1
COPY --from=builder /FinUslugi/matertial_school /FinUslugi/matertial_school
COPY --from=builder /FinUslugi/env /FinUslugi/env

WORKDIR /FinUslugi
EXPOSE 8080/tcp
CMD ["/FinUslugi/matertial_school"]