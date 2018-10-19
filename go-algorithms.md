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

2. 已知一颗二叉树的先序遍历序列为ABCDEFG，中序遍历为CDBAEGF，能否唯一确定一颗二叉树？如果可以，请画出这颗二叉树。

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

func buildTree(pre, in []int) *TreeNode {

	if len(in) == 0 {
		return nil
	}

	res := &TreeNode{
		Val: pre[0],
	}

	if len(in) == 1 {
		return res
	}

	idx := indexOf(res.Val, in)

	res.Left = buildTree(pre[1:idx+1], in[:idx])
	res.Right = buildTree(pre[idx+1:], in[idx+1:])

	return res
}

func indexOf(val int, nums []int) int {
	for i, v := range nums {
		if v == val {
			return i
		}
	}

	return 0
}
```

3. 二分查找及其变种

```go
package main

import (
	"fmt"
	"time"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func LinearSearch(array []int, t int) bool {
	i := 0
	for i < len(array) {
		if array[i] == t {
			return true
		}
		i++
	}
	return false
}

func BinarySearch(array []int, t int) bool {
	left := 0
	right := len(array) - 1

	for left <= right {
		mid := (left + right) / 2
		if array[mid] < t {
			left = mid + 1
		} else if array[mid] > t {
			right = mid - 1
		} else {
			return true
		}
	}

	return false
}

func main() {
	array := makeRange(0, 1000000000)
	time1 := time.Now()
	bool := LinearSearch(array, 1000000001)
	time2 := time.Now()
	fmt.Println("time2-time1: ", time2.Sub(time1))
	fmt.Println("bool: ", bool)
	time3 := time.Now()
	bool1 := BinarySearch(array, 1000000001)
	time4 := time.Now()
	fmt.Println("time4-time3: ", time4.Sub(time3))
	fmt.Println("bool1: ", bool1)
}
```

```go
package main

import "fmt"

func searchMatrix(matrix [][]int, t int) bool {
        i := 0
        j := len(matrix[0]) - 1
        for i < len(matrix) && j >= 0 {
                if matrix[i][j] == t {
                        return true
                } else if matrix[i][j] > t {
                        j -= 1
                } else {
                        i += 1
                }
        }
        return false
}

func main() {
        matrix := [][]int{
                []int{1, 3, 5, 7},
                []int{10, 11, 16, 20},
                []int{23, 30, 34, 50},
        }
        fmt.Println(searchMatrix(matrix, 16))
}
```

4. 翻转链表

```go
package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	if nil == head {
		return nil
	}

	dummy := head
	tmp := head

	for head != nil && head.Next != nil {
		dummy = head.Next
		head.Next = dummy.Next
		dummy.Next = tmp
		tmp = dummy
	}

	return dummy
}

func recursive(head *ListNode) {
	if head == nil {
		return
	}

	fmt.Println(head.Val)
	recursive(head.Next)
}

func recursiveArray(array []int) {
	if len(array) == 0 {
		return
	}

	fmt.Println(array[0])
	recursiveArray(array[1:])
}

func main() {
	head := &ListNode{1, nil}
	head.Next = &ListNode{2, nil}
	head.Next.Next = &ListNode{3, nil}
	tmp := reverseList(head)
	for tmp != nil {
		fmt.Println(tmp.Val)
		tmp = tmp.Next
	}
	recursive(tmp)
	array := []int{1, 2, 3, 4, 5}
	recursiveArray(array)
}
```
