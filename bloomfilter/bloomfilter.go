package bloomfilter

type filter struct{
    bitmap []uint8
    queryNum uint32
    setNum uint32
    collisionCount uint32
}


