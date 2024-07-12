FROM golang:1.16 as build

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build -o kvs

FROM scratch

COPY --from=build /src/kvs .

# Create a directory in the container to store the certificates
COPY --from=build /src/certificates /certificates

EXPOSE 8080

CMD ["/kvs"]