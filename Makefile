
src = $(wildcard *.go)
prog = route2cloud
.PHONY: build_rpi install test all

# Default values for remote system
# Set environment variables or add to commandline:
# "make REMOTE_HOST=192.168.1.182 install_rpi"
REMOTE_HOST ?= raspberrypi
REMOTE_USER ?= pi

build:
	$(info "BUILD for this machine")
	go build

go_fmt:
	go fmt

build_rpi:
	env GOOS=linux GOARCH=arm GOARM=5 go build

install:
	go install route2cloud

generate_systemd:
	./generate_systemd.sh $(REMOTE_USER)

install_rpi: build_rpi generate_systemd
	$(info "REMOTE_HOST: $(REMOTE_HOST)" )
	$(info REMOTE_USER: $(REMOTE_USER) )
	ssh $(REMOTE_USER)@$(REMOTE_HOST) "sudo systemctl stop route2cloud@$(REMOTE_USER).service"
	ssh $(REMOTE_USER)@$(REMOTE_HOST) "mkdir -p bin/templates"

	scp $(prog) $(REMOTE_USER)@$(REMOTE_HOST):bin/$(prog)
	scp -r templates/* $(REMOTE_USER)@$(REMOTE_HOST):bin/templates/
	scp route2cloud@pi.service $(REMOTE_USER)@$(REMOTE_HOST):
	ssh $(REMOTE_USER)@$(REMOTE_HOST) "sudo cp route2cloud@pi.service /etc/systemd/system/"

#	go install route2cloud

test:
	go test

