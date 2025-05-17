package context

type ContextKey string

const (
	UserID  ContextKey = "user_id"
	Logging ContextKey = "logging"
)
