.PHONY: dependencies
dependencies:
	go get github.com/c9s/goprocinfo/linux
	go get github.com/fzzy/radix
	go get github.com/josegonzalez/go-radixurl
	go get github.com/ogier/pflag
	go get github.com/Sirupsen/logrus
	go get github.com/vaughan0/go-ini

.PHONY: install
install: dependencies config build
	sudo cp ./metricsd /usr/local/bin/metricsd

.PHONY: clean
clean:
	rm -f ./metricsd

.PHONY: clean-data
clean-data:
	sudo service elasticsearch stop
	sudo rm -rf /var/lib/elasticsearch/elasticsearch/
	sudo service elasticsearch start
	sudo rm -rf /var/lib/graphite/whisper/servers /var/lib/graphite/whisper/vagrant

.PHONY: config
config:
	sudo mkdir -p /etc/metricsd/
	sudo rm -rf /etc/metricsd/metricsd.ini
	sudo cp config.ini /etc/metricsd/metricsd.ini

.PHONY: build
build:
	go build

.PHONY: run
run: config
	go run main.go

.PHONY: start
start: build config
	./metricsd --config="/etc/metricsd/metricsd.ini" --loglevel=debug

.PHONY: fmt
fmt:
	@find . -name "*.go" -exec go fmt {} \;

.PHONY: tags
tags:
	gotags -tag-relative=true -R=true -sort=true -f="tags" -fields=+l .

.PHONY: ip
ip:
	@ifconfig | grep "inet " | grep -v 127 | grep -v "addr:10.0" | cut -d':' -f2 | cut -d' ' -f1

.PHONY: collector
collector:
	cd templaters && go run collector.go --names=$(names)
