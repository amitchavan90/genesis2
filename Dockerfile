# Build the Go API
FROM golang:1.14.4-buster AS go_builder
ARG GOSUMDB_DEFAULT=""
ENV GOSUMDB=$GOSUMDB_DEFAULT
ARG GOPROXY_DEFAULT=""
ENV GOPROXY=$GOPROXY_DEFAULT
ADD . /app
WORKDIR /app/server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main ./cmd/platform

# Build the React application
FROM node:lts AS js_builder
ARG FONTAWESOME_TOKEN_DEFAULT=""
ENV FONTAWESOME_TOKEN=$FONTAWESOME_TOKEN_DEFAULT
COPY --from=go_builder /app/web ./
RUN npm install
RUN npm run build

# Final stage build, this will be the container
# that we will deploy to production
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=go_builder /main ./
COPY --from=go_builder /app/server/rest_manifest_tpl_list.html ./
COPY --from=js_builder /dist ./web
RUN chmod +x ./main
CMD ["./main", "serve"]
