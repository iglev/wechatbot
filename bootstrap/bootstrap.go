package bootstrap

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/iglev/wechatbot/handlers"
	"github.com/iglev/wechatbot/pkg/logger"
)

func Run() {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	handler, err := handlers.NewHandler()
	if err != nil {
		logger.Danger("register error: %v", err)
		return
	}
	bot.MessageHandler = handler

	// 注册登陆二维码回调
	bot.UUIDCallback = handlers.QrCodeCallBack

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// 执行热登录
	err = bot.HotLogin(reloadStorage)
	if err != nil {
		logger.Warning("hot login fail, try scanned login")
		// 热登录失败，尝试扫码登录
		bot.LoginCallBack = func(body openwechat.CheckLoginResponse) {
			errDumpTo := bot.DumpTo(reloadStorage)
			if errDumpTo != nil {
				logger.Warning("storage hot reload info fail!!!, err=%+v", errDumpTo)
			}
			logger.Info("storage hot reload info success!!!")
		}
		err = bot.Login()
		if err != nil {
			logger.Warning("login error: %v", err)
			return
		}
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	err = bot.Block()
	logger.Info("END err=%+v", err)
}
