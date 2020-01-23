FROM ubuntu:16.04
USER root
RUN apt-get update
RUN apt-get install wget
RUN wget https://releases.hashicorp.com/terraform/0.12.18/terraform_0.12.18_linux_amd64.zip
RUN unzip terraform_0.12.18_linux_amd64.zip
RUN mv terraform /usr/bin/
RUN terraform --version