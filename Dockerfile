FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./project/cmd/main/main.go 
RUN CGO_ENABLED=0 GOOS=linux go build -o migrationsBuild ./project/cmd/migrations/mig.go 

FROM alpine:3.21

RUN apk --no-cache add ca-certificates postgresql-client busybox

WORKDIR /app 

COPY --from=builder /app/main . 
COPY --from=builder /app/migrationsBuild . 
RUN chmod +x main migrationsBuild 

COPY --from=builder ./app/project/config ./config 
COPY --from=builder ./app/project/migrations ./migrations

COPY --from=builder ./app/scripts/entrypoint.sh .
RUN chmod +x ./entrypoint.sh  

EXPOSE 8080 

CMD [ "./entrypoint.sh" ]
