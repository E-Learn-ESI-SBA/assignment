FROM golang:alpine3.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags 'popstgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY . .
RUN go build -o main .
RUN chmod +x main
# Create Sys user named assignments , with group assignments
RUN addgroup -S assignment && adduser -S assignments -G assignments
EXPOSE 8080
USER assignments
CMD [ "./main" ]








# Set the Current Working Directory inside the container