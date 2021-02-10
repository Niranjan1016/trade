FROM golang:1.14.3-alpine As build

WORKDIR $GOPATH/src/github.com/niranjan1016/trade 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/trade .

FROM scratch

COPY --from=build /bin/trade /bin/trade 

EXPOSE 8080

ENTRYPOINT ["/bin/trade"]
