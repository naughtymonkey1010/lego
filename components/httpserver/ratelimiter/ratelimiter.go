package ratelimiter

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	"gitlab.yidian-inc.com/image/lego/pkg/app"
	"gitlab.yidian-inc.com/image/lego/util"
)

//限流模块

var RateLimiter = RateLimit{
	RB: make(map[string]*ratelimit.Bucket),
	RLResFunc: func(c *gin.Context) {
		res := util.Response{
			Code: 403,
			Msg:  "request trigger rate limiter",
		}
		requestId := c.Request.Header.Get(app.App.GetRequestId())
		if len(requestId) > 1 {
			res.RequestId = requestId
		}
		c.AbortWithStatusJSON(403, res)
	},
}

type RateLimit struct {
	//存储url对应的bucket
	RB map[string]*ratelimit.Bucket
	//定义返回的 url
	RLResFunc func(c *gin.Context)
	mutex     sync.Mutex
}

//添加对应的路径和限速速率
func (r *RateLimit) AddRateLimit(path string, capacity int64) {
	defer r.mutex.Unlock()
	r.mutex.Lock()
	r.RB[path] = ratelimit.NewBucketWithQuantum(1*time.Second, capacity, capacity)
}

//自定义限速返回
func (r *RateLimit) AddCustomResponseFunc(f func(*gin.Context)) {
	defer r.mutex.Unlock()
	r.mutex.Lock()
	r.RLResFunc = f
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		//获取限流
		bucket, ok := RateLimiter.RB[path]
		if ok {
			take := bucket.TakeAvailable(1)
			//如果未获取到take中断停止
			if take == 0 {
				RateLimiter.RLResFunc(c)
				return
			}
		}
		c.Next()
	}
}
