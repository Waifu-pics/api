FROM golang:alpine AS gobuilder
WORKDIR /src/server
COPY ./src/server /src/server
RUN cd /src/server && go build -o goapp

FROM node:lts-alpine AS vuebuilder
COPY . .
WORKDIR /src/client
RUN npm install
RUN npm run build

FROM alpine
WORKDIR /app
COPY --from=gobuilder /src/server/goapp /app
COPY --from=vuebuilder /src/client/dist /app/dist

ENTRYPOINT ./goapp