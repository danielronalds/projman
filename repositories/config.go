package repositories

// TODO: make user configurable
type ConfigRepository struct {}

func NewConfigRepository() ConfigRepository {
	return ConfigRepository{}
}

func (r ConfigRepository) Theme() string {
	return "bw"
}

func (r ConfigRepository) Layout() string {
	return "reverse"
}

func (r ConfigRepository) ProjectDirs() []string {
	return []string {
		"/home/danielr/Projects/",
		"/home/danielr/.config/",
	}
}
