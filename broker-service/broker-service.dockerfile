#base Go Image
# Build Image
# FROM golang:1.19-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o brokerApp  ./cmd/api

# RUN chmod +x /app/brokerApp

# Run the Code on the Image

FROM alpine:latest

RUN mkdir /app

# For Build
# COPY --from=builder /app/brokerApp /app
COPY brokerApp /app

CMD [ "/app/brokerApp" ]
