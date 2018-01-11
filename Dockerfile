FROM golang:1.9.2-alpine as server
WORKDIR /go/src/github.com/Stratoscale/logserver
COPY . .
RUN go build -o /logserver

FROM node:8.9.3-alpine as client
RUN apk add --no-cache git
COPY ./client /client
COPY ./.git /.git
WORKDIR /client
RUN yarn
RUN npm run build

FROM alpine:3.7
COPY --from=server /logserver /usr/bin/logserver
COPY --from=client /client/dist /client/dist
ENTRYPOINT ["logserver"]

