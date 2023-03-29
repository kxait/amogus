build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/amogus .

copy:
	test -n "$(PVM_PATH)" # $$PVM_PATH
	cp bin/linux-amd64/amogus ${PVM_PATH}/

clean:
	rm -rf bin