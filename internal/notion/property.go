package notion

import "time"

type Property struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Checkbox     bool         `json:"checkbox,omitempty"`
	CreatedBy    *User        `json:"created_by,omitempty"`
	CreatedTime  *time.Time   `json:"created_time,omitempty"`
	Date         *Date        `json:"date,omitempty"`
	Email        *string      `json:"email,omitempty"`
	Files        []FileObject `json:"files,omitempty"`
	Formula      *Formula     `json:"formula,omitempty"`
	Icon         any          `json:"icon,omitempty"`
	LastEditedBy *User        `json:"last_edited_by,omitempty"`
	Title        []*RichText  `json:"title,omitempty"`
	RichText     []*RichText  `json:"rich_text,omitempty"`
	Url          *string      `json:"url,omitempty"`
	Select       *Select      `json:"select,omitempty"`
}

type Date struct {
	Start    string  `json:"start"`
	End      *string `json:"end,omitempty"`
	TimeZone *string `json:"time_zone,omitempty"`
}

type Formula struct {
	Type   string   `json:"type"`
	Bool   *bool    `json:"boolean,omitempty"`
	Date   *Date    `json:"date,omitempty"`
	Number *float64 `json:"number,omitempty"`
	String *string  `json:"string,omitempty"`
}

type RichText struct {
	Type        string      `json:"type"`
	Text        *Text       `json:"text,omitempty"`
	PlainText   string      `json:"plain_text,omitempty"`
	Href        *string     `json:"href,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type Text struct {
	Content string  `json:"content"`
	Link    *string `json:"link,omitempty"`
}

type Select struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
