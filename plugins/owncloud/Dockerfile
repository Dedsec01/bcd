FROM bytesized/debian-base
MAINTAINER test@test.com

RUN apt-get update \
    && apt-get install -y apt-utils \
    && { \
        echo debconf debconf/frontend select Noninteractive; \
        echo mysql-community-server mysql-community-server/data-dir \
            select ''; \
        echo mysql-community-server mysql-community-server/root-pass \
            password 'JohnUskglass'; \
        echo mysql-community-server mysql-community-server/re-root-pass \
            password 'JohnUskglass'; \
        echo mysql-community-server mysql-community-server/remove-test-db \
            select true; \
    } | debconf-set-selections \
    && apt-get -y install apache2 mysql-server libapache2-mod-php5 php5-gd php5-json php5-mysql php5-curl  php5-intl php5-mcrypt php5-imagick

RUN apt-get -y install wget
RUN apt-get -y install curl
RUN apt-get -y install bzip2
RUN apt-get -y install unzip
RUN wget https://download.owncloud.org/community/owncloud-9.1.2.tar.bz2
RUN wget https://owncloud.org/owncloud.asc
RUN tar xjvf owncloud-9.1.2.tar.bz2
RUN cp -r owncloud /var/www

VOLUME /config /data /media
