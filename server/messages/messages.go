package messages

import (
	"time"

	"github.com/piratey7007/rediss/server/types"
)

func RpStartup(options types.ConnectionOptions) string {
	var builder CustomBuilder

	prefix := func() string {
		return "RedissServer: " + time.Now().Format("02 Jan 2006 15:04:05.000")
	}

	builder.WriteString(prefix() + " * oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo")
	builder.WriteFString("%s * Rediss, just started", prefix())

	builder.WriteString("                _._")
	builder.WriteString("           _.-``__ ''-._")
	builder.WriteString("      _.-``    `.  `_.  ''-._           Rediss Clone")
	builder.WriteString("  .-`` .-```.  ```\\/    _.,_ ''-._")
	builder.WriteString(" (    '      ,       .-`  | `,    )     Running in standalone mode")
	builder.WriteFString(" |`-._`-...-` __...-.``-._|'` _.-'|     Port: %s", options.Port)
	builder.WriteString(" |    `-._   `._    /     _.-'    |")
	builder.WriteString("  `-._    `-._  `-./  _.-'    _.-'")
	builder.WriteString(" |`-._`-._    `-.__.-'    _.-'_.-'|")
	builder.WriteString(" |    `-._`-._        _.-'_.-'    |           https://redis.io")
	builder.WriteString("  `-._    `-._`-.__.-'_.-'    _.-'")
	builder.WriteString(" |`-._`-._    `-.__.-'    _.-'_.-'|")
	builder.WriteString(" |    `-._`-._        _.-'_.-'    |")
	builder.WriteString("  `-._    `-._`-.__.-'_.-'    _.-'")
	builder.WriteString("      `-._    `-.__.-'    _.-'")
	builder.WriteString("          `-._        _.-'")
	builder.WriteString("              `-.__.-'")

	builder.WriteString(prefix() + " * Server initialized")
	builder.WriteString(prefix() + " * Ready to accept connections")

	return builder.String()
}

func RpClose() string {
	var builder CustomBuilder

	builder.WriteString("59377:M 20 Oct 2023 09:39:35.780 # User requested shutdown...")
	builder.WriteString("59377:M 20 Oct 2023 09:39:35.787 * DB saved on disk")
	builder.WriteString("59377:M 20 Oct 2023 09:39:35.788 # Redis is now ready to exit, bye bye...")

	return builder.String()
}
