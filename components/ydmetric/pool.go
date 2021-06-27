package ydmetric

import "github.com/panjf2000/ants/v2"

//消费协程池
func NewConsumerPool(cap int, f func(args interface{})) (*ants.PoolWithFunc, error) {
	//通用协程池
	return ants.NewPoolWithFunc(cap, f)
}
