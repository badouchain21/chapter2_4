package main

import (
	"fmt"
	"math/rand"
	"sort"
)

//3 构造一个map以学号为 key 存储学生信息(姓名、性别、学号、年龄、成绩)，产生10 个学生信息并存入 map，然后将所有学生信息取出，按成绩排序(由高到低)存入一个切片中，然后按顺序输出学生信息。
type Student struct {
	name string
	sex string
	number string
	age string
	score int
}

type Students [] Student

func (s Students) Len() int { return len(s) }
func (s Students) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Students) Less(i, j int) bool { return s[i].score < s[j].score }

func main() {
	SortScore()
}

/**
学生成绩排序
*/
func SortScore()  {

	var sexs = [2]string{"男","女"}
	var names = [10]string{"张三","李四","王","赵","牛","何","刘","孟","齐","徐"}
	var students= make(map[int]Student)
	//产生10个学生信息
	for r := 1; r <= 10; r++ {
		s1 := Student {names[rand.Intn(10)],sexs[rand.Intn(2)],fmt.Sprint(rand.Intn(100)),fmt.Sprint(rand.Intn(100)),rand.Intn(100)}
		students [ r ] = s1
	}
	Studentss := make([]Student, 10,10)
	index :=0;
	for _,v:= range students{
		Studentss[index] = v
		index++
	}
	sort.Sort(sort.Reverse(Students(Studentss))) //按照 score 的由高到底排序
	fmt.Println(Studentss)

}
