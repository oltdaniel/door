#!/bin/bash

while true; do
    curl http://localhost:8080/stats > stats_$(date "+%F_%H:%M").txt
    sleep 3600
done
