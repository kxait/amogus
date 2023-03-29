build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/amogus .

clean:
	rm -rf bin