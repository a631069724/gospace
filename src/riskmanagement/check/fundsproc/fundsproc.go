package fundsproc

import (
	"riskmanagement/conf"
	"riskmanagement/log"
	"riskmanagement/subject"
	"riskmanagement/tranpack"
	"riskmanagement/tranpack/position"

	"github.com/liangdas/mqant/utils"
)

type FundsProcer struct {
	redis *conf.Redis
	/*观察者对象 注册fundsprocer到观察者 观察者将通知事件通过update方法调用*/
	subject.Observable
}

func NewFundsProcer(ob subject.Observable, redis *conf.Redis) *FundsProcer {
	p := &FundsProcer{
		redis:      redis,
		Observable: ob,
	}
	ob.AddObserver(p)
	return p
}

func (this *FundsProcer) ClosePosition(p position.Position) {
	reqData, _ := tranpack.ClosePosition(p, "risk", "21")
	/*组包并发送，可以将redis发送的代码封装成redis客户端对象进行处理，以后有时间可以优化*/
	pool := utils.GetRedisFactory().GetPool(this.redis.Uri).Get()
	defer pool.Close()
	log.Info("Position:%s Close Data:%b", p.GetPositionId(), reqData)
	pool.Do("lpush", this.redis.WriteId, reqData)
}

//func (this *FundsProcer) DeleteFollowRelation(fwUid string) {
/*删除跟单关系，可以将数据库删除操作的代码封装成数据库客户端对象进行处理，以后有时间可以优化*/
//}

func (this *FundsProcer) Update(e interface{}) {
	p := e.(position.Position)
	this.ClosePosition(p)
	//this.DeleteFollowRelation(p.GetUserId())
}
