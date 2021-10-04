FROM chromedp/headless-shell:latest AS build
WORKDIR /src
RUN apt-get update && apt-get install -y build-essential git wget
COPY . .
RUN wget https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin && go build -o /src/worldcup .

FROM chromedp/headless-shell:latest AS bin
COPY assets /app/assets
COPY templates /app/templates
COPY --from=build /src/worldcup /app
WORKDIR /app
RUN apt-get update && apt-get install -y dumb-init
RUN chmod +x worldcup
ENTRYPOINT ["dumb-init", "--"]
CMD ["/app/worldcup"]
