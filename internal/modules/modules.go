// Package modules 登记应用已安装的 lynx-start 模块。
// 本文件由 lyx add 维护（安装清单见仓库根 lynx-start.json）；锚点注释勿删。
package modules

import (
	"github.com/lynx-go/lynx-start/kernel"

	demo "github.com/lynx-go/lynx-start-demo/internal/api/demo"
// lyx:module-import（lyx add 注入锚点，勿删）
)

// Installed 返回全部已装模块（静态描述 + 动态装配）。
func Installed() []kernel.InstalledModule {
	return []kernel.InstalledModule{
		{
			Module:   demo.Module(),
			Assemble: demo.Assemble,
		},
		// lyx:module-entry（lyx add 注入锚点，勿删）
	}
}
