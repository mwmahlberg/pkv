# Copyright ©2016 Markus W Mahlberg <markus@mahlberg.io>
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

CC=$(shell which go)
BUILD=$(CC) build $(CFLAGS)
GENERATE=$(CC) generate
INSTALL=$(CC) install
CLEAN=$(CC) clean

.PHONY: all clean build install commit
	
all: clean build
	
clean:
	$(CLEAN)
	$(RM) cmd/bindata.go

cmd/bindata.go: 
	$(GENERATE)

build: cmd/bindata.go
	$(BUILD)
	
install: cmd/bindata.go
	$(INSTALL)

commit: clean cmd/bindata.go