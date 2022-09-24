FROM apline:latest

RUN mkdir /app

COPY  listenerApp /app

CMD [ "/app/listenerApp" ]