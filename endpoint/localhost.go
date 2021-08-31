package endpoint

type LocalHostRecipe struct{}

var LocalHost = &LocalHostRecipe{}

func (r *LocalHostRecipe) Marathon() string {
	return "127.0.0.1:8080"
}

func (r *LocalHostRecipe) MesosMaster() string {
	return "127.0.0.1:5050"
}

func (r *LocalHostRecipe) Chronos() string {
	return "127.0.0.1:4400"
}
