FROM ubuntu:20.04

RUN apt-get update -y && apt-get install -y --force-yes --no-install-recommends python3-pip python3-dev build-essential
RUN apt-get autoclean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*

WORKDIR /app
RUN pip3 install flask pillow bfi
COPY . /app
RUN cd brainfuck && make && make PREFIX="/usr" install
ENTRYPOINT [ "python3" ]
CMD [ "run.py" ]
