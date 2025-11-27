FROM golang:1.24-alpine AS build
ENV TZ=Asia/Tokyo

WORKDIR /opt/app

COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .
RUN go build -o /bin/bot ./main.go


FROM gcr.io/distroless/base:nonroot AS runner
ENV TZ=Asia/Tokyo
ENV GOENV=production

COPY --from=build /bin/bot /bin/bot

USER nonroot
ENTRYPOINT ["/bin/bot"]
