package gosafe

import (
	"github.com/xinpianchang/xservice/v2/pkg/log"
	"go.uber.org/zap"
)

func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("gosafe panic", zap.Any("err", err))
			}
		}()

		f()
	}()
}
