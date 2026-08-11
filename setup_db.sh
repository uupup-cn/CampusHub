#!/bin/bash
sudo pg_ctlcluster 16 main start 2>&1
sudo -u postgres psql -c "CREATE USER chb WITH PASSWORD 'chb_test_2024' CREATEDB;" 2>&1
sudo -u postgres psql -c "CREATE DATABASE chb_test OWNER chb;" 2>&1
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE chb_test TO chb;" 2>&1
echo "DONE"
