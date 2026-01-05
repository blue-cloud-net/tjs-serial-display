package models

import "time"

// UpgradeProgress 升级进度信息
type UpgradeProgress struct {
	Current    int64         // 已发送字节数
	Total      int64         // 总字节数
	Percentage float64       // 百分比
	Speed      int64         // 传输速度 (字节/秒)
	Elapsed    time.Duration // 已用时间
	Remaining  time.Duration // 预计剩余时间
}

// UpgradeProgressCallback 升级进度回调函数类型
type UpgradeProgressCallback func(progress *UpgradeProgress)
