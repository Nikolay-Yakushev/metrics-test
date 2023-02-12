package closer

import (
	"context"
	"sort"
	"sync"

	ports "github.com/Nikolay-Yakushev/lychee/ports/driver"
	"go.uber.org/zap"
)

// Sigletone implemetation. Used sync.Once instead of `mutex`
// Seems to be more clear and sophisticated way to implemet this pattern
// Compred to
//https://refactoring.guru/design-patterns/singleton/go/example

var singleinstance *Closer
var once sync.Once 

type Closer struct{
	l          *zap.Logger
	mu         sync.Mutex
	closeables []ports.Closeable
}

func(c *Closer) sortCloseables(){
	sort.SliceStable(c.closeables, func(idx1, idx2 int) bool {
		return c.closeables[idx1].GetCloseOrder() < c.closeables[idx2].GetCloseOrder()
	})
}

func New(logger *zap.Logger) *Closer {

	var cl []ports.Closeable
	var closer *Closer

	 once.Do(func(){
		closer = &Closer{
			l: logger,
			closeables: cl,
		}
		singleinstance = closer
	})
	
	return singleinstance
}

func(c * Closer) Add(cl ports.Closeable){
	c.mu.Lock()
	defer c.mu.Unlock()

	contains := func (closeables []ports.Closeable, cl ports.Closeable) bool {

		for ix := range closeables{
			if closeables[ix] == cl{
				return true
			}
		}
		return false
	}
	if !contains(c.closeables, cl){
		c.closeables = append(c.closeables, cl)
	}
	
}


func(c *Closer) Close(ctx context.Context){
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sortCloseables()
	for ix := range c.closeables{
		errorCh := make(chan error, 1)
		closeable := c.closeables[ix]

		go func(entity ports.Closeable) {
			c.l.Sugar().Debugf("%s is being stopped\n", entity.GetDescription())
			err := entity.Stop(ctx)
			errorCh <- err

		}(closeable)

		select {
			case <-ctx.Done():
				return 

			case err := <-errorCh:
				descr := closeable.GetDescription()

				if err != nil{
					c.l.Sugar().Errorf("Failed to stop `%s`\n", descr)
					continue
				}
				c.l.Sugar().Infof("`%s` has been stopped successfully\n", descr)
		}
	}
	c.l.Sugar().Infof("Successefully stopped all component\n")
}


