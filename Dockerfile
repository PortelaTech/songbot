FROM golang:1.11
ADD . .
WORKDIR .
RUN go install ./...
RUN apk add --no-cache ca-certificates && chmod +x code
EXPOSE 80/tcp
CMD [ "./code" ]

