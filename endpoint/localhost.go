package endpoint

type LocalHostRecipe struct{}

var LocalHost = &LocalHostRecipe{}

func (r *LocalHostRecipe) Marathon() string {
	return "0.0.0.0:8080"
}

func (r *LocalHostRecipe) MesosMaster() string {
	return "0.0.0.0:5050"
}

func (r *LocalHostRecipe) Chronos() string {
	return "0.0.0.0:4400"
}
