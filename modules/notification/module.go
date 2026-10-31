// lynx-mods:bridge —— 由 add 生成（remove 删除）。转发装配入口：
// modules/<name>/internal 受 Go 可见性保护，应用唯一入口是本包。
// 改模块装配请改 internal/api；本文件属注入物，不参与 update 合并。
package notification

import (
	"github.com/lynx-go/lynx-start"

	api "github.com/lynx-go/lynx-start-demo/modules/notification/internal/api"
)

// Module 转发模块静态装配描述。
func Module() kernel.Module { return api.Module() }

// Assemble 转发模块动态装配。
func Assemble(deps kernel.Deps) (kernel.Assembly, error) { return api.Assemble(deps) }
