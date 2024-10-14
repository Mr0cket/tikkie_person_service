# Fix the build  version
FROM golang:1.22 as build
WORKDIR /build
# Copy dependencies list
COPY go.mod go.sum ./
# Build with optional lambda.norpc tag
COPY . ./
ARG GOARCH=amd64
ARG GOOS=linux
ENV GOARCH=${GOARCH} GOOS=${GOOS}
RUN go build -tags lambda.norpc -o main ./cmd/lambda


# Final image
FROM public.ecr.aws/lambda/provided:al2023-x86_64
COPY --from=build /build/main /main
ENTRYPOINT [ "/main" ]
