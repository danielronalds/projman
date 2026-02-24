package worktree

type mockSessionLauncher struct {
	returnErr  error
	calledName string
	calledDir  string
}

func (m *mockSessionLauncher) LaunchSession(name, dir string) error {
	m.calledName = name
	m.calledDir = dir
	return m.returnErr
}

type mockSelecter struct {
	returnSelected string
	returnErr      error
	calledOptions  []string
}

func (m *mockSelecter) Select(options []string) (string, error) {
	m.calledOptions = options
	return m.returnSelected, m.returnErr
}
