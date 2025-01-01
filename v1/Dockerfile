FROM python:3.11-bullseye

WORKDIR /app

COPY . .

RUN pip install -r requirements.txt

RUN chmod +x si-validator.py

RUN ln --symbolic /app/si-validator.py /usr/bin/si-validator

WORKDIR /tmp

ENTRYPOINT ["si-validator"]