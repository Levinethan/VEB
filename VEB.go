package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type VEB struct {
	u , min ,max int //u 数据  min最小  max 最大
	summary *VEB
	cluster []*VEB
}

func (V VEB)Max ()int  {
	return V.max
}

func (V VEB)Min() int  {
	return V.min
}
//起到一个附加作用 log2 log2
func LowerSqrt(u int)int{
	return int(math.Pow(2.0,math.Floor(math.Log2(float64(u))/2)))
}
func HigherSqrt(u int)int{
	return int(math.Pow(2.0,math.Ceil(math.Log2(float64(u))/2)))
}

//计算x  存储的深度
func (V VEB)High(x int) int{
	return int(  math.Floor(float64(x) /float64(LowerSqrt(V.u)))   )
}

func (V VEB)Low (x int) int  {
	return x%LowerSqrt(V.u)
}

//计算并返回x和y的索引
func (V VEB)Index (x , y int) int  {
	return x *LowerSqrt(V.u) + y
}
//创建VEB树
func CreateTree(size int)*VEB {
	if size < 0 {
		return nil
	}

	x := math.Ceil(math.Log2(float64(size)))   //size = 100 2^6=64

	u := int(math.Pow(2,x))  //2的x次方

	V := new(VEB)  //新建一个节点
	V.min,V.max=-1,-1
	V.u = u

	if u == 2{
		return V  //构造完成
	}

	clustercount := HigherSqrt(u)  //计算cluster的数量
	clustersize := LowerSqrt(u)  //计算cluster的大小
	for i:=0;i<clustercount;i++{
		V.cluster= append(V.cluster,CreateTree(clustersize))
	}

	summarysize := HigherSqrt(u)
	V.summary = CreateTree(summarysize)

	return V
}

//判断节点是否存在
func (V VEB)IsMember(x int) bool  {
	if x==V.min ||x==V.max{
		return true
	}else if V.u==2{
		return false
	}else{
		return V.cluster[V.High(x)].IsMember(V.Low(x))
	}
}

//插入 删除
func (V *VEB)Insert(x int)  {
	if V.min == -1{
		V.min,V.max=x,x
	}else {
		if x < V.min{
			V.min , x = x,V.min  //exchange data
		}

		if V.u>2{
			if V.cluster[V.High(x)].Min()==-1{
				V.summary.Insert(V.High(x))
				V.cluster[V.High(x)].min,V.cluster[V.High(x)].max = V.Low(x),V.Low(x)
			}else {
				V.cluster[V.High(x)].Insert(V.Low(x))
			}
		}

		if x > V.max{
			V.max = x
		}
	}


}

func (V *VEB)Delete(x int) {
	if V.summary == nil || V.summary.Min() ==-1{
		//无非空簇
		if x == V.min && x == V.max{
			V.min ,V.max = -1 ,-1
		}else if x ==V.min{
			V.min = V.max  //重合删除
			//两个元素  x最小
		}else {
			V.max = V.min
		}


	}else {
		//存在非空簇
		if x == V.min{
			//取得最小  在cluster
			y := V.Index(V.summary.min,V.cluster[V.summary.min].min)
			V.min = y //备份
			//取得最接近的  然后赋值给v。min   然后对cluster 做改变
			V.cluster[V.High(x)].Delete(V.Low(y))
			if V.cluster[V.High(y)].min == -1{   //仅有的数据
				V.summary.Delete(V.High(y))
			}
		}else if x==V.max{
			y := V.Index(V.summary.max,V.cluster[V.summary.max].max)
			V.cluster[V.High(y)].Delete(V.Low(y))
			if V.cluster[V.High(y)].min == -1{
				V.summary.Delete(V.High(y))
			}
			if V.summary == nil || V.summary.min == -1 {
				if V.min == y{
					V.min,V.max = -1,-1
				}else {
					V.max = V.min //重合删除
				}
			}else {
				V.max= V.Index(V.summary.max,V.cluster[V.summary.max].max)
			}
		}else {
			V.cluster[V.High(x)].Delete(V.Low(x))
			if V.cluster[V.High(x)].min == -1 {
				V.summary.Delete(V.High(x))
			}
		}
	}

}
//找到x后续位置
func (V VEB)Successor(x int) int {
	if V.u==2{
		if x == 0 && V.max==1{
			return 1
		}else {
			return -1
		}
	}else if V.min != -1 && x < V.min {
		return V.min  //如果min!=-1  找不到
	}else {
		maxlow := V.cluster[V.High(x)].Max()  //最大值 用来约束范围
		if maxlow != -1 &&V.Low(x) < maxlow{  //定位cluster节点
			offset := V.cluster[V.High(x)].Successor(V.Low(x))
			return V.Index(V.High(x),offset)
		}else {  //进入子节点
			succCluster := V.summary.Successor(V.High(x))
			if succCluster==-1{
				return -1
			}else {
				offset := V.cluster[succCluster].Min()
				return V.Index(succCluster,offset)
			}
		}
	}
}

//找到x的前驱位置
func (V VEB)Predecessor(x int) int {
	if V.u==2{
		if x == 1 && V.max==0{
			return 0
		}else {
			return -1  //找不到
		}
	}else if V.max != -1 && x > V.max {
		return V.max  //如果min!=-1  找不到
	}else {
		minlow := V.cluster[V.High(x)].Min() //最大值 用来约束范围
		if minlow != -1 &&V.Low(x) > minlow{  //定位cluster节点
			offset := V.cluster[V.High(x)].Predecessor(V.Low(x))
			return V.Index(V.High(x),offset)
		}else {  //进入子节点
			preCluster := V.summary.Predecessor(V.High(x))
			if preCluster==-1{
				if V.min!=-1 && x >V.min{
					return V.min
				}else {
					return -1
				}
			}else {
				offset := V.cluster[preCluster].Max()
				return V.Index(preCluster,offset)
			}
		}
	}
}

//统计树的节点
func (V VEB)Count() int {
	if V.u == 2{
		return 1
	}
	sum := 1
	for i := 0 ; i<len(V.cluster);i ++{
		sum +=V.cluster[i].Count()   //如果是0 不需要操作  如果是1 需要相加
	}
	sum +=V.summary.Count()
	return sum
}

//递归实现   打印层级
func (V VEB)PrintFunc(level int,clusterNo int,summary bool)  {
	space := " "
	for i:=0;i<level;i++{
		space+="\t"
	}
	if level ==0{
		fmt.Printf("%vR:{u:%v,min:%v,max:%v,cluster:%v}\n",space,V.u,V.min,V.max,len(V.cluster))

	}else {
		if summary{
			fmt.Printf("%vS:{u:%v,min:%v,max:%v,cluster:%v}\n",space,V.u,V.min,V.max,len(V.cluster))

		}else {
			fmt.Printf("%vC[%v]:{u:%v,min:%v,max:%v,cluster:%v}\n",space,clusterNo,V.u,V.min,V.max,len(V.cluster))
		}
	}
	if len(V.cluster) > 0{
		V.summary.PrintFunc(level+1,0,true)
		for i:= 0 ;i<len(V.cluster);i++{
			V.cluster[i].PrintFunc(level+1,i,false)
		}
	}


}

func (V VEB)Print()  {
	V.PrintFunc(0,0,false)
}

func (V *VEB)Clear()  {
	V.min,V.max = -1 , 1
	if V.u==2{
		return
	}
	for i:=0;i<len(V.cluster);i++{
		V.cluster[i].Clear()
	}
	V.summary.Clear()

	//双层递归
}

func (V *VEB)Fill()  {
	for i:=0;i<V.u;i++{
		V.Insert(i)
	}
}

//返回元素
func (V VEB)Members() []int  {  //[]int 返回整数类型数组
	members := []int{}
	for i:=0;i<V.u;i++{
		if V.IsMember(i){
			members=append(members,i)
		}
	}
	return members
}

func arrayContains(arr[] int ,value int)bool{
	for i:=0;i<len(arr);i++{
		if arr[i]==value{
			return true
		}
	}
	return false
}
func  makerandom(max int)[]int{
	myrand:=rand.New(rand.NewSource(int64(time.Now().Nanosecond()*1)))//时间随机数
	keys:=[]int{}



	keyNo:=myrand.Intn(max)//创建随机数
	for i:=0;i<keyNo;i++{
		mykey:=myrand.Intn(max-1)//取得随机数
		if arrayContains(keys,mykey)==false{
			keys=append(keys,mykey)
		}
	}
	sort.Ints(keys)//排序

	return keys

}
func main()  {
	maxUp := 10
	for i:=1 ; i<maxUp;i++{
		u := int(math.Pow(2,float64(i)))
		V := CreateTree(u)
		keys := makerandom(u)//随机数组
		fmt.Println("keys",keys)

		for j:=0;j<len(keys);j++{
			V.Insert(keys[j])//插入数据
		}
		for j :=0 ; j < u ; j++ {
			if j > 0 && j <len(keys){

				fmt.Println(V.IsMember(keys[j]))
			}

		}
		V.Print()
		fmt.Println("----------------------------------------")
	}
}