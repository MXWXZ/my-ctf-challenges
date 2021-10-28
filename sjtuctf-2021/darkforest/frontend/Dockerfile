FROM python:3

WORKDIR /app
RUN pip install --no-cache-dir flask redis requests

COPY . .
EXPOSE 8000
CMD ["python3","main.py"]
