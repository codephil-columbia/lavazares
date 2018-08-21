FROM postgres:9.6

USER postgres

ADD entities.sql /docker-entrypoint-initdb.d/
ADD test_data.sql /docker-entrypoint-initdb.d/