package identity

import (
	"context"

	appidentity "github.com/lynx-go/lynx-start-demo/internal/app/identity"

	"github.com/lynx-go/grpcapi/authn"
	"github.com/lynx-go/grpcapi/principal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BearerAuthenticator 把 Authorization: Bearer <token> 变成 Principal
// （grpcapi authn SPI），是本模块对整个应用的横切贡献——其他模块的
// 方法只要不声明 anonymous，就自动要求本认证器注入的 principal。
//
// SPI 约定：(nil, nil) = 不认识这种凭证族（非 Bearer），链上继续尝试
// 下一个认证器；(nil, err) = 认识但拒绝。
type BearerAuthenticator struct {
	uc *appidentity.UseCase
}

// NewBearerAuthenticator 构造认证器（Assemble 调用）。
func NewBearerAuthenticator(uc *appidentity.UseCase) *BearerAuthenticator {
	return &BearerAuthenticator{uc: uc}
}

var _ authn.Authenticator = (*BearerAuthenticator)(nil)

// Name 实现 authn.Authenticator。
func (a *BearerAuthenticator) Name() string { return "identity-bearer" }

// Authenticate 实现 authn.Authenticator。
func (a *BearerAuthenticator) Authenticate(ctx context.Context, c authn.Credentials) (*principal.Principal, error) {
	token := c.BearerToken()
	if token == "" {
		return nil, nil // 非 Bearer 凭证：交给链上下一个认证器
	}
	user, err := a.uc.Authenticate(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	return &principal.Principal{
		ID:    user.ID,
		Kind:  principal.KindUser,
		Attrs: map[string]string{"email": user.Email, "display_name": user.DisplayName},
	}, nil
}
