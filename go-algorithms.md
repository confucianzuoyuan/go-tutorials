1. 翻书问题或者跳台阶问题。每次可以翻1页书或者2页书，一本N页的书共有多少中翻法。

```go
package main

import (
        "fmt"
)

// O(2^N)
func Fibonacci(n int) int {
        if n == 0 || n == 1 {
                return 1
        }
        if n > 1 {
                return Fibonacci(n-1) + Fibonacci(n-2)
        }
        return 0
}

// O(N)
// dynamic programming
func Fibonacci1(n int) int {
        array := make([]int, n+1)

        array[0] = 1
        array[1] = 1
        i := 2
        for {
                array[i] = array[i-1] + array[i-2]
                i++
                if i > n {
                        break
                }
        }

        return array[n]
}

func main() {
        fmt.Println(Fibonacci1(45))
}

// Fibonacci(100) = Fibonacci(99) + Fibonacci(98) = Fibonacci(98) + Fibonacci(97) + Fibonacci(98)
```
