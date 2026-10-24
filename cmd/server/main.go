// lynx-start-demo 服务入口：装配已安装的 lynx-start 模块并托管生命周期。
// 本文件由 lynx-start init 生成，后续基本不再变——模块增减只改 internal/modules。
package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx/contrib/schedule"
	"github.com/lynx-go/lynx/contrib/zap"

	grpcapiauthn "github.com/lynx-go/grpcapi/authn"
	grpcapiinterceptor "github.com/lynx-go/grpcapi/interceptor"
	grpcapipolicy "github.com/lynx-go/grpcapi/policy"
	grpcapiprotopolicy "github.com/lynx-go/grpcapi/policy/protopolicy"
	grpcapiserver "github.com/lynx-go/grpcapi/server"
	"github.com/lynx-go/lynx-start/kernel"

	"github.com/lynx-go/lynx-start-demo/internal/bootstrap"
	"github.com/lynx-go/lynx-start-demo/internal/config"
	"github.com/lynx-go/lynx-start-demo/internal/modules"
)

func main() {
	runner := lynx.NewRunner(setup,
		lynx.WithName("lynx-start-demo"),
		lynx.WithShutdownTimeout(30*time.Second),
		lynx.WithBindConfigFunc(config.BindConfig),
		lynx.WithLoggerProvider(func(app lynx.AppContext) (*slog.Logger, error) {
			return zap.NewLogger(app)
		}),
	)
	runner.Run()
}

func setup(app lynx.App) error {
	logger := app.Logger()

	cfg, err := config.Load(app.Config())
	if err != nil {
		return err
	}
	deps, cleanup, err := bootstrap.NewDeps(app, cfg, logger)
	if err != nil {
		return err
	}
	app.OnPostStop(cleanup)

	installed := modules.Installed()
	agg, err := kernel.AssembleAll(deps, installed)
	if err != nil {
		return err
	}

	// 策略集：全部已装模块的 proto 注解（fail-closed 启动断言由 grpcapi 承担）。
	set, err := buildPolicies(app.Context(), installed)
	if err != nil {
		return err
	}

	// 模块贡献的认证后增强器（access 的 roles/permissions claims 等）。
	authnOpts := grpcapiauthn.Options{Enrichers: agg.Enrichers}

	// 模块贡献的拦截器槽位（audit/ratelimit/usage 等）插入默认链。
	var chain *grpcapiinterceptor.Chain
	if len(agg.Stages) > 0 {
		chain, err = grpcapiinterceptor.AssembleWith(grpcapiinterceptor.DefaultChainOptions(), agg.Stages...)
		if err != nil {
			return err
		}
	}

	srv, err := grpcapiserver.New(app.Context(), grpcapiserver.Config{
		Policy:               set,
		Services:             agg.Registrations,
		ExtraRoutes:          agg.ExtraRoutes,
		Authenticators:       agg.Authenticators,
		AuthnOptions:         authnOpts,
		Chain:                chain,
		GRPCAddr:             cfg.Server.GRPCAddr,
		HTTPAddr:             cfg.Server.HTTPAddr,
		AllowNoAuthenticator: len(agg.Authenticators) == 0,
		Logger:               logger,
	})
	if err != nil {
		return err
	}
	app.Register(bootstrap.NewAPIService(srv, logger))

	if len(agg.Tasks) > 0 {
		scheduler, err := schedule.NewScheduler(agg.Tasks)
		if err != nil {
			return err
		}
		app.Register(scheduler)
	}
	app.Register(agg.Services...)
	return nil
}

func buildPolicies(ctx context.Context, installed []kernel.InstalledModule) (*grpcapipolicy.Set, error) {
	files := kernel.PolicyFiles(kernel.Descriptors(installed))
	if len(files) == 0 {
		return grpcapipolicy.NewSet()
	}
	return grpcapipolicy.Build(ctx, grpcapipolicy.FromFiles(files...), grpcapipolicy.Options{
		Provider: grpcapiprotopolicy.New(grpcapiprotopolicy.Options{}),
		Registry: grpcapipolicy.DefaultRegistry(),
	})
}
