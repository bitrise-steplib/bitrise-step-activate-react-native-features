package step

type CacheFeature interface {
	ActivateArgs() []string
	ServiceArgs() []string // nil if no background service
}
