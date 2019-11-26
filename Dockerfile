FROM golang:1.12.9


MAINTAINER Marian Craciunescu

LABEL Description="Symantecbeat "

RUN \
    apt-get update \
      && apt-get install -y --no-install-recommends \
         netcat \
         python-pip \
         virtualenv \
      && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip
RUN pip install --upgrade setuptools
RUN pip install --upgrade docker-compose==1.23.2

RUN mkdir /plugin
COPY symantecbeat.yml /plugin/
#COPY fields.yml /plugin/
COPY symantecbeat /plugin/symantecbeat
COPY ecs/ecs_translating_mapping.csv /plugin/ecs_translating_mapping.csv

WORKDIR /plugin

ENTRYPOINT ["/plugin/symantecbeat" ,"-e" ]
