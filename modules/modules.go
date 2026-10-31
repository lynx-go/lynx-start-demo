// Package modules 登记应用已安装的 lynx-mods 模块。
// 本文件由 lynx-mods add 维护（安装清单见仓库根 lynx-mods.json）；锚点注释勿删。
package modules

import (
	"github.com/lynx-go/lynx-start"
	access "github.com/lynx-go/lynx-start-demo/modules/access"
	audit "github.com/lynx-go/lynx-start-demo/modules/audit"
	demo "github.com/lynx-go/lynx-start-demo/modules/demo"
	identity "github.com/lynx-go/lynx-start-demo/modules/identity"
	notification "github.com/lynx-go/lynx-start-demo/modules/notification"
	organization "github.com/lynx-go/lynx-start-demo/modules/organization"
	session "github.com/lynx-go/lynx-start-demo/modules/session"
	task "github.com/lynx-go/lynx-start-demo/modules/task"
	// lynx-mods:module-import（lynx-mods add 注入锚点，勿删）
)

// Installed 返回全部已装模块（静态描述 + 动态装配）。
func Installed() []kernel.InstalledModule {
	return []kernel.InstalledModule{
		{
			Module:   demo.Module(),
			Assemble: demo.Assemble,
		},
		{
			Module:   identity.Module(),
			Assemble: identity.Assemble,
		},
		{
			Module:   session.Module(),
			Assemble: session.Assemble,
		},
		{
			Module:   access.Module(),
			Assemble: access.Assemble,
		},
		{
			Module:   audit.Module(),
			Assemble: audit.Assemble,
		},
		{
			Module:   organization.Module(),
			Assemble: organization.Assemble,
		},
		{
			Module:   notification.Module(),
			Assemble: notification.Assemble,
		},
		{
			Module:   task.Module(),
			Assemble: task.Assemble,
		},
		// lynx-mods:module-entry（lynx-mods add 注入锚点，勿删）
	}
}
