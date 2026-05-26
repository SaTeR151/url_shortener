package inmemory

type (
	code string

	defaultURL string
)

func (u code) String() string {
	return string(u)
}

func (u defaultURL) String() string {
	return string(u)
}

func inCode(s string) code {
	return code(s)
}

func inDefaultURL(s string) defaultURL {
	return defaultURL(s)
}
