package main

import "fmt"

type EvictionAlgo interface {
	evict(c *Cache)
	update(c *Cache)
}

type algo struct {
}

func (l *algo) update(c *Cache) {
	for k, v := range c.record {
		c.record[k] = v + 1
	}
}

type Fifo struct {
	algo
}

func (l *Fifo) evict(c *Cache) {
	oldNum := 0
	oldKey := ""

	for k, v := range c.record {
		if oldNum == 0 || v > oldNum {
			oldNum = v
			oldKey = k
		}
	}

	delete(c.storage, oldKey)
	delete(c.record, oldKey)
	fmt.Println("Evicting by fifo strategy: ", oldKey)
}

type Cache struct {
	capacity     int
	maxCapacity  int
	storage      map[string]string
	record       map[string]int
	evictionAlgo EvictionAlgo
}

func (c *Cache) add(key, value string) {
	if c.capacity == c.maxCapacity {
		c.evictionAlgo.evict(c)
		c.capacity--
	}
	c.capacity++
	c.storage[key] = value
	c.record[key] = 1
	c.evictionAlgo.update(c)
}

func (c *Cache) delete(key string) {
	delete(c.storage, key)
}

func initCache(e EvictionAlgo) *Cache {
	storage := make(map[string]string)
	records := make(map[string]int)
	return &Cache{
		storage:      storage,
		record:       records,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

type Lifo struct {
	algo
}

func (l *Lifo) evict(c *Cache) {
	youngNum := 100000000
	youngKey := ""

	for k, v := range c.record {
		if v < youngNum {
			youngNum = v
			youngKey = k
		}
	}

	delete(c.storage, youngKey)
	delete(c.record, youngKey)
	fmt.Println("Evicting by lifo strategy: ", youngKey)
}
func main() {

	//cache := initCache(&LRU{})
	//cache.add("a", "1")
	//cache.add("b", "2")
	//cache.add("c", "3")

	println()
	println("Strategy pattern")
	//lifo := &Lifo{}
	//cache1 := initCache(lifo)

	//cache1.add("a", "1")
	//cache1.add("b", "2")
	//cache1.add("c", "3")

	fifo := &Fifo{}
	cache2 := initCache(fifo)

	cache2.add("a", "1")
	cache2.add("b", "2")
	cache2.add("c", "3")
	cache2.add("d", "4")
	cache2.add("e", "5")

}
