FROM golang:1.24-rc-alpine AS build

WORKDIR /usr/src/app

COPY . .

# Prevent image errors
RUN if [ ! -f config.json ]; then cp config_example.json config.json; fi

RUN go build -v -ldflags "-s -w"

FROM alpine AS runner
WORKDIR /usr/src/app
COPY --from=build /usr/src/app/config.json .
COPY --from=build /usr/src/app/pb-purger /usr/local/bin/

ENTRYPOINT [ "pb-purger" ]