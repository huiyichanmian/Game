package main

import (
	"fmt"
	"os"
	"math/rand"
	"Clib"
	"time"
)

//游戏界面大小
const WIDE int = 20
const HIGH int = 20

//食物坐标
var food Food
//分数
var score int = 0
//蛇的移动偏移值
var dx int = 0
var dy int = 0

//坐标父类
type Postion struct {
	X int
	Y int
}

//蛇子类
type Snake struct {
	pos  [WIDE * HIGH]Postion //坐标位置
	size int                  //长度
	dir  byte                 //方向
}

//食物子类
type Food struct {
	Postion
}

//绘制界面
func DrawUI(p Postion, ch byte) {

	//通过C函数更改控制台光标
	Clib.GotoPostion(p.X*2+4, p.Y+2)
	fmt.Fprintf(os.Stderr, "%c", ch)

}

//蛇初始化方法
func (s *Snake) Init() {
	//初始化蛇的长度
	s.size = 2
	//初始化蛇的坐标位置
	s.pos[0].X = WIDE / 2
	s.pos[0].Y = HIGH / 2

	s.pos[1].X = WIDE/2 - 1
	s.pos[1].Y = HIGH / 2

	//初始化蛇的方向
	//U 上 L 左 R 右 D 下 P 暂停
	s.dir = 'R'

	// 输出初始画面
	fmt.Fprintln(os.Stderr,
		`
  #-----------------------------------------#
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  #-----------------------------------------#
`)

	//食物初始化
	food = Food{Postion{rand.Intn(WIDE), rand.Intn(HIGH)}}
	//判断  食物和蛇 障碍物重合（需完善）
	DrawUI(food.Postion, '#')
	//Clib.Direction()
	//调用C函数Direction()通过键盘控制方向
	//go独立程序 用来接收键盘输入的值 非阻塞式调用
	go func() {
		for {
			switch Clib.Direction() {
			//方向上  W|w|↑
			case 72, 87, 119:
				s.dir = 'U'
				//方向左
			case 65, 97, 75:
				s.dir = 'L'
				//方向右
			case 100, 68, 77:
				s.dir = 'R'
				//方向下
			case 83, 115, 80:
				s.dir = 'D'
				//暂停  空格键
			case 32:
				s.dir = 'P'
			}
		}
	}()
}

//开始游戏
func (s *Snake) PlayGame() {
	for {
		//程序更新周期  333毫秒  3可以表示等级 (在关卡中 表示蛇的速度 需要完善)
		time.Sleep(time.Second / 3)
		//time.Sleep(time.Millisecond * 300)

		//暂停
		if s.dir == 'P' {
			continue
		}

		//蛇头和墙的碰撞判断
		if s.pos[0].X < 0 || s.pos[0].X >= WIDE || s.pos[0].Y < 0 || s.pos[0].Y >= HIGH {
			//将游戏光标设置在末尾
			Clib.GotoPostion(0, 23)
			return
		}

		//蛇头和身体碰撞
		for i := 1; i < s.size; i++ {
			if s.pos[0].X == s.pos[i].X && s.pos[0].Y == s.pos[i].Y {
				//将游戏光标设置在末尾
				Clib.GotoPostion(0, 23)
				return
			}
		}

		//蛇和食物的碰撞
		if s.pos[0].X == food.X && s.pos[0].Y == food.Y {
			//蛇身体增加
			s.size++
			//分数增加
			score++
			//生成新的食物
			food = Food{Postion{rand.Intn(WIDE), rand.Intn(HIGH)}}
			DrawUI(food.Postion, '#')
		}

		//方向控制
		switch s.dir {
		case 'U':
			dx = 0
			dy = -1
		case 'L':
			dx = -1
			dy = 0
		case 'R':
			dx = 1
			dy = 0
		case 'D':
			dx = 0
			dy = 1

		}
		//记录最后一节坐标
		lp := s.pos[s.size-1] //Postion
		//顺序为每一节蛇坐标赋值
		for i := s.size - 1; i > 0; i-- {
			//用前一节坐标位当前节赋值
			s.pos[i] = s.pos[i-1]
			DrawUI(s.pos[i], '*')
		}
		//将蛇移动过的位置 设为空格
		DrawUI(lp, ' ')

		//更新蛇头坐标
		s.pos[0].X += dx
		s.pos[0].Y += dy
		DrawUI(s.pos[0], '@')
	}
}

func main() {

	//创建随机数种子
	rand.Seed(time.Now().UnixNano())
	//去掉控制台光标
	Clib.HideCursor()

	var s Snake
	//蛇的初始化
	s.Init()

	//开始游戏
	s.PlayGame()
	fmt.Println("分数：",score)

	//20秒后自动关闭界面
	time.Sleep(time.Second*20)
	//fmt.Scan(&score)

}
