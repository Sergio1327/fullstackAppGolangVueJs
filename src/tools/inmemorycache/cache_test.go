package inmemorycache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	ID int
}

func TestCache(t *testing.T) {
	r := require.New(t)
	expiration := time.Second * 10

	cacheStorate := NewInmemoryCacheRepository[*TestStruct, int](expiration)

	for i := 0; i < 1000000; i++ {
		ts := &TestStruct{i}
		cacheStorate.Add(ts.ID, ts)
	}

	for i := 0; i < 1000000; i++ {
		ts, exists := cacheStorate.Get(i)
		r.True(exists)
		r.NotNil(ts)
		r.Equal(i, ts.ID)
	}
}

func TestClearCache(t *testing.T) {
	r := require.New(t)
	expiration := time.Second * 1

	cacheStorate := NewInmemoryCacheRepository[*TestStruct, int](expiration)

	for i := 0; i < 1000000; i++ {
		ts := &TestStruct{i}
		cacheStorate.Add(ts.ID, ts)
	}

	time.Sleep(time.Second * 1)

	t.Run("очистка", func(t *testing.T) {
		cacheStorate.ClearExpired()
		r.Empty(cacheStorate.cache)
		r.Empty(cacheStorate.lifeTime)
	})

}

func TestCacheExpiration(t *testing.T) {
	r := require.New(t)
	expiration := time.Second * 3

	cacheStorate := NewInmemoryCacheRepository[*TestStruct, int](expiration)

	for i := 0; i < 1000000; i++ {
		ts := &TestStruct{i}
		cacheStorate.Add(ts.ID, ts)
	}

	for i := 0; i < 1000000; i++ {
		ts, exists := cacheStorate.Get(i)
		r.True(exists)
		r.NotNil(ts)
		r.Equal(i, ts.ID)
	}

	time.Sleep(time.Second * 3)

	for i := 0; i < 1000000; i++ {
		ts, exists := cacheStorate.Get(i)
		r.False(exists)
		r.Nil(ts)
	}
}

func BenchmarkCache(t *testing.B) {
	r := require.New(t)
	expiration := time.Second * 1000

	cacheStorate := NewInmemoryCacheRepository[*TestStruct, int](expiration)

	t.ResetTimer()

	t.Run("создание", func(b *testing.B) {
		for i := 0; i < t.N; i++ {
			ts := &TestStruct{i}
			cacheStorate.Add(ts.ID, ts)
		}
	})

	t.Run("считывание", func(b *testing.B) {
		for i := 0; i < t.N; i++ {
			ts, exists := cacheStorate.Get(i)
			r.True(exists)
			r.NotNil(ts)
			r.Equal(i, ts.ID)
		}
	})

	t.Run("считывание параллельное", func(b *testing.B) {
		wg := &sync.WaitGroup{}
		parallel := 100
		wg.Add(parallel)

		for i := 0; i < parallel; i++ {
			go func() {
				for i := 0; i < t.N; i++ {
					ts, exists := cacheStorate.Get(i)
					r.True(exists)
					r.NotNil(ts)
					r.Equal(i, ts.ID)
				}
				cacheStorate.ClearExpired()

				wg.Done()
			}()
		}

		wg.Wait()

	})

}
