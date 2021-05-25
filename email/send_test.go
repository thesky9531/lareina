package email

import "testing"

func TestNewMailEntity(t *testing.T) {
	m := NewMailEntity()
	m.SetMailEntity("", "", "", "", "", "", "")
	err := m.SendToMail()
	if err != nil {

	}
}
