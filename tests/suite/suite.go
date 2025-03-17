package suite

import (
	"context"
	ssov1 "github.com/mefedraw/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpcAuth/internal/config"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const grpcHost = "localhost"

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	basePath, err := os.Getwd()
	if err != nil {
		slog.Error("failed to get current working directory")
	}
	if filepath.Base(basePath) == "tests" {
		basePath = filepath.Dir(basePath)
	}

	cfgPath := filepath.Join(basePath, "config", "local.yaml")
	cfg := config.MustLoadByPath(cfgPath)

	ctx, cancelCtx := context.WithCancel(context.Background())

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connect error:", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
