GO = go
GOFLAGS = -ldflags "-s -w"
BIN = corsairtweets

# Important since it is hosted on a very old system.
export CGO_ENABLED=0

all: $(BIN)

$(BIN):
	$(GO) build $(GOFLAGS) -o $@ .

clean:
	rm -f $(BIN)

.PHONY: all $(BIN) clean

-include deploy.mk
