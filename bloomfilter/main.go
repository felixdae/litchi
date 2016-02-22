package main

// import "bloomfilter"
import "fmt"

func main(){
    hh := []uint8{1,2,3,4,5,6,7,8,9,1,2,3,4,10,6,6,4}
    fmt.Println(FeedHash1(hh, 10000000));
    fmt.Println(FeedHash2(hh, 10000000));
}
