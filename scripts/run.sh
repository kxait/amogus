#!/bin/sh

make clean
make 
make copy

#scp bin/linux-amd64/amogus ip-172-31-5-153:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-9-11:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-9-86:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-5-248:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-8-233:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-12-242:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-7-42:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-14-161:pvm3/bin/`pvmgetarch`
#scp bin/linux-amd64/amogus ip-172-31-8-12:pvm3/bin/`pvmgetarch`

amogus ./hashes
