package infra

type Infra struct {
	Kc    *KafkaConsumer
	Rpc   *RPCConn
	Minio *Minio
	Redis *Redis
}

func NewInfra(kc *KafkaConsumer, rpc *RPCConn, minio *Minio, redis *Redis) *Infra {
	return &Infra{
		Kc:    kc,
		Rpc:   rpc,
		Minio: minio,
		Redis: redis,
	}
}
