package scanners

import (
	"fmt"

	"github.com/bestserversio/spy/internal/servers"
	"github.com/rumblefrog/go-a2s"
)

func QueryA2s(server *servers.Server) error {
	var err error

	if server.Ip == nil || server.Port == nil {
		err = fmt.Errorf("server IP or port is nil")

		return err
	}

	// Format IP/port address string.
	connStr := fmt.Sprintf("%s:%d", *server.Ip, *server.Port)

	cl, err := a2s.NewClient(connStr)

	if err != nil {
		return err
	}

	// Query server information.
	info, err := cl.QueryInfo()

	if err != nil {
		return err
	}

	// Set current users.
	if server.CurUsers == nil {
		server.CurUsers = new(int)
	}

	*server.CurUsers = int(info.Players)

	// Set max users.
	if server.MaxUsers == nil {
		server.MaxUsers = new(int)
	}

	*server.MaxUsers = int(info.MaxPlayers)

	// Set bots.
	if server.Bots == nil {
		server.Bots = new(int)
	}

	*server.Bots = int(info.Bots)

	// Set map name.
	if server.MapName == nil {
		server.MapName = new(string)
	}

	*server.MapName = info.Map

	// Check if secure.
	if server.Secure == nil {
		server.Secure = new(bool)
	}

	*server.Secure = info.VAC

	// Check OS.
	if server.Os == nil {
		server.Os = new(string)
	}

	switch info.ServerOS.String() {
	case "Windows":
		*server.Os = "windows"

	case "Linux":
		*server.Os = "linux"

	case "Mac":
		*server.Os = "Mac"
	}

	// Check if dedicated.
	if server.Dedicated == nil {
		server.Dedicated = new(bool)
	}

	if info.ServerType.String() == "Dedicated" {
		*server.Dedicated = true
	} else {
		*server.Dedicated = false
	}

	// Check for password.
	if server.Password == nil {
		server.Password = new(bool)
	}

	if info.Visibility {
		*server.Password = true
	} else {
		*server.Password = false
	}

	return err
}
