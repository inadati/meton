package env

type Zookeeper struct {
	MYID    int
	SERVERS string
}
type MesosMaster struct {
	MESOS_HOSTNAME string
	MESOS_IP       string
	MESOS_ZK       string
}

type Marathon struct {
	MARATHON_HOSTNAME      string
	MARATHON_HTTPS_ADDRESS string
	MARATHON_HTTP_ADDRESS  string
	MARATHON_MASTER        string
	MARATHON_ZK            string
}

type MesosSlave struct {
	MESOS_HOSTNAME string
	MESOS_IP       string
	MESOS_MASTER   string
}

type Chronos struct {
	CHRONOS_MASTER   string
	CHRONOS_ZK_HOSTS string
}
