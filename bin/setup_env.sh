#!/bin/bash

set -euo pipefail

docker-compose -f docker-compose.dev.yml up -d
