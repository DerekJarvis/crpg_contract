#!/bin/bash
rm --force crpg.zip
./firefly/ff reset dev -f
COMPOSE_HTTP_TIMEOUT=180 ./firefly/ff start dev
./2_package.sh
./3_deploy.sh
spd-say -i -50 -p +30 -t female2 "Fire fly ready, finished in $SECONDS seconds"
