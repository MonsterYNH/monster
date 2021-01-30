FROM --platform=arm64 nginx:latest

RUN rm -rf /usr/share/nginx/html/*
COPY ./dist/* /usr/share/nginx/html/
COPY ./docker-entrypoint.d/prepare.sh /docker-entrypoint.d/
RUN chmod +x /docker-entrypoint.d/prepare.sh