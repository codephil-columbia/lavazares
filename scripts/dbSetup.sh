#!/usr/bin/sh
sudo -u postgres createuser --superuser codephil
sudo -u postgres createdb -E UTF8 -l en_US.UTF8 -T template0 -O codephil lavazaresDB
sudo -u postgres psql -d lavazaresDB <<EOL
ALTER ROLE codephil WITH PASSWORD 'password';
EOL

sudo -u postgres psql lavazaresDB < /scripts/user.sql