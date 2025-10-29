package main

import (
	"WoopsBBS/api"
	"WoopsBBS/global/Database/mysql"
)

func main() {
	mysql.GormInit()
	api.GinInit()

}
