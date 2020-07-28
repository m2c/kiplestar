package kiplestar

import (
	"context"
	redisv8 "github.com/go-redis/redis/v8"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/config"
	"github.com/m2c/kiplestar/iris"
	"github.com/m2c/kiplestar/kafka"
	"github.com/m2c/kiplestar/kipledb"
	"github.com/m2c/kiplestar/redis"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//we need create the single object but thread safe
var kipleInstance *kipleSever

var once sync.Once

type kipleSever struct {
	app   iris.App
	redis redis.Redis
	db    []kipledb.KipleDB
	kafka kafka.Kafka
}
type Server_Option int

const (
	Mysql_service = iota + 1
	Redis_service
)

//create the single object
func GetKipleServerInstance() *kipleSever {
	once.Do(func() {
		kipleInstance = new(kipleSever)
	})
	return kipleInstance
}
func (slf *kipleSever) Default() {
	slf.app.Default()

}
func (slf *kipleSever) WaitClose() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			// kill -SIGINT XXXX 或 Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			os.Kill,
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXXD
			//^
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			slog.Infof("wait for close server")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for _, db := range slf.db {
				db.StopDb()
			}
			slf.app.GetIrisApp().Shutdown(ctx)

		}
	}()
	slf.app.Start()
}
func (slf *kipleSever) New() {
	slf.app.New()
}

//return app
func (slf *kipleSever) App() *iris.App {
	return &slf.app
}
func (slf *kipleSever) DB(name string) *kipledb.KipleDB {
	for _, v := range slf.db {
		if v.Name() == name {
			return &v
		}
	}
	return nil
}
func (slf *kipleSever) Redis() *redisv8.Client {
	return slf.redis.Redis()
}
func (slf *kipleSever) LoadCustomizeConfig(slfConfig interface{}) error {
	return config.LoadCustomizeConfig(slfConfig)
}

//need call this function after Option
func (slf *kipleSever) StartServer(opt ...Server_Option) (err error) {
	for _, v := range opt {
		switch v {
		case Mysql_service:
			slf.db = make([]kipledb.KipleDB, len(config.Configs.DataBase))
			for i, v := range config.Configs.DataBase {
				err = slf.db[i].StartDb(v)
			}
			if err != nil {
				return err
			}
		case Redis_service:
			err = slf.redis.StartRedis()
		}
		if err != nil {
			return err
		}
	}

	return
}

func (slf *kipleSever) KafkaService(ctx context.Context, topic string, callBackChan chan []byte) {
	slf.kafka.KafkaReceiver(ctx, topic, callBackChan)
}
