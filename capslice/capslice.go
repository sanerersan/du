package capslice

import (
	"sync"
)

type FileCapacityInfo struct {
	Name string
	Path string
	Cap int64
}

type CapSlice struct {
	sortAsName bool
	sync.Mutex
	capInfo []*FileCapacityInfo
}

func (cs *CapSlice) Len() int {
	return len(cs.capInfo)
}

func (cs *CapSlice)Less(i, j int) bool{
	if cs.sortAsName {
		return cs.capInfo[i].Name < cs.capInfo[j].Name
	}

	return cs.capInfo[i].Cap < cs.capInfo[j].Cap
}

func (cs *CapSlice) Swap(i, j int) {
	cs.capInfo[i],cs.capInfo[j] = cs.capInfo[j],cs.capInfo[i]
}

func (cs *CapSlice)SetSortByName() {
	cs.sortAsName = true
}

func (cs *CapSlice) Append(capInfos... *FileCapacityInfo) {
	cs.Lock()
	defer cs.Unlock()
	capLen := len(cs.capInfo)
	cs.tryGrow(len(capInfos))

	copy(cs.capInfo[capLen:],capInfos)
}

func (cs *CapSlice) tryGrow(n int) {
	capLen := len(cs.capInfo)
	capCap := cap(cs.capInfo)

	if (capLen + n <= capCap) {
		cs.capInfo = cs.capInfo[:capLen + n]
		return 
	}

	newCapInfo := make([]*FileCapacityInfo,capCap * 2 + n)[:capLen + n]
	copy(newCapInfo,cs.capInfo)
	cs.capInfo = newCapInfo
	return
}

func (cs *CapSlice) GetCaps() []*FileCapacityInfo {
	cs.capInfo = cs.capInfo[:len(cs.capInfo)]
	return cs.capInfo
}

func NewCapSlice() *CapSlice {
	return &CapSlice{}
}