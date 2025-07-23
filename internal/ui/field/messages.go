package field

type FieldDoneMsg struct{}
type FieldDoneUpMsg struct{}
type FieldDoneDownMsg struct{}

// FieldOpenMsg used for the complex fields that can open a modal or a new screen
type FieldOpenMsg struct {
	Label string
	Value any // could be *Repository, *Permission, etc.
}
