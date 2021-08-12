FROM python:3-alpine

WORKDIR /workdir
COPY docs/requirements.txt /workdir
RUN pip install -r /workdir/requirements.txt

EXPOSE 8000

CMD ["mkdocs", "serve", "--dev-addr", "0.0.0.0:8000"]
