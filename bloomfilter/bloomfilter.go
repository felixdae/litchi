package bloomfilter

import "math"
import "sync"
import "errors"

type Filter struct{
    bitmap []uint8
    hashNum uint32
    numOfBit uint32
    mu *sync.RWMutex

    checkCount uint32
    setCount uint32
    collisionCount uint32
}

const MAXBUKETSIZE = 500 * 1024 * 1024 * 8

func feedHash1(in []uint8, bound uint32) (out uint32) {
    return 1 % bound;
}

func feedHash2(in []uint8, bound uint32) (out uint32) {
    return 2 % bound;
}

func RequiredBitNum(numOfSet uint32, collisionProb float64) uint64 {
    m := numOfSet
    p := collisionProb
    n := (uint32)(m * math.Log(1/p) / (math.Log(2) * math.Log(2)))
    return n
}

func RequiredHashNum(numOfSet uint32, collisionProb float64) uint32 {
    m := numOfSet
    n := RequiredBitNum(numOfSet, collisionProb);
    k := (uint32)(math.Log(2) * n / m)
    return k
}

func NewFilter(numOfSet uint32, collisionProb float64) (flt *Filter, err error) {
    n := RequiredBitNum(numOfSet, collisionProb)
    if n > MAXBUKETSIZE {
        return nil, errors.New("too buge a instance")
    }
    k := RequiredHashNum(numOfSet, collisionProb)
    bitmap := make([]uint8, n / 8)
    return &Filter{bitmap, k, n, 0, 0, 0, &sync.RWMutex{}}, nil
}

func (flt *Filter) setAtIndex(idx uint32) bool {
    mask := (1 << (idx % 8))
    flt.bitmap[idx / 8] |= mask
}

func (flt *Filter) checkAtIndex(idx uint32) bool {
    mask := (1 << (idx % 8))
    targetByte := flt.bitmap[idx / 8]
    if (targetByte & mask) == mask {
        return true
    }
    return false
}

func (flt *Filter) check(in []uint8) bool {
    flt.checkCount++;
    for i := 0; i < flt.hashNum; i++ {
        hashIndex := (feedHash1(in, flt.numOfBit) + feedHash2(in, flt.numOfBit) * i) % flt.numOfBit
        if !flt.checkAtIndex(hashIndex) {
            return false
        }
    }
    return true
}

func (flt *Filter) set(in []uint8) {
    //flt.mu.Lock()
    //defer flt.mu.Unlock()

    flt.checkCount++;
    for i := 0; i < flt.hashNum; i++ {
        hashIndex := (feedHash1(in, flt.numOfBit) + feedHash2(in, flt.numOfBit) * i) % flt.numOfBit
        if !flt.checkAtIndex(hashIndex) {
            flt.setAtIndex(hashIndex)
        }
    }
    return
}

func (flt *Filter) CAS(in []uint8) bool {
    if flt.Check(in) {
        return 1
    }

    flt.mu.Lock()
    defer flt.mu.Unlock()
    if flt.check(in) {
        return 1
    }
    flt.set(in)
    return 0
}

func (flt *Filter) Check(in []uint8) bool {
    flt.mu.RLock()
    defer flt.mu.RUnlock()
    return flt.check(in)
}
