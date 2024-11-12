# Copyright 2024 Notedown Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: hygiene verify tidy features gifs dirty format licenser dev

hygiene: tidy format licenser

tidy:
	nix develop --command go mod tidy

features:
	nix develop --command go run features/main.go

gifs: 
	nix develop --command go run features/main.go --generate-gifs

dirty:
	nix develop --command git diff --exit-code

format:
	nix develop --command gofmt -w .

licenser:
	nix develop --command licenser apply -r "Notedown Authors"

dev:
	nix develop --command sandbox/run.sh

run:
	nix develop --command go run main.go

install:
	nix develop --command go install -ldflags "-X 'github.com/notedownorg/task/cmd.CommitHash=$(shell git rev-parse HEAD)'"


