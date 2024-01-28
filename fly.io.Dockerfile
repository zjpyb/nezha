FROM ghcr.io/naiba/nezha-dashboard:latest

COPY ./script/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY ./resource/ /dashboard/resource

COPY ./script/config.yaml /dashboard/config.yaml