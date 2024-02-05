FROM openjdk:21-jdk-slim as builder

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
RUN apt update && \
    apt install -y build-essential maven

COPY ./src /app/src
COPY ./curl-8.5.0.tar.gz /app/curl-8.5.0.tar.gz

WORKDIR /app

RUN tar -xzvf curl-8.5.0.tar.gz && cd curl-8.5.0 && \
    ./configure --disable-shared --without-ssl --disable-libcurl-option --disable-alt-svc --disable-headers-api --disable-ldap --disable-ldaps --disable-ftp --disable-file --disable-rtsp --disable-dict --disable-gopher --disable-mqtt --disable-smtp --disable-pop3 --disable-tftp --disable-telnet --disable-imap --disable-ipv6 --disable-unix-sockets && \
    make -j
RUN cd src && mvn package

FROM openjdk:21-jdk-slim as prod

COPY ./flag /flag
COPY ./readflag /readflag
COPY --from=0 /app/curl-8.5.0/src/curl /app/new_curl
COPY --from=0 /app/src/target/ezjava-0.0.1-SNAPSHOT.jar /app/ezjava-0.0.1-SNAPSHOT.jar

RUN useradd -M -s /bin/bash ctf
RUN chmod +x /app/new_curl && chown root:root /flag && chmod 0600 /flag
RUN chmod u+s /readflag && chmod +x /readflag

WORKDIR /app
USER ctf
ENV IPFS_GATEWAY=http://kubo:8080
ENTRYPOINT ["java","-jar","ezjava-0.0.1-SNAPSHOT.jar"]