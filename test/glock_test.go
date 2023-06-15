package test

import (
	"fmt"
	"github.com/paceew/go-redisson/glock"
	"github.com/paceew/go-redisson/pkg/log"
	"github.com/paceew/go-redisson/redis"
	"sync"
	"testing"
	"time"
)

var gl glock.Glock

func TestGlock(t *testing.T) {
	redisentry := redis.NewRedisEntry("127.0.0.1:6379", "")
	log.SetDefaultLogConfig("./", "glock.log", "", 300, 1000, log.TRACE, false)
	logger := log.NewFieLogger("", log.Fields{}, log.TRACE)
	gl = glock.NewGlock(redisentry, glock.WithLogger(logger), glock.WithPrefix("testll"), glock.WithWatchDogTimeout(10*time.Second))
	//testGMutexWithTTl(t)
	//testGMutexWithoutTTL(t)
	//testGetReLockWithTTL(t)
	//testGReLockWithoutTTL(t)
	//testGRWMutexWithoutTTL(t)
	testGRWMutexSinger(t)
}

func testGMutexWithTTl(t *testing.T) {
	mu := gl.GetMutex("testMu")
	fmt.Println("try lock ...")
	err := mu.TryLock(30*time.Millisecond, 500*time.Second)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	fmt.Println("lock ok ")
	time.Sleep(20 * time.Second)
	fmt.Println("unlock ...")
	err = mu.UnLock()
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	fmt.Println("unlock ok ")
}

func testGMutexWithoutTTL(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		mu := gl.GetMutex("testMu2")
		err := mu.TryLock(30 * time.Millisecond)
		if err != nil {
			fmt.Printf("lock1 err:%s", err.Error())
		} else {
			fmt.Println("lock1 ok ")
		}
		time.Sleep(13 * time.Second)
		err = mu.UnLock()
		if err != nil {
			fmt.Printf("unlock1 err:%s", err.Error())
		} else {
			fmt.Println("unlock1 ok ")
		}
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	wg.Add(3)
	go func() {
		mu := gl.GetMutex("testMu2")
		err := mu.TryLock(15 * time.Second)
		if err != nil {
			fmt.Printf("lock2 err:%s", err.Error())
		} else {
			fmt.Println("lock2 ok ")
		}
		time.Sleep(5 * time.Second)
		err = mu.UnLock()
		if err != nil {
			fmt.Printf("unlock2 err:%s", err.Error())
		} else {
			fmt.Println("unlock2 ok ")
		}
		wg.Done()
	}()
	go func() {
		mu := gl.GetMutex("testMu2")
		err := mu.TryLock(20 * time.Second)
		if err != nil {
			fmt.Printf("lock3 err:%s", err.Error())
		} else {
			fmt.Println("lock3 ok ")
		}
		time.Sleep(5 * time.Second)
		err = mu.UnLock()
		if err != nil {
			fmt.Printf("unlock3 err:%s", err.Error())
		} else {
			fmt.Println("unlock3 ok ")
		}
		wg.Done()
	}()
	go func() {
		mu := gl.GetMutex("testMu2")
		err := mu.TryLock(7 * time.Second)
		if err != nil {
			fmt.Printf("lock4 err:%s", err.Error())
		} else {
			fmt.Println("lock4 ok ")
		}
		time.Sleep(5 * time.Second)
		err = mu.UnLock()
		if err != nil {
			fmt.Printf("unlock4 err:%s", err.Error())
		} else {
			fmt.Println("unlock4 ok ")
		}
		wg.Done()
	}()

	wg.Wait()
}

func testGetReLockWithTTL(t *testing.T) {
	rl := gl.GetReLock("testRe")
	fmt.Println("try lock ...")
	err := rl.TryLock(30*time.Millisecond, 500*time.Second)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	fmt.Println("lock ok ")
	time.Sleep(20 * time.Second)
	fmt.Println("unlock ...")
	err = rl.UnLock()
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	fmt.Println("unlock ok ")
}

func testGReLockWithoutTTL(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rl := gl.GetReLock("testRe2")
		err := rl.TryLock(30 * time.Millisecond)
		if err != nil {
			fmt.Printf("lock1 err:%s", err.Error())
		} else {
			fmt.Println("lock1 ok ")
		}
		time.Sleep(13 * time.Second)
		err = rl.UnLock()
		if err != nil {
			fmt.Printf("unlock1 err:%s", err.Error())
		} else {
			fmt.Println("unlock1 ok ")
		}
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	wg.Add(4)
	go func() {
		rl := gl.GetReLock("testRe2")
		err := rl.TryLock(15 * time.Second)
		if err != nil {
			fmt.Printf("lock2 err:%s", err.Error())
		} else {
			fmt.Println("lock2 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rl.UnLock()
		if err != nil {
			fmt.Printf("unlock2 err:%s", err.Error())
		} else {
			fmt.Println("unlock2 ok ")
		}
		wg.Done()
	}()
	go func() {
		rl := gl.GetReLock("testRe2")
		err := rl.TryLock(20 * time.Second)
		if err != nil {
			fmt.Printf("lock3 err:%s", err.Error())
		} else {
			fmt.Println("lock3 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rl.UnLock()
		if err != nil {
			fmt.Printf("unlock3 err:%s", err.Error())
		} else {
			fmt.Println("unlock3 ok ")
		}
		wg.Done()
	}()
	go func() {
		rl := gl.GetReLock("testRe2", "127.0.0.1:6666")
		err := rl.TryLock(7 * time.Second)
		if err != nil {
			fmt.Printf("lock4 err:%s", err.Error())
		} else {
			fmt.Println("lock4 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rl.UnLock()
		if err != nil {
			fmt.Printf("unlock4 err:%s", err.Error())
		} else {
			fmt.Println("unlock4 ok ")
		}
		wg.Done()
	}()
	go func() {
		rl := gl.GetReLock("testRe2", "127.0.0.1:6666")
		err := rl.TryLock(15 * time.Second)
		if err != nil {
			fmt.Printf("lock5 err:%s", err.Error())
		} else {
			fmt.Println("lock5 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rl.UnLock()
		if err != nil {
			fmt.Printf("unlock5 err:%s", err.Error())
		} else {
			fmt.Println("unlock5 ok ")
		}
		wg.Done()
	}()

	wg.Wait()
}

func testGRWMutexWithoutTTL(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rwm := gl.GetRWMutex("testRWMu")
		err := rwm.TryRLock(30 * time.Millisecond)
		if err != nil {
			fmt.Printf("lock1 err:%s", err.Error())
		} else {
			fmt.Println("lock1 ok ")
		}
		time.Sleep(13 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock1 err:%s", err.Error())
		} else {
			fmt.Println("unlock1 ok ")
		}
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	wg.Add(4)
	go func() {
		rwm := gl.GetRWMutex("testRWMu")
		err := rwm.TryRLock(15 * time.Second)
		if err != nil {
			fmt.Printf("lock2 err:%s", err.Error())
		} else {
			fmt.Println("lock2 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock2 err:%s", err.Error())
		} else {
			fmt.Println("unlock2 ok ")
		}
		wg.Done()
	}()
	go func() {
		rwm := gl.GetRWMutex("testRWMu")
		err := rwm.TryLock(20 * time.Second)
		if err != nil {
			fmt.Printf("lock3 err:%s", err.Error())
		} else {
			fmt.Println("lock3 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock3 err:%s", err.Error())
		} else {
			fmt.Println("unlock3 ok ")
		}
		wg.Done()
	}()
	go func() {
		rwm := gl.GetRWMutex("testRWMu")
		err := rwm.TryRLock(17 * time.Second)
		if err != nil {
			fmt.Printf("lock4 err:%s", err.Error())
		} else {
			fmt.Println("lock4 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock4 err:%s", err.Error())
		} else {
			fmt.Println("unlock4 ok ")
		}
		wg.Done()
	}()
	go func() {
		rwm := gl.GetRWMutex("testRWMu")
		err := rwm.TryLock(15 * time.Second)
		if err != nil {
			fmt.Printf("lock5 err:%s", err.Error())
		} else {
			fmt.Println("lock5 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock5 err:%s", err.Error())
		} else {
			fmt.Println("unlock5 ok ")
		}
		wg.Done()
	}()

	wg.Wait()
}

func testGRWMutexSinger(t *testing.T) {
	rwm := gl.GetRWMutex("testRWMu")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := rwm.TryRLock(30 * time.Millisecond)
		if err != nil {
			fmt.Printf("lock1 err:%s", err.Error())
		} else {
			fmt.Println("lock1 ok ")
		}
		time.Sleep(13 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock1 err:%s", err.Error())
		} else {
			fmt.Println("unlock1 ok ")
		}
		wg.Done()
	}()

	time.Sleep(1 * time.Second)
	wg.Add(4)
	go func() {
		err := rwm.TryRLock(15 * time.Second)
		if err != nil {
			fmt.Printf("lock2 err:%s", err.Error())
		} else {
			fmt.Println("lock2 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock2 err:%s", err.Error())
		} else {
			fmt.Println("unlock2 ok ")
		}
		wg.Done()
	}()
	go func() {
		err := rwm.TryLock(20 * time.Second)
		if err != nil {
			fmt.Printf("lock3 err:%s", err.Error())
		} else {
			fmt.Println("lock3 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock3 err:%s", err.Error())
		} else {
			fmt.Println("unlock3 ok ")
		}
		wg.Done()
	}()
	go func() {
		err := rwm.TryRLock(17 * time.Second)
		if err != nil {
			fmt.Printf("lock4 err:%s", err.Error())
		} else {
			fmt.Println("lock4 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock4 err:%s", err.Error())
		} else {
			fmt.Println("unlock4 ok ")
		}
		wg.Done()
	}()
	go func() {
		err := rwm.TryLock(15 * time.Second)
		if err != nil {
			fmt.Printf("lock5 err:%s", err.Error())
		} else {
			fmt.Println("lock5 ok ")
		}
		time.Sleep(5 * time.Second)
		err = rwm.UnLock()
		if err != nil {
			fmt.Printf("unlock5 err:%s", err.Error())
		} else {
			fmt.Println("unlock5 ok ")
		}
		wg.Done()
	}()

	wg.Wait()
}
