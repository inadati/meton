package container

type SingleContainerOperation interface {
	up() error
}

type AllContainerOperation interface {
	down() error
}

func Up(fork SingleContainerOperation) error {
	return fork.up()
}

func DownAll(fork AllContainerOperation) error {
	return fork.down()
}
