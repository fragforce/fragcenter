package plugin

type TrackedPlugins struct {
	plugMakers []NewCellFunc
}

var (
	tracked *TrackedPlugins
)

func NewTrackedPlugins() *TrackedPlugins {
	return &TrackedPlugins{
		plugMakers: make([]NewCellFunc, 0),
	}
}

func init() {
	tracked = NewTrackedPlugins()
}

func (t *TrackedPlugins) Register(c NewCellFunc) {
	t.plugMakers = append(t.plugMakers, c)
}

func Register(c NewCellFunc) {
	tracked.Register(c)
}
