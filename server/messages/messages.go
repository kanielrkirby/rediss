package messages

import "strings"

func RedisPrompt() string {
	var builder strings.Builder

	builder.WriteString("                _._")
	builder.WriteString("           _.-``__ ''-._")
	builder.WriteString("      _.-``    `.  `_.  ''-._           Redis 7.2.1 (00000000/0) 64 bit")
	builder.WriteString("  .-`` .-```.  ```\\/    _.,_ ''-._")
	builder.WriteString(" (    '      ,       .-`  | `,    )     Running in standalone mode")
	builder.WriteString(" |`-._`-...-` __...-.``-._|'` _.-'|     Port: 6379")
	builder.WriteString(" |    `-._   `._    /     _.-'    |     PID: 93873")
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

	return builder.String()
}
