#!/bin/sh

# runs test suites: 1 slave per host, 2 slaves per host, 4 slaves per hosts
# hashrates are saved as hashrate-hosts-slaves

if [ -z ${HOSTS+x} ]; then
    echo "set HOSTS"
    exit 1
fi

hosts=${HOSTS}

SLAVES=$((hosts)) 
SLAVES=$((SLAVES)) ./make-config.sh > amogus.yaml
echo "$((hosts)) hosts, $((SLAVES)) slaves total"
./run.sh
mv hashrate hashrate-$((hosts))-$((SLAVES))

SLAVES=$((hosts * 2)) 
SLAVES=$((SLAVES)) ./make-config.sh > amogus.yaml
echo "$((hosts)) hosts, $((SLAVES)) slaves total"
./run.sh
mv hashrate hashrate-$((hosts))-$((SLAVES))

SLAVES=$((hosts * 4)) 
SLAVES=$((SLAVES)) ./make-config.sh > amogus.yaml
echo "$((hosts)) hosts, $((SLAVES)) slaves total"
./run.sh
mv hashrate hashrate-$((hosts))-$((SLAVES))

