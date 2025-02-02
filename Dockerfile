FROM golang:1.23.5

WORKDIR /app
ENTRYPOINT [ "tail", "-f", "/dev/null" ]

