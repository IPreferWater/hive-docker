FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y imagemagick

RUN useradd -m laboratoire_user

USER laboratoire_user

CMD ["convert","/home/laboratoire_user/duck.webp", "/home/laboratoire_user/wonderfull_duck.tif"]