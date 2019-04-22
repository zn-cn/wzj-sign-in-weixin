package constant

const (

	/****************************************** table name ****************************************/
	TableUser = "user"

	RedisUserCurCoordinate    = "user:%s:coordinate" // %s -> openid, value: %f,%f
	RedisUserCoordinateFormat = "%f,%f"
	// %s -> openid, value: 0 -> 默认, 1-> 设置openid并启动任务 2->添加标签，经纬度 3->设置当前坐标为标签
	RedisUserStatus = "user:%s:status"

	RedisUserTask            = "user:task:%s"             // %s -> openid, value: 微助教openid
	RedisUserDisCourseNotice = "user:course:notice:%s:%d" // %s -> openid, %d -> course ID
)
