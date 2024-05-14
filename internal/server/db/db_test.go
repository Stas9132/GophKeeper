package db

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"testing"
)

func TestDB_Register_Login_Logout(t *testing.T) {
	d, err := NewDB(slog.Default())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = d.M.Down()
	}()
	t.Run("Register", func(t *testing.T) {
		got, err := d.Register(context.Background(), "user1", "password1")
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Fatal(got)
		}
	})
	t.Run("Logout", func(t *testing.T) {
		got, err := d.Logout(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Fatal(got)
		}
	})
	t.Run("Login", func(t *testing.T) {
		got, err := d.Login(context.Background(), "user1", "password1")
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Fatal(got)
		}
	})
	t.Run("Logout", func(t *testing.T) {
		got, err := d.Logout(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Fatal(got)
		}
	})
	t.Run("LoginWithWrongPassword", func(t *testing.T) {
		got, err := d.Login(context.Background(), "user1", "wrang_password1")
		if err != nil {
			t.Fatal(err)
		}
		if got {
			t.Fatal(got)
		}
	})
	t.Run("LoginWithWrongUserPassword", func(t *testing.T) {
		got, err := d.Login(context.Background(), "wrang_user1", "wrang_password1")
		if err != nil {
			t.Fatal(err)
		}
		if got {
			t.Fatal(got)
		}
	})

}

func TestDB_Put_Get(t *testing.T) {
	d, err := NewDB(slog.Default())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = d.M.Down()
	}()
	var save uuid.UUID
	t.Run("Put", func(t *testing.T) {
		got, err := d.PutMeta(context.Background(), "user1", "obj1", 1)
		if err != nil {
			t.Fatal(err)
		}
		if u, e := uuid.Parse(got.ObjId); e != nil {
			t.Fatal(e)
		} else {
			save = u
		}
	})
	t.Run("PutDuplicate", func(t *testing.T) {
		got, err := d.PutMeta(context.Background(), "user1", "obj1", 1)
		if err != nil {
			t.Fatal(err)
		}
		if u, e := uuid.Parse(got.ObjId); e != nil || u != save {
			t.Fatal(e, u, save)
		}
	})
	t.Run("PutDuplicateAnotherUser", func(t *testing.T) {
		got, err := d.PutMeta(context.Background(), "user2", "obj1", 1)
		if err != nil {
			t.Fatal(err)
		}
		if u, e := uuid.Parse(got.ObjId); e != nil || u == save {
			t.Fatal(e, u, save)
		}
	})
	t.Run("Get", func(t *testing.T) {
		got, err := d.GetMeta(context.Background(), "user1", "obj1")
		if err != nil {
			t.Fatal(err)
		}
		if u, e := uuid.Parse(got.ObjId); e != nil || u != save || got.ObjType != 1 {
			t.Fatal(e, u, save)
		}
	})
	t.Run("GetNotExisObject", func(t *testing.T) {
		got, err := d.GetMeta(context.Background(), "user1", "obj2")
		if err == nil {
			t.Fatal(got)
		}
	})
	t.Run("GetAnotherUser", func(t *testing.T) {
		got, err := d.GetMeta(context.Background(), "user2", "obj1")
		if err != nil {
			t.Fatal(err)
		}
		if u, e := uuid.Parse(got.ObjId); e != nil || u == save || got.ObjType != 1 {
			t.Fatal(e, u, save)
		}
	})
	d.Close()
	t.Run("PutCllassedDB", func(t *testing.T) {
		got, err := d.PutMeta(context.Background(), "user2", "obj1", 1)
		if err == nil {
			t.Fatal(got)
		}
	})

}
