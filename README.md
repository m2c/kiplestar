#### Introduction

​	This project is the package of the basic project framework of the company’s go iris, including configuration file reading, database transaction package, graceful shutdown, basic return parameter package, http package and many other new features. Currently, business systems are recommended to use this System. At the beginning of the project, in order to unify the framework and technology stack of the company, the project will inevitably have bugs or issues. If there are related problems, welcome to point out. I believe that on top of our continuous optimization and iterative version, he will become more The better.

#### History

| Version | Date       | type   | description |
| :------ | :--------- | :----- | ----------- |
| 1.0.0   | 2020-07-20 | common |             |

#### Architecture

Software architecture description.

#### Guide

1. go get github.com/m2c/kiplestar

#### Instructions

1. Conditions of use: go SDK version >=14
2. application.yaml is a public file
3. application-dev.yaml represents the dev environment, application-prod.yaml represents the production environment, application-xxx.yaml represents the relationship between the xxx environment and the user's custom environment, and there is an option profile in the application-dev.yaml configuration : dev means the configuration file that activates the xxx environment

#### Get Start


```go
func main() {

	server := kiplestar.GetKipleServerInstance()
	server.Default()

	err := server.StartServer(kiplestar.Mysql_service)
	if err != nil {
		panic(err.Error())
	}

	router.RegisterGlobalModel(server.App().GetIrisApp())

	server.WaitClose()
}
```

#### Core function Example

1.

#### Project structure

```sybase
├─commons //common package
│  ├─error//error wrap
│  ├─log  // log 
│  ├─time // time format
│  └─utils //often utils tools
├─config   //configration file parse
├─database //gorm datasource
├─iris //iris
├─kafka //kafka client 
├─kipledb //dao base struct
├─middleware// often middlewares
└─redis  //redis

```

#### Contact :

mark：mark.jiang@greenpacket.com.cn

seven: seven.zhang@greenpacket.com.cn

