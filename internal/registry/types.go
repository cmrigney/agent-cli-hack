package registry

type Registry struct {
	Agents map[string]bool `yaml:"agents"`
}
