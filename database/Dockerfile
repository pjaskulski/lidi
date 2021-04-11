FROM mysql:8.0.20

ENV MYSQL_ROOT_PASSWORD: password \
    MYSQL_DATABASE: dictionary \
    MYSQL_USER: web \
    MYSQL_PASSWORD: pass

WORKDIR /database
COPY ./database .

RUN mysql < database.sql

EXPOSE 3306