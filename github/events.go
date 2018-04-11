package github

// RepositoryEventReader retrieves the events of one GitHub repository.
type RepositoryEventReader struct {
	owner         string
	repository    string
	authenticator Authenticator
}

// NewRepositoryEventReader creates an event reader for the given
// owner and repository.
func NewRepositoryEventReader(owner, repository string) *RepositoryEventReader {
	return &RepositoryEventReader{
		owner:      owner,
		repository: repository,
	}
}

// SetAuthenticator allows the event reader to authenticate against
// GitHub if needed.
func (r *RepositoryEventReader) SetAuthenticator(authenticator Authenticator) {
	r.authenticator = authenticator
}
