package container

type mesos struct {
	Master mesosMaster
	Slave  mesosSlave
}

var Mesos = &mesos{
	Master: *master,
	Slave:  *slave,
}
