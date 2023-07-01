package global

type Data struct {
	Token         string
	AdminServer   string
	AdminChannel  string
	InfoChannel   string
	WarnChannel   string
	ErrChannel    string
	UpdateChannel string
	APIEndpoint   string
	Debugging     bool
	ActiveSession bool
}
