FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY mailerApp /app

CMD [ "/app/mailerApp" ]