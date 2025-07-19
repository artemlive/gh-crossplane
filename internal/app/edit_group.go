package app

type editGroupModel struct {
	tabs       []string
	activeTab  int
	tabContent []string // Placeholder for content in each tab
}

func NewEditGroupModel() editGroupModel {
	return editGroupModel{
		tabs:       []string{"General", "Features", "Protections", "Permissions"},
		activeTab:  0,
		tabContent: make([]string, 4), // placeholder content per tab
	}
}
