FROM postgres:9.6

RUN apt-get update

RUN apt-get -y install sudo

ADD entities.sql /scripts/
ADD test_data.sql /scripts/

ADD dbSetup.sh /docker-entrypoint-initdb.d/

RUN sudo adduser postgres sudo
USER postgres

RUN /etc/init.d/postgresql start
