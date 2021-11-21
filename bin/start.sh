#!/bin/bash

# Exit on fail
set -e

./bin/inject_port_into_nginx.sh

nginx -c /etc/nginx/conf.d/default.conf

# Start server
/app/backend
