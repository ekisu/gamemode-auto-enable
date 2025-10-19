BINARY_NAME=gamemode-auto-enable
CMD_PATH=./cmd/gamemode-auto-enable
SERVICE_FILE=assets/gamemode-auto-enable.service
INSTALL_PATH=/usr/local/bin
SERVICE_INSTALL_PATH=/usr/lib/systemd/user

.PHONY: all install clean

all: 
	go build -o $(BINARY_NAME) $(CMD_PATH)

install: all
	install -Dm755 $(BINARY_NAME) $(DESTDIR)$(INSTALL_PATH)/$(BINARY_NAME)
	install -Dm644 $(SERVICE_FILE) $(DESTDIR)$(SERVICE_INSTALL_PATH)/$(BINARY_NAME).service

clean:
	rm -f $(BINARY_NAME)
