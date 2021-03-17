GOCMD       = go
GOBUILD     = $(GOCMD) build
GORUN       = $(GOCMD) run
BINARY_FILE = bin/match
RM          = rm

all: build

build:
	$(GOBUILD) -gcflags -m -v -o $(BINARY_FILE) .

run:
	$(GORUN) -gcflags -m -v .

clean:
	$(RM) -f $(BINARY_FILE)
