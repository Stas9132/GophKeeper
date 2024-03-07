package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	keeper.KeeperClient
	logger.Logger
	token string
}

var ErrInvFormatCommand = errors.New("invalid format command")

func NewClient(l logger.Logger) (*Client, error) {
	con, err := grpc.Dial(config.ListenAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		KeeperClient: keeper.NewKeeperClient(con),
		Logger:       l,
	}, nil
}

func (c *Client) Health() error {
	ctx := context.Background()
	if len(c.token) > 0 {
		ctx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"Authorization": c.token}))
	}
	h, err := c.KeeperClient.Health(ctx, &keeper.Empty{})
	if err != nil {
		return err
	}
	c.Info(h.String())
	return nil
}

func (c *Client) Register(flds []string) error {
	if len(flds) != 3 {
		fmt.Println("usage: register <user name> <password>")
		return ErrInvFormatCommand
	}
	var header metadata.MD
	_, err := c.KeeperClient.Register(context.Background(), &keeper.AuthMain{
		User:     flds[1],
		Password: flds[2],
	}, grpc.Header(&header))
	if err != nil {
		return err
	}
	c.token = header["authorization"][0]
	return nil
}

func (c *Client) Login(flds []string) error {
	if len(flds) != 3 {
		fmt.Println("usage: register <user name> <password>")
		return ErrInvFormatCommand
	}
	var header metadata.MD
	_, err := c.KeeperClient.Login(context.Background(), &keeper.AuthMain{
		User:     flds[1],
		Password: flds[2],
	}, grpc.Header(&header))
	if err != nil {
		return err
	}
	c.token = header["authorization"][0]
	return nil
}

func (c *Client) Logout() error {
	var header metadata.MD
	_, err := c.KeeperClient.Logout(context.Background(), &keeper.Empty{}, grpc.Header(&header))
	if err != nil {
		return err
	}
	c.token = header["authorization"][0]
	return nil
}
