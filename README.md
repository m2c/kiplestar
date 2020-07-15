# kiplestar golang base framework

//use GetKipleServerInstance() get kiplestar instance
server := kiplestar.GetKipleServerInstance()

//use server.Default() create Default server
server.Default()

//then call server.StartServer() start server and you can use kiplestar.Redis_service,kiplestar.Mysql_service start redis and mysql
server.StartServer()

//simple
func main() {
	// log init
	server := kiplestar.GetKipleServerInstance()
	//http
	server.Default()
	router.RegisterGlobalModel(server.App().GetIrisApp());
	err := server.StartServer()
	if err != nil {
		fmt.Println(err.Error())
	}
}
