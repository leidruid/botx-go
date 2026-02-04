package models

import "github.com/leidruid/botx-go/optional"

// ButtonTextAlign defines text alignment.
type ButtonTextAlign string

const (
	AlignLeft   ButtonTextAlign = "left"
	AlignCenter ButtonTextAlign = "center"
	AlignRight  ButtonTextAlign = "right"
)

// Button represents UI button.
type Button struct {
	Label           string
	Command         optional.Optional[string]
	Data            map[string]any
	TextColor       optional.Optional[string]
	BackgroundColor optional.Optional[string]
	Align           ButtonTextAlign
	Silent          bool
	WidthRatio      optional.Optional[int]
	Alert           optional.Optional[string]
	ProcessOnClient optional.Optional[bool]
	Link            optional.Optional[string]
}

// ButtonRow is a list of buttons.
type ButtonRow []Button

// Markup represents buttons layout.
type Markup struct {
	Rows []ButtonRow
}

func NewMarkup() *Markup {
	return &Markup{}
}

func (m *Markup) AddButton(btn Button, newRow bool) {
	if newRow || len(m.Rows) == 0 {
		m.Rows = append(m.Rows, ButtonRow{btn})
		return
	}
	m.Rows[len(m.Rows)-1] = append(m.Rows[len(m.Rows)-1], btn)
}

// API models for BotX

type apiButtonOptions struct {
	Silent          optional.Optional[bool]   `json:"silent,omitempty"`
	FontColor       optional.Optional[string] `json:"font_color,omitempty"`
	BackgroundColor optional.Optional[string] `json:"background_color,omitempty"`
	Align           optional.Optional[string] `json:"align,omitempty"`
	HSize           optional.Optional[int]    `json:"h_size,omitempty"`
	ShowAlert       optional.Optional[bool]   `json:"show_alert,omitempty"`
	AlertText       optional.Optional[string] `json:"alert_text,omitempty"`
	Handler         optional.Optional[string] `json:"handler,omitempty"`
	Link            optional.Optional[string] `json:"link,omitempty"`
}

type apiButton struct {
	Command optional.Optional[string] `json:"command,omitempty"`
	Label   string                    `json:"label"`
	Data    map[string]any            `json:"data"`
	Opts    apiButtonOptions          `json:"opts"`
}

type APIMarkup struct {
	Rows [][]apiButton `json:"__root__"`
}

// ToAPI converts markup to BotX API format.
func (m *Markup) ToAPI() APIMarkup {
	rows := make([][]apiButton, 0, len(m.Rows))
	for _, row := range m.Rows {
		apiRow := make([]apiButton, 0, len(row))
		for _, btn := range row {
			showAlert := optional.None[bool]()
			if btn.Alert.Set {
				showAlert = optional.Some(true)
			}
			handler := optional.None[string]()
			if btn.ProcessOnClient.Set && btn.ProcessOnClient.Value {
				handler = optional.Some("client")
			}
			if btn.Link.Set {
				handler = optional.Some("client")
			}

			apiRow = append(apiRow, apiButton{
				Command: btn.Command,
				Label:   btn.Label,
				Data:    btn.Data,
				Opts: apiButtonOptions{
					Silent:          optional.Some(btn.Silent),
					FontColor:       btn.TextColor,
					BackgroundColor: btn.BackgroundColor,
					Align:           optional.Some(string(btn.Align)),
					HSize:           btn.WidthRatio,
					ShowAlert:       showAlert,
					AlertText:       btn.Alert,
					Handler:         handler,
					Link:            btn.Link,
				},
			})
		}
		rows = append(rows, apiRow)
	}
	return APIMarkup{Rows: rows}
}
