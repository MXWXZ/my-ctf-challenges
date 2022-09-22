FROM golang:alpine as builder
ADD ohf /ohf
WORKDIR /ohf
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o ohf ohf_main_to_be_deployed.go

FROM node:alpine as prod
RUN npm install less -g
ADD ohf /ohf
COPY --from=0 /ohf/ohf /ohf/ohf
COPY readflag /readflag
RUN addgroup ctf && adduser -s /bin/sh -G ctf -H -D ctf
RUN echo "flag{123456}" >/flag
RUN chown root:root /flag && chmod 0600 /flag
RUN chmod u+s /readflag && chmod +x /readflag
WORKDIR /ohf
USER ctf
ENTRYPOINT ["./ohf"]
