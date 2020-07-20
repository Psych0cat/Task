FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go mod vendor
RUN make
CMD ["./main"]
