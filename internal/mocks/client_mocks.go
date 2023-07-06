package mocks

// MockDBClient is a mock implementation of the SurrealDBClient interface.
type MockSurrealDBClient struct {
	CloseFunc  func()
	CreateFunc func(thing string, data interface{}) (interface{}, error)
	QueryFunc  func(sql string, vars interface{}) (interface{}, error)
	SigninFunc func(vars interface{}) (interface{}, error)
	UseFunc    func(namespace string, database string) (interface{}, error)
}

func (m *MockSurrealDBClient) Close() {
	m.CloseFunc()
}

func (m *MockSurrealDBClient) Query(sql string, vars interface{}) (interface{}, error) {
	return m.QueryFunc(sql, vars)
}

func (m *MockSurrealDBClient) Signin(vars interface{}) (interface{}, error) {
	return m.SigninFunc(vars)
}

func (m *MockSurrealDBClient) Use(namespace string, database string) (interface{}, error) {
	return m.UseFunc(namespace, database)
}

func (m *MockSurrealDBClient) Create(thing string, data interface{}) (interface{}, error) {
	return m.CreateFunc(thing, data)
}
