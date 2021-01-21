package kiplestar

import (
	"context"
	redisv8 "github.com/go-redis/redis/v8"
	irisv12 "github.com/kataras/iris/v12"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/commons/utils"
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
	redis []redis.Redis
	db    []kipledb.KipleDB
	kafka kafka.Kafka
	Oss   utils.OSSClient
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
		config.Init()
		//Automatically loads the configured services
		//except mysql redis
		kipleInstance.initService()
	})
	return kipleInstance
}
func (slf *kipleSever) Default() {
	slf.app.Default()

}

func GetOss() utils.OSSClient {
	return kipleInstance.Oss
}

func (slf *kipleSever) initService() {
	if config.Configs.Oss.OssBucket != "" {
		slf.Oss = utils.OSSClientInstance(config.Configs.Oss.OssBucket, config.Configs.Oss.AccessKeyID, config.Configs.Oss.AccessKeySecret, config.Configs.Oss.OssEndPoint)
	}
}

func (slf *kipleSever) RegisterController(f func(app *irisv12.Application)) {
	f(slf.app.GetIrisApp())
}

func (slf *kipleSever) WaitClose(params ...irisv12.Configurator) {
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
	slf.app.Start(params...)
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
func (slf *kipleSever) Redis(name string) *redisv8.Client {
	for _, v := range slf.redis {
		if v.Name() == name {
			return v.Redis()
		}
	}
	return nil
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
			slf.redis = make([]redis.Redis, len(config.Configs.Redis))
			for i, v := range config.Configs.Redis {
				err = slf.redis[i].StartRedis(v)
			}
			if err != nil {
				return err
			}
			//err = slf.redis.StartRedis()
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
