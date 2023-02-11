package context

import (
	"context"

	"github.com/si-bas/go-rest-boilerplate/pkg/logger/tag"
)

type key string

const (
	// LoggingTagKey is reserved name in the context for auto logging
	LoggingTagKey key = "logging_tag"
)

// AddLoggingTag func to add context logging tag
func AddLoggingTag(ctx context.Context, tagsToAdd ...tag.Tag) context.Context {
	allTags := ctx.Value(LoggingTagKey)
	if len(tagsToAdd) == 0 {
		return ctx
	}
	if contextTags, ok := allTags.(map[string]string); ok && contextTags != nil {
		return context.WithValue(ctx, LoggingTagKey, mergeTags(contextTags, tagsToAdd...))
	}
	return context.WithValue(ctx, LoggingTagKey, mergeTags(nil, tagsToAdd...))
}

func mergeTags(contextTags map[string]string, tagsToAdd ...tag.Tag) map[string]string {
	if contextTags == nil {
		contextTags = make(map[string]string)
	}
	for _, tag := range tagsToAdd {
		contextTags[tag.Key] = tag.Value
	}
	return contextTags
}

// InjectRequestID to inject msgID to logging tag
func InjectRequestID(ctx context.Context, msgID string) context.Context {
	return AddLoggingTag(ctx, tag.Tag{
		Key:   tag.RequestIDKey,
		Value: msgID,
	})
}

// GetAllLoggingTagInTagStr to get all tag str from logging tag
func GetAllLoggingTagInTagStr(ctx context.Context) []tag.Tag {
	allTags := ctx.Value(LoggingTagKey)
	contextTags, ok := allTags.(map[string]string)
	if !ok || contextTags == nil {
		return nil
	}

	var tags []tag.Tag
	for k, v := range contextTags {
		tags = append(tags, tag.Tag{
			Key:   k,
			Value: v,
		})
	}
	return tags
}

// GetTagValue to get a value from specific key
func GetTagValue(ctx context.Context, tagKey string) string {
	allTags := ctx.Value(LoggingTagKey)
	if contextTags, ok := allTags.(map[string]string); ok && contextTags != nil {
		return contextTags[tagKey]
	}
	return ""
}
