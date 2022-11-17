package internal

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type ZapOption struct {
}

func (option *ZapOption) Apply(ctx context.Context) context.Context {
	return ctxzap.ToContext(ctx, GetSingleInst().zapLogger)
}
