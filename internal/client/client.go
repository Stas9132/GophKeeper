package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/internal/storage"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	keeper.KeeperClient
	logger.Logger
	user       string
	token      string
	s3         Storage
	storedKeys []string
	sync.Mutex
}

type Storage interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}

var ErrInvFormatCommand = errors.New("invalid format command")
var ErrObjectNotFound = errors.New("object not found")

func NewClient(l logger.Logger) (*Client, error) {
	s3, err := storage.NewS3(l)
	if err != nil {
		return nil, err
	}
	con, err := grpc.Dial(config.ListenAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		KeeperClient: keeper.NewKeeperClient(con),
		Logger:       l,
		s3:           s3,
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
	c.user = flds[1]
	c.token = header["authorization"][0]
	return nil
}

func (c *Client) Login(flds []string) error {
	if len(flds) != 3 {
		fmt.Println("usage: login <user name> <password>")
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
	c.user = flds[1]
	c.token = header["authorization"][0]
	return nil
}

func (c *Client) Logout() error {
	var header metadata.MD
	_, err := c.KeeperClient.Logout(context.Background(), &keeper.Empty{}, grpc.Header(&header))
	if err != nil {
		return err
	}
	c.user = ""
	c.token = ""
	return nil
}

func (c *Client) Put(flds []string) error {
	if len(flds) != 3 {
		fmt.Println("usage: put <key> <data>")
		return ErrInvFormatCommand
	}
	key, data := flds[1], flds[2]
	c.Lock()
	defer c.Unlock()
	info, err := c.s3.PutObject(context.Background(), c.user, key, strings.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		c.Error("PutObject receive: " + err.Error())
		return err
	}
	c.Info("Uploaded " + key + " of size: " + strconv.FormatInt(info.Size, 10) + " succesfully.")
	c.storedKeys = append(c.storedKeys, key)

	return nil
}

func (c *Client) Get(flds []string) (string, error) {
	if len(flds) != 2 {
		fmt.Println("usage: get <key>")
		return "", ErrInvFormatCommand
	}
	key := flds[1]
	c.Lock()
	defer c.Unlock()
	for _, skey := range c.storedKeys {
		if key == skey {
			goto keyExist
		}
	}
	return "", ErrObjectNotFound
keyExist:
	object, err := c.s3.GetObject(context.Background(), c.user, key, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(object)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) List() ([]string, error) {
	c.Lock()
	defer c.Unlock()
	return c.storedKeys, nil
}

func (c *Client) SyncList() error {
	c.Lock()
	defer c.Unlock()
	s, err := c.Sync(context.Background(), &keeper.SyncMain{Keys: c.storedKeys})
	if err != nil {
		return err
	}
	c.storedKeys = s.GetKeys()
	return nil
}
