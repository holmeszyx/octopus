gofiles=crypt.go \
		main.go


octopus: $(gofiles)
	go build -o octopus $(gofiles)

build: octopus

clean:
	-rm octopus