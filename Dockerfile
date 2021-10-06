FROM python:3.10.0

COPY requirements-freeze.txt .
RUN pip install -r requirements-freeze.txt

RUN useradd -m app
COPY --chown=app src /app
WORKDIR /app
USER app

CMD ["/app/serve.py"]
