// Copyright 2016 Google, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	auth "grpc-login-server/go-client/api/v1"
	"log"
)



func main() {
	serverAddr := "127.0.0.1:16888"

	var (
		username   = flag.String("username", "", "Username to use.")
	)
	flag.Parse()


	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ac := auth.NewAuthClient(conn)

	fmt.Println("enter password:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}

	req := &auth.LoginRequest{
		Username: *username,
		Password: string(password),
	}
	lm, err := ac.Login(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

/*
	req := &auth.RegistrationRequest{
		Username: *username,
		Password: string(password),
	}
	lm, err := ac.Registration(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
*/
	fmt.Println("Client:", lm.Token)
}
