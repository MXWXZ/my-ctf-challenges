from php:7.4-apache
ADD php-ffi /php-ffi
ADD flag.c /
ADD flag.h /
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN apt update && apt install -y musl-dev libffi-dev build-essential autoconf
RUN cd / && gcc --share -fPIC flag.c -o flag.so && rm flag.c
RUN cd /php-ffi && phpize && ./configure --with-ffi && make && make install && cd / && rm -rf php-ffi
RUN cd / && chmod 444 flag.h && chmod 555 flag.so
