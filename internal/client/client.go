package client

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
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
	user  string
	token string
	s3    S3
}

type S3 interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}

var ErrInvFormatCommand = errors.New("invalid format command")

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

var dirFS = os.DirFS(".")

const ChunkSize = 1024

func (c *Client) Put(flds []string) error {
	if len(flds) != 4 {
		fmt.Println("usage: put <key> <type> <data>")
		return ErrInvFormatCommand
	}
	keyName, data := flds[1], flds[3]
	Type, _ := strconv.Atoi(flds[2])
	var dataReader io.Reader
	var dataLen int64
	switch keeper.TypeCode(Type) {
	case keeper.TypeCode_TYPE_LOGIN_PASSWORD:
		dataReader = strings.NewReader(data)
		dataLen = int64(len(data))
	case keeper.TypeCode_TYPE_TEXT:
		dataReader = strings.NewReader(data)
		dataLen = int64(len(data))
	case keeper.TypeCode_TYPE_BIN:
		f, err := dirFS.Open(data)
		if err != nil {
			return err
		}
		defer f.Close()
		fst, err := f.Stat()
		if err != nil {
			return err
		}
		dataReader = f
		dataLen = fst.Size()
	case keeper.TypeCode_TYPE_CARD:
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
	block, err := aes.NewCipher(config.AESKey)
	if err != nil {
		return err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(dataReader)
	if err != nil {
		return err
	}
	bytes = aesgcm.Seal(nil, config.AESnonce, bytes, nil)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", c.token))
	putClient, err := c.KeeperClient.Put(ctx)
	if err != nil {
		return err
	}
	for start := 0; start < len(bytes); start += ChunkSize {
		end := start + ChunkSize
		if end > len(bytes) {
			end = len(bytes)
		}
		if err = putClient.Send(&keeper.ObjMain{
			Name:    key.Name,
			Type:    keeper.TypeCode(key.Type),
			EncData: bytes[start:end],
			Size:    int64(len(bytes)),
		}); err != nil {
			c.Error("Chank send errror", "error", err)
			return err
		}
	}
	if err = putClient.CloseSend(); err != nil {
		c.Error("CloseSend", "error", err)
	}

	c.Info("Uploaded " + key.Name + " of size: " + strconv.FormatInt(dataLen, 10) + " succesfully.")

	return nil
}

func (c *Client) Get(flds []string) (string, error) {
	if len(flds) != 2 {
		fmt.Println("usage: get <key>")
		return "", ErrInvFormatCommand
	}
	keyName := flds[1]
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", c.token))
	out, err := c.KeeperClient.Get(ctx, &keeper.ObjMain{
		Name: keyName,
	})
	if err != nil {
		return "", err
	}

	object, err := c.s3.GetObject(ctx, c.user, out.S3Link, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(object)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(config.AESKey)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	bytes, err = aesgcm.Open(nil, config.AESnonce, bytes, nil)
	if err != nil {
		return "", err
	}

	return out.GetType().String() + ": " + string(bytes), nil
}
