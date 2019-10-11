package util

// Environment variables that represent the executing environment
type Environment struct {
	Paths       map[string]string
	Urls        map[string]string
	ServiceList map[string]interface{}
}

// Path returns a path from Environment.Paths or "" if not found
func (s Environment) Path(name string) string {
	val, _ := s.Paths[name]

	return val
}

// URL returns a url from Environment.Urls or "" if not found
func (s Environment) URL(name string) string {
	val, _ := s.Urls[name]

	return val
}
