package nertivia

type buttonPayload struct {
	Message string `json:"message"`
	TempID string `json:"tempID"`
	Buttons []*button `json:"buttons"`
}

func (bp *buttonPayload) add(id string, name string) {
	btn := new(button)
	btn.ID = id
	btn.Name = name
	bp.Buttons = append(bp.Buttons, btn)
}

type button struct {
	ID string `json:"id"`
	Name string `json:"name"`
}
