package server

import (
	"context"
	"fmt"

	kratos "github.com/ory/kratos-client-go"
)

type contextKey struct{ name string }
var KeyKratosSession contextKey = contextKey { name: "kratos_session" }


func contextSetSession(ctx context.Context, session *kratos.Session) context.Context {
	return context.WithValue(ctx, KeyKratosSession, session)
}

func contextGetSession(ctx context.Context) (*kratos.Session, bool) {
	val, ok := ctx.Value(KeyKratosSession).(*kratos.Session)
	return val, ok
}

func contextMustSession(ctx context.Context) *kratos.Session {
	val, ok := contextGetSession(ctx)
	if !ok {
		panic(fmt.Sprintf("expected %s to be present in context", KeyKratosSession.name))
	}

	return val
}
