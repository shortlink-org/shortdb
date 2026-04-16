package repl

import tea "charm.land/bubbletea/v2"

// isSelectCLIKeyMsg matches f1 (workspace tab).
func isSelectCLIKeyMsg(msg tea.KeyPressMsg) bool {
	return msg.String() == "f1"
}

// isSelectObservableKeyMsg matches f2 (workspace tab).
func isSelectObservableKeyMsg(msg tea.KeyPressMsg) bool {
	return msg.String() == "f2"
}

func keyScrollsTranscript(msg tea.KeyPressMsg) bool {
	switch msg.String() {
	case "pgup", "pgdown":
		return true
	default:
		return false
	}
}
