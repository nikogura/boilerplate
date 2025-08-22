FROM golang:1.24.5 AS builder

WORKDIR /{{.ProjectName}}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build .

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder --chown=nonroot:nonroot /{{.ProjectName}}/{{.ProjectName}} /app/

ENTRYPOINT ["/app/{{.ProjectName}}", "server"]