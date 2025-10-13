package types

type Resume struct {
	ID       int64               `json:"id" db:"id"`
	FilePath string              `json:"file_path" db:"file_path"`
	Raw      map[string]string   `json:"raw" db:"raw"`
	Sections map[string]Sections `json:"sections" db:"sections"`
	Metadata map[string]string   `json:"metadata" db:"metadata"`
}

type Sections struct {
	Type    SectionType `json:"type"`
	Content interface{} `json:"content"`
}

type SectionType string

const (
	ContactSection  SectionType = "contact"
	TimelineSection SectionType = "timeline"
	ListSection     SectionType = "list"
	FreeformSection SectionType = "freeform"
)

type ContactContent struct {
	Name     string            `json:"name"`
	Email    []string          `json:"email"`
	Number   []string          `json:"number"`
	Location string            `json:"location"`
	Social   map[string]string `json:"social"`
}

type TimelineContent struct {
	Entries []TimelineEntry `json:"entries"`
}

type TimelineEntry struct {
	Organization string            `json:"organization"`
	Location     string            `json:"location"`
	Title        string            `json:"title"`
	StartDate    string            `json:"start_date"`
	EndDate      string            `json:"end_date"`
	Details      []string          `json:"details"`
	Metadata     map[string]string `json:"metadata"`
}

type ListContent struct {
	Categories []ListCategory `json:"categories"`
}

type ListCategory struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

type FreeformContent struct {
	Entries []FreeformEntry `json:"entries"`
}

type FreeformEntry struct {
	Heading string   `json:"heading"`
	Content []string `json:"content"`
}

type ParsedResume struct {
	Raw      map[string]string      `json:"raw"`
	Sections map[string]interface{} `json:"sections"`
	Metadata map[string]string      `json:"metadata"`
}
