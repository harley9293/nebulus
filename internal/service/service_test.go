package service

import (
	"sort"
	"testing"

	"github.com/harley9293/nebulus/pkg/def"
)

type Service struct {
	def.DefaultHandler
}

func (m *Service) In0Out0() {
}

func (m *Service) In1Out0(a string) {
}

func (m *Service) In0Out1() string {
	return "nice"
}

func (m *Service) In1Out1(a string) string {
	return a + " nice"
}

func (m *Service) In2Out2(a, b int) (int, int) {
	return a / b, a % b
}

func (m *Service) SliceParam(a []int) []int {
	sort.Ints(a)
	return a
}

func (m *Service) MapParam(a map[int]string) map[int]string {
	a[10000] = "gid"
	return a
}

type Person struct {
	Name string
	Age  int
}

func (m *Service) StructParam() Person {
	return Person{"harry", 18}
}

func (m *Service) PointerParam() *Person {
	return &Person{"harry", 18}
}

func init() {
	Register("Test", &Service{})
}

func TestService_SendCall(t *testing.T) {
	Send("Test.In0Out0")
	Call("Test.In0Out0")

	Send("Test.In1Out0", "hello world")
	Call("Test.In1Out0", "hello world")

	Send("Test.In0Out1")
	rsp := ""
	Call("Test.In0Out1", &rsp)

	if rsp != "nice" {
		t.Fatal(rsp)
	}

	Send("Test.In1Out1", "nice")
	rsp = ""
	Call("Test.In1Out1", "nice", &rsp)

	if rsp != "nice nice" {
		t.Fatal(rsp)
	}

	Send("Test.In2Out2", 7, 3)

	var a, b int
	Call("Test.In2Out2", 7, 3, &a, &b)

	if a != 2 || b != 1 {
		t.Fail()
	}
}

func TestService_ParamType(t *testing.T) {
	// slice
	var rSlice []int
	Call("Test.SliceParam", []int{3, 9, 4, 7, 5}, &rSlice)
	if len(rSlice) != 5 || rSlice[0] != 3 || rSlice[4] != 9 {
		t.Fatalf("%v", rSlice)
	}

	// map
	var rMap map[int]string
	Call("Test.MapParam", map[int]string{1: "hello"}, &rMap)
	if len(rMap) != 2 || rMap[1] != "hello" || rMap[10000] != "gid" {
		t.Fatalf("%v", rMap)
	}

	// struct
	var rStruct Person
	Call("Test.StructParam", &rStruct)
	if rStruct.Name != "harry" || rStruct.Age != 18 {
		t.Fatalf("%v", rStruct)
	}

	// pointer
	var rPointer *Person
	Call("Test.PointerParam", &rPointer)
	if rPointer.Name != "harry" || rPointer.Age != 18 {
		t.Fatalf("%v", rPointer)
	}
}

func TestService_ParamError(t *testing.T) {
	err := Call("Test.In2Out2", 1, 2, 3, 4, 5)
	if err == nil {
		t.Fail()
	}

	var a, b int
	err = Call("Test.In2Out2", "test", "test", &a, &b)
	if err == nil {
		t.Fail()
	}

	err = Call("Test.In2Out2", 7, 5, a, b)
	if err == nil {
		t.Fail()
	}

	var c, d string
	err = Call("Test.In2Out2", 7, 5, &c, &d)
	if err == nil {
		t.Fail()
	}
}
