FROM keymetrics/pm2:12-slim

RUN mkdir -p /app
COPY src /app
RUN chmod 0755 -R /app
RUN chmod 0777 /app/app/public/css/style.css
RUN rm -rf /app/.git
COPY flag /
RUN chmod 0444 /flag

RUN useradd -s /bin/bash -m 0ops
WORKDIR /app/
USER 0ops

EXPOSE 3000
CMD pm2-runtime app.js -i $CPUS
