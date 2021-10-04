FROM ubuntu:20.04
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai
RUN apt-get update -y && apt-get install -y tzdata
RUN apt-get install -y --force-yes --no-install-recommends cmake python3-pip python3-dev build-essential
RUN apt-get autoclean && apt-get -y autoremove && rm -rf /var/lib/apt/lists/*
RUN pip install flask gunicorn numpy face_recognition pillow

ADD ./app.py /app/app.py
WORKDIR /app
RUN mkdir upload
EXPOSE 8000
CMD ["python3","app.py"]
#CMD ["gunicorn", "-b", "0.0.0.0:8000","--timeout","0", "app:app"]
