# amogus

ඞ

### usage
`amogus hashes_path [config_path] [output_path]`
- hashes_path is just a list of hash strings in the required format
- config path is the yaml file with the parameters
- output_path is where the cracked passwords will be written to

### cfg
config file is like
```
length_start: 6
length_end: 6
characters: abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPRSTUVWXYZ012345678
mode: sha512
slaves: 4
chunk_size: 500000
test_suite_sample_size: 12
```
most parameters are self explanatory, `chunk_size` is the amount of hashes that will be generated and checked against the input file, use this to manage RAM usage
test_suite_sample_size is optional - if defined, will create a file `hashrate` in your current dir and write the hashrate there every 5 seconds. once the number of samples reaches the parameter value, the program exits

### currently implemented modes
- sha512
- sha256
- shadow
  - sha512

### requirements
- preferably a linux environment
- `pvm` and `pvm-dev` installed on your system
- the paths set up correctly for pvm
- `go`, `gcc`, etc

### running test suite
this project was made for testing the performance of parallel computing. there is an automated system for running the program with different numbers of nodes per host:
- hosts (1 node per host)
- 2*hosts (2 nodes per host)
- 4*hosts (4 nodes per host)

where hosts are the var passed to `scripts/run-test-suite.sh`, for example:
- `HOSTS=10 ./scripts/run-test-suite.sh` (the script uses relative paths so it's important to run it from the project root)

this will run the program 3 different times, collect 14 samples for each of the different modes and exit. the hashrate samples will be saved as `hashrate-X-Y` where X is count of hosts and Y is count of total slaves


to change different parameters during the test, modify the config template in `scripts/make-config.sh`

for tasks required before running the program, add them to `scripts/run.sh` (e.g. copying the executable to different hosts)
