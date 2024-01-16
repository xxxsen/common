package versionmgr

type SQLActionStep struct {
	Name   string
	Action ActionFunc
}

type SQLAction struct {
	Steps []SQLActionStep
}

type VersionAction struct {
	OnInit      SQLAction //当前版本为初始化版本的时候的操作, 例如新安装, 必须可重入
	OnUpgrade   SQLAction //从上一个版本进行升级执行的操作, 必须可重入
	Version     uint64    //当前版本, 可以用日期进行标记, 例如20230506110004
	LastVersion uint64    //上一个版本, 只能从上一个版本升级, 当为空字符串的时候, 则认为当前版本为初始版本
}

type Version struct {
	Id      uint64 `ddb:"id"`
	Version uint64 `ddb:"version"`
	ExtInfo string `ddb:"ext_info"` //预留字段
}
