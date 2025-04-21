package infra

type Infra struct {
	Kc    *KafkaConsumer
	Rpc   *RPCConn
	Minio *Minio
}

func NewInfra(kc *KafkaConsumer, rpc *RPCConn, minio *Minio) *Infra {
	return &Infra{
		Kc:    kc,
		Rpc:   rpc,
		Minio: minio,
	}
}
