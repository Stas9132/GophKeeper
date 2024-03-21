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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"strconv"
	"strings"
)

type Keys struct {
	Name string
	Type int
}

type Client struct {
	keeper.KeeperClient
	logger.Logger
	user       string
	token      string
	s3         Storage
	storedKeys []Keys
}

type Storage interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}

var ErrInvFormatCommand = errors.New("invalid format command")
var ErrObjectNotFound = errors.New("object not found")

func NewClient(l logger.Logger, tlsCred credentials.TransportCredentials) (*Client, error) {
	s3, err := storage.NewS3(l)
	if err != nil {
		return nil, err
	}
	con, err := grpc.Dial(config.ListenAddress, grpc.WithTransportCredentials(tlsCred))
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
	if len(flds) != 4 {
		fmt.Println("usage: put <key> <type> <data>")
		return ErrInvFormatCommand
	}
	keyName, data := flds[1], flds[3]
	Type, _ := strconv.Atoi(flds[2])
	var dataReader io.Reader
	var dataLen int64
	switch keeper.SyncMain_TypeCode(Type) {
	case keeper.SyncMain_TYPE_LP:
		dataReader = strings.NewReader(data)
		dataLen = int64(len(data))
	case keeper.SyncMain_TYPE_TEXT:
		dataReader = strings.NewReader(data)
		dataLen = int64(len(data))
	case keeper.SyncMain_TYPE_BIN:
		f, err := os.Open(data)
		if err != nil {
			return err
		}
		defer f.Close()
		fst, err := f.Stat()
		if err != nil {
			return err
		}
		dataLen = fst.Size()
	case keeper.SyncMain_TYPE_CARD:
		dataReader = strings.NewReader(data)
		dataLen = int64(len(data))
	default:
		fmt.Println("Types:")
		fmt.Println("	1:Login/Password")
		fmt.Println("	2:Text")
		fmt.Println("	3:Binary file")
		fmt.Println("	4:Bank's Card")
		return ErrInvFormatCommand
	}
	key := Keys{
		Name: keyName,
		Type: Type,
	}
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", c.token))
	info, err := c.s3.PutObject(ctx, c.user, key.Name, dataReader, dataLen, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		c.Error("PutObject receive: " + err.Error())
		return err
	}
	c.Info("Uploaded " + key.Name + " of size: " + strconv.FormatInt(info.Size, 10) + " succesfully.")
	c.storedKeys = append(c.storedKeys, key)

	return nil
}

func (c *Client) Get(flds []string) (string, error) {
	if len(flds) != 2 {
		fmt.Println("usage: get <key>")
		return "", ErrInvFormatCommand
	}
	keyName := flds[1]
	for _, skey := range c.storedKeys {
		if keyName == skey.Name {
			goto keyExist
		}
	}
	return "", ErrObjectNotFound
keyExist:
	object, err := c.s3.GetObject(context.Background(), c.user, keyName, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(object)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) List() ([]Keys, error) {
	return c.storedKeys, nil
}

func (c *Client) SyncList() error {
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", c.token))
	keys := make([]*keeper.SyncMain_KeysMain, len(c.storedKeys))
	for i, skey := range c.storedKeys {
		keys[i] = &keeper.SyncMain_KeysMain{
			Name: skey.Name,
			Type: keeper.SyncMain_TypeCode(skey.Type),
		}
	}
	s, err := c.Sync(ctx, &keeper.SyncMain{Keys: keys})
	if err != nil {
		return err
	}
	c.storedKeys = make([]Keys, len(s.GetKeys()))
	for i, key := range s.GetKeys() {
		c.storedKeys[i] = Keys{
			Name: key.GetName(),
			Type: int(key.GetType()),
		}
	}
	return nil
}
