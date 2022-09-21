FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY mailerApp /app
COPY ./templates /app/templates

CMD [ "/app/mailerApp" ]