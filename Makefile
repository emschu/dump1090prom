# dump1090prom
# Copyright (C) 2025 emschu[aet]mailbox.org
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public
# License along with this program.
# If not, see <https://www.gnu.org/licenses/>.

SHELL := /bin/bash

GO := GO111MODULE=on go
GO_PATH = $(shell $(GO) env GOPATH)
APP_VERSION_DOT = "1.0.2"

.PHONY: all
all: help

.DEFAULT_GOAL:=help

.PHONY: help
help: ## show this help
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: setup
setup: ## setup for ci
	python -m pip install pandas lxml requests

.PHONY: build
build: ## build the application
	$(GO) build -o dump1090prom

.PHONY: release
release: ## build release binaries
	mkdir -p release
	env GOOS=linux GOARCH=amd64 $(GO) build -o release/dump1090prom-linux-amd64
	env GOOS=linux GOARCH=arm $(GO) build -o release/dump1090prom-linux-arm
	env GOOS=linux GOARCH=arm64 $(GO) build -o release/dump1090prom-linux-arm64
	env GOOS=linux GOARCH=386 $(GO) build -o release/dump1090prom-linux-386
	env GOOS=windows GOARCH=amd64 $(GO) build -o release/dump1090prom-windows-amd64.exe
	env GOOS=windows GOARCH=arm64 $(GO) build -o release/dump1090prom-windows-arm64.exe
	env GOOS=darwin GOARCH=arm64 $(GO) build -o release/dump1090prom-darwin-arm64.bin

.PHONY: airlines
airlines: ## Fetch airlines
	cd dev; python fetch_airlines.py

.PHONY: version
version: ## Set version
	sed -r -i 's/VERSION\s*=\s*"([0-9]+.[0-9]+.[0-9]+)"/VERSION = "'$(APP_VERSION_DOT)'"/g' main.go
	sed -r -i 's/Version: ([0-9]+.[0-9]+.[0-9]+)/Version: '$(APP_VERSION_DOT)'/g' README.md

.PHONY: test
test: ## Run tests
	$(GO) test