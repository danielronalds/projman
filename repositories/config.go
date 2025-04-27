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
		"/data/data/com.termux/files/home/Projects/",
		"/data/data/com.termux/files/home/.config/",
	}
}
