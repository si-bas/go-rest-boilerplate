package shared

import "context"

func GetContextValueAsString(ctx context.Context, key string) string {
	val, ok := ctx.Value(key).(string)
	if ok {
		return val
	}

	return ""
}
