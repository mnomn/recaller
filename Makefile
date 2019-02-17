
src = $(wildcard *.go)
prog = route2cloud
.PHONY: build_rpi install test all


# Default values for remote system
# Set environment variables or add to commandline:
# "make REMOTE_IP=192.168.1.182 install_rpi"
REMOTE_IP ?= 0.0.0.0
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
	$(info "REMOTE_IP: $(REMOTE_IP)" )
	$(info REMOTE_USER: $(REMOTE_USER) )
	ssh $(REMOTE_USER)@$(REMOTE_IP) "mkdir -p bin/templates"
	scp $(prog) $(REMOTE_USER)@$(REMOTE_IP):bin/$(prog)
	scp -r templates/* $(REMOTE_USER)@$(REMOTE_IP):bin/templates/
	scp route2cloud@pi.service $(REMOTE_USER)@$(REMOTE_IP):
	ssh $(REMOTE_USER)@$(REMOTE_IP) "sudo cp route2cloud@pi.service /etc/systemd/system/"

#	go install route2cloud

test:
	bin/route2cloud -conf route2cloud.json

