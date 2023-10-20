package messages

import (
	"time"
)

type Options struct {
	Version    string
	Bits       string
	Commit     string
	Modified   string
	Pid        string
	Port       string
	ConfigFile bool
}

func RpStartup(options Options) string {
	var builder CustomBuilder

	prefix := func() string {
		return options.Pid + ": " + time.Now().Format("02 Jan 2006 15:04:05.000")
	}

	configLine := func() string {
		if options.ConfigFile {
			return prefix() + " # Warning: no config file specified, using the default config. In order to specify a config file use redis-server"
		} else {
			return prefix() + " # Warning: no config file specified, using the default config. In order to specify a config file use redis-server"
		}
	}

	builder.WriteString(prefix() + " * oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo")
	builder.WriteFString("%s * Redis version=%s, bits=%s, commit=%s, modified=%s, options.Pid=%s, just started", prefix(), options.Version, options.Bits, options.Commit, options.Modified, options.Pid)
	builder.WriteString(configLine())

	builder.WriteString("                _._")
	builder.WriteString("           _.-``__ ''-._")
	builder.WriteFString("      _.-``    `.  `_.  ''-._           Redis %s %s %s bit", options.Version, options.Bits, options.Commit)
	builder.WriteString("  .-`` .-```.  ```\\/    _.,_ ''-._")
	builder.WriteString(" (    '      ,       .-`  | `,    )     Running in standalone mode")
	builder.WriteFString(" |`-._`-...-` __...-.``-._|'` _.-'|     Port: %s", options.Port)
	builder.WriteFString(" |    `-._   `._    /     _.-'    |     PID: %s", options.Pid)
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
