// Code generated by protoc-gen-goext. DO NOT EDIT.

package vision

func (m *TextAnnotation) SetPages(v []*Page) {
	m.Pages = v
}

func (m *Page) SetWidth(v int64) {
	m.Width = v
}

func (m *Page) SetHeight(v int64) {
	m.Height = v
}

func (m *Page) SetBlocks(v []*Block) {
	m.Blocks = v
}

func (m *Page) SetEntities(v []*Entity) {
	m.Entities = v
}

func (m *Entity) SetName(v string) {
	m.Name = v
}

func (m *Entity) SetText(v string) {
	m.Text = v
}

func (m *Block) SetBoundingBox(v *Polygon) {
	m.BoundingBox = v
}

func (m *Block) SetLines(v []*Line) {
	m.Lines = v
}

func (m *Line) SetBoundingBox(v *Polygon) {
	m.BoundingBox = v
}

func (m *Line) SetWords(v []*Word) {
	m.Words = v
}

func (m *Line) SetConfidence(v float64) {
	m.Confidence = v
}

func (m *Word) SetBoundingBox(v *Polygon) {
	m.BoundingBox = v
}

func (m *Word) SetText(v string) {
	m.Text = v
}

func (m *Word) SetConfidence(v float64) {
	m.Confidence = v
}

func (m *Word) SetLanguages(v []*Word_DetectedLanguage) {
	m.Languages = v
}

func (m *Word) SetEntityIndex(v int64) {
	m.EntityIndex = v
}

func (m *Word_DetectedLanguage) SetLanguageCode(v string) {
	m.LanguageCode = v
}

func (m *Word_DetectedLanguage) SetConfidence(v float64) {
	m.Confidence = v
}
