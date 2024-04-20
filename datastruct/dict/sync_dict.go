package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func (dict *SyncDict) Len() int {
	length := 0

	dict.m.Range(func(key, value any) bool {
		length += 1
		return true
	})

	return length
}

func (dict *SyncDict) Get(key string) (val any, exists bool) {
	return dict.m.Load(key)
}

func (dict *SyncDict) Put(key string, val any) (result int) {
	_, exists := dict.m.Load(key)
	dict.m.Store(key, val)

	if exists {
		// modify not insert
		return 0
	}

	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val any) (result int) {
	_, exists := dict.m.Load(key)

	if exists {
		return 0
	}

	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExists(key string, val any) (result int) {
	_, exists := dict.m.Load(key)

	if exists {
		dict.m.Store(key, val)
		return 1
	}

	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, exists := dict.m.Load(key)

	dict.m.Delete(key)

	if exists {
		return 1
	}

	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		return consumer(key.(string), value)
	})
}

func (dict *SyncDict) Keys() []string {
	ret := make([]string, dict.Len())

	idx := 0
	dict.m.Range(func(key, value any) bool {
		ret[idx] = key.(string)
		idx += 1
		return true
	})

	return ret
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	ret := make([]string, limit)

	for i := range limit {
		// sync.Map.Range 是无序的
		dict.m.Range(func(key, value any) bool {
			ret[i] = key.(string)
			return false
		})
	}

	return ret
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	ret := make([]string, limit)

	idx := 0
	dict.m.Range(func(key, value any) bool {
		ret[idx] = key.(string)
		idx += 1
		return idx < limit
	})

	return ret
}

func (dict *SyncDict) clear() {
	*dict = *MakeSyncDict()
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}
