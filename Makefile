CC=$(shell which go)
BUILD=$(CC) build $(CFLAGS)
GENERATE=$(CC) generate
INSTALL=$(CC) install
CLEAN=$(CC) clean

.PHONY: all clean build install commit
	
all: clean generate build
	
clean:
	$(CLEAN)
	$(RM) cmd/bindata.go

cmd/bindata.go: 
	$(GENERATE)

build:
	$(BUILD)
	
install: cmd/bindata.go
	$(INSTALL)

commit: clean cmd/bindata.go