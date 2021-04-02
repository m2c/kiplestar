package kiplestar

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
	irisv12 "github.com/kataras/iris/v12"
	"github.com/m2c/kiplestar/commons"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/commons/utils"
	"github.com/m2c/kiplestar/config"
	"github.com/m2c/kiplestar/iris"
	"github.com/m2c/kiplestar/kafka"
	"github.com/m2c/kiplestar/kipledb"
	"github.com/m2c/kiplestar/redis"
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
		// read config file
		config.Init()
		//Automatically loads the configured services
		//except mysql redis
		kipleInstance.initService()
	})
	return kipleInstance
}

func GetKipleServerCustomInstance(path string) *kipleSever {
	once.Do(func() {
		kipleInstance = new(kipleSever)
		config.InitCustomPath(path)
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

func (slf *kipleSever) RegisterErrorCodeAndMsg(arr map[commons.ResponseCode]string) {
	if len(arr) == 0 {
		return
	}
	for k, v := range arr {
		commons.CodeMsg[k] = v
	}
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
			slog.Log.Sync()
		}
	}()
	err := slf.app.Start(params...)
	if err != nil {
		panic(err)
	}
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
func (slf *kipleSever) LoadCustomizeConfig(slfConfig interface{}) {
	err := config.LoadCustomizeConfig(slfConfig)
	if err != nil {
		panic(err)
	}
}

//need call this function after Option, if Dependent service is not started return panic.
func (slf *kipleSever) StartServer(opt ...Server_Option) {
	var err error
	for _, v := range opt {
		switch v {
		case Mysql_service:
			slf.db = make([]kipledb.KipleDB, len(config.Configs.DataBase))
			for i, v := range config.Configs.DataBase {
				err = slf.db[i].StartDb(v)
				if err != nil {
					panic(err)
				}
			}
		case Redis_service:
			slf.redis = make([]redis.Redis, len(config.Configs.Redis))
			for i, v := range config.Configs.Redis {
				err = slf.redis[i].StartRedis(v)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (slf *kipleSever) KafkaService(ctx context.Context, topic string, callBackChan chan []byte) {
	slf.kafka.KafkaReceiver(ctx, topic, callBackChan)
}
