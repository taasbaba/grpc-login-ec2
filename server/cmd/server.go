package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"grpc-login-server/server/internal/logger"
	"grpc-login-server/server/internal/protocol/grpc"
	v1 "grpc-login-server/server/internal/service/v1"
	"net"
)

// RunServer runs gRPC server
func RunServer() error {
	ctx := context.Background()
	gRPCport := "16888"

	// initialize logger
	if err := logger.Init(/*cfg.LogLevel*/-1, "2006-01-02T15:04:05.999999999Z07:00"); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	//host,_ := externalIP()
	host := "HostIP:DBPort"

	fmt.Println("host:", host)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		"root",
		"password",
		host,
		"game",
		)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1AuthServer := v1.NewAuthServer(db)
	return grpc.RunServer(ctx, v1AuthServer, gRPCport)//cfg.GRPCPort)
}


func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
