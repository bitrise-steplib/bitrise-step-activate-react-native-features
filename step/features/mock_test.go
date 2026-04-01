package features_test

type MockEnvRepo struct {
	vars map[string]string
}

func NewMockEnvRepo() *MockEnvRepo {
	return &MockEnvRepo{
		vars: make(map[string]string),
	}
}

func (m *MockEnvRepo) Get(key string) string {
	return m.vars[key]
}

func (m *MockEnvRepo) Set(key, value string) error {
	m.vars[key] = value
	return nil
}

func (m *MockEnvRepo) Unset(key string) error {
	delete(m.vars, key)
	return nil
}

func (m *MockEnvRepo) List() []string {
	var list []string
	for k, v := range m.vars {
		list = append(list, k+"="+v)
	}
	return list
}
