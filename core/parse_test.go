package core

import (
	"fmt"
	"testing"
)

var file1 = "C:\\Users\\leon\\Desktop\\pr-cvbs-stor07_241218_0000.nmon"

func TestTableHeader(t *testing.T) {
	data, err := GetAllCategory(file1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)
}

func TestGetAAA(t *testing.T) {
	data, err := GetAAA(file1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)
}

func TestGetBBBP(t *testing.T) {
	data, err := GetBBBP(file1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)
}

func TestGetAllCategory(t *testing.T) {
	data, err := GetAllCategory(file1)
	if err != nil {

		t.Fatal(err)
	}

	fmt.Printf("%+v \n", data)
}
