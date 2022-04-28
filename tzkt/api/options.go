package api

// Option -
type Option func(*API)

// WithAuth -
func WithAuth(user, privateKey string) Option {
	return func(api *API) {
		api.user = user
		api.privateKey = privateKey
	}
}
