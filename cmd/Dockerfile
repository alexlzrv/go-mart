FROM --platform=linux/x86_64 ubuntu:latest

RUN mkdir /app
WORKDIR /app

ADD cmd/accrual/accrual_linux_amd64 /app

EXPOSE 8090

CMD ./accrual_linux_amd64 -a :8090