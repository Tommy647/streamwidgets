package widgets

// Widget watcher
type Widget struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// New widget
func New() *Widget { return &Widget{} }
