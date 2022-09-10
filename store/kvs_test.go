package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/chmikata/go_todo_app/config"
	"github.com/chmikata/go_todo_app/entity"
	"github.com/chmikata/go_todo_app/testutil"
	"github.com/go-redis/redis/v8"
)

func TestNewKvs(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	port := 36379
	if _, defined := os.LookupEnv("CI"); defined {
		port = 6379
	}
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "localhost", port)})
	cli.Ping(context.Background())
	tests := map[string]struct {
		args    args
		want    *KVS
		wantErr bool
	}{
		"ok case": {
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{
					RedisHost: "localhost",
					RedisPort: 36379,
				},
			},
			want:    nil,
			wantErr: false,
		},
		"error case": {
			args: args{
				ctx: context.Background(),
				cfg: &config.Config{
					RedisHost: "localhost",
					RedisPort: 36378,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			got, err := NewKvs(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKvs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// hack because it is a DB connection
			tt.want = got
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKvs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	kvs := &KVS{Cli: cli}
	key := "TestKVS_Save"
	uid := entity.UserId(1234)
	ctx := context.Background()

	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	if err := kvs.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, but got: %v", err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	kvs := &KVS{Cli: cli}

	t.Run("ok case", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_ok"
		uid := entity.UserId(1234)
		ctx := context.Background()
		cli.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})
		got, err := kvs.Load(ctx, key)
		if err != nil {
			t.Fatalf("want no error, but got %v", err)
		}
		if got != uid {
			t.Errorf("want %d, but got %d", uid, got)
		}
	})

	t.Run("notFound case", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_notFound"
		ctx := context.Background()
		got, err := kvs.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v(value == %d)", ErrNotFound, err, got)
		}
	})
}
