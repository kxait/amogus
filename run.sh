#!/bin/sh

make clean
make 
make copy
scp bin/linux-amd64/amogus aws1:pvm3/bin/`pvmgetarch`
scp bin/linux-amd64/amogus aws2:pvm3/bin/`pvmgetarch`
amogus ./hashes
