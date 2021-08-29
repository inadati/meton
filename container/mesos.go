package container

type mesos struct {
	Master MesosMasterRecipe
	Slave  MesosSlaveRecipe
}

var Mesos = &mesos{
	Master: *master,
	Slave:  *slave,
}
