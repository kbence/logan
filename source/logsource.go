package source

// LogSource provides an interface to access logs
type LogSource interface {
	GetCategories() []string
}
