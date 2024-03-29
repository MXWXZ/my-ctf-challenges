FROM ubuntu:xenial

ENV DEBIAN_FRONTEND noninteractive

# nginx & php
RUN echo "deb http://ppa.launchpad.net/nginx/stable/ubuntu xenial main" > /etc/apt/sources.list.d/nginx.list
RUN echo "deb http://ppa.launchpad.net/ondrej/php/ubuntu xenial main" > /etc/apt/sources.list.d/php.list

# install packages
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 4F4EA0AAE5267A6C && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 00A6F0A3C300EE8C && \
    apt-get update && \
    apt-get -y --force-yes --no-install-recommends install \
    supervisor nginx php5.6-cli php5.6-fpm php5.6-common php5.6-mysql php5.6-mbstring

# configure NGINX as non-daemon
RUN echo "daemon off;" >> /etc/nginx/nginx.conf

# configure php-fpm as non-daemon
RUN sed -i -e "s/;daemonize\s*=\s*yes/daemonize = no/g" /etc/php/5.6/fpm/php-fpm.conf

# clear apt cache and remove unnecessary packages
RUN apt-get autoclean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

# copy config file for Supervisor
COPY config/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# backup default default config for NGINX
RUN cp /etc/nginx/sites-available/default /etc/nginx/sites-available/default.bak

# copy local defualt config file for NGINX
COPY config/nginx/default /etc/nginx/sites-available/default

# php-fpm5.6 will not start if this directory does not exist
RUN mkdir /run/php

RUN rm -rf /var/www/html/* && mkdir /var/www/html/uploads && chmod 777 /var/www/html/uploads
ADD index.html /var/www/html
ADD upload.php /var/www/html
ADD flag /

# NGINX ports
EXPOSE 80 443

CMD ["/usr/bin/supervisord"]
