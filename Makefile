
src = $(wildcard *.go)
prog = recaller
.PHONY: build_rpi install test all

# Default values for remote system
# Set environment variables or add to commandline:
# "make REMOTE_HOST=192.168.1.182 install_rpi"
REMOTE_HOST ?= raspberrypi

build:
	$(info "BUILD for this machine")
	go build cmd/recaller

go_fmt:
	go fmt

build_rpi:
	env GOOS=linux GOARCH=arm GOARM=5 go build cmd/recaller

install:
	go install recaller

install_rpi: build_rpi
	$(info "REMOTE_HOST: $(REMOTE_HOST)" )
	scp $(prog) "$(REMOTE_HOST):"
	scp scripts/install_r2c.sh "$(REMOTE_HOST):"
	ssh $(REMOTE_HOST) "sudo bash install_r2c.sh -r"

test:
	go test
