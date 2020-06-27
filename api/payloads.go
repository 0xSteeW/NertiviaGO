package nertivia

type buttonPayload struct {
	Buttons []button `json:"buttons"`
}

func (bp buttonPayload) add(id string, name string) {
	btn := new(button)
	btn.id = id
	btn.name = name
	bp.Buttons = append(bp.Buttons, *btn)
}

type button struct {
	id string
	name string
}
