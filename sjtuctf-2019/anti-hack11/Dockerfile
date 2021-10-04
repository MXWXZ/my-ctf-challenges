FROM ubuntu:latest

RUN apt-get update -y && \
    apt-get install -y --force-yes --no-install-recommends python3-pip python3-dev build-essential
RUN apt-get autoclean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

WORKDIR /app
RUN pip3 install flask pillow
COPY . /app
RUN g++ BP.cpp Network.cpp -o Network
ENTRYPOINT [ "python3" ]
CMD [ "run.py" ]