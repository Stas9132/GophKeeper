package client

import (
	"context"
	"fmt"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	con *grpc.ClientConn
	keeper.KeeperClient
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Dial() error {
	var err error
	c.con, err = grpc.Dial(config.ListenAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c.KeeperClient = keeper.NewKeeperClient(c.con)
	return nil
}

func (c *Client) Close() error {
	return c.con.Close()
}

func (c *Client) Health() error {
	h, err := c.KeeperClient.Health(context.Background(), &keeper.Empty{})
	if err != nil {
		return err
	}
	fmt.Println(h)
	return nil
}
