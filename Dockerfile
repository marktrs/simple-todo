####################################################################
# Builder Stage                                                    #
####################################################################
FROM golang:1.20-alpine3.16 AS builder
# Set working directory.
WORKDIR /app
# Copy all the code and stuff to compile everything
COPY . .
# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go


####################################################################
# Final Stage                                                      #
####################################################################
# Moving the binary to the 'final Image' to make it smaller
FROM alpine:3.16 as release
# Set working directory.
WORKDIR /app
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main .
COPY .env .
COPY start.sh .
COPY wait-for.sh .

RUN mkdir temp

# Exposes port 3000 because our program listens on that port
EXPOSE 3000
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
