package core

import (
	"testing"
)

func TestExcel(t *testing.T) {
	categories, err := GetAllCategory(file1)
	if err != nil {
		t.Fatal(err)
	}

	data, err := GetCategoryAllData(file1)

	if err != nil {
		t.Fatal(err)
	}

	err = AddToExcel("nmon.xlsx", categories, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExcel2(t *testing.T) {
	BBBPData, err := GetBBBP(file1)
	if err != nil {
		t.Fatal(err)
	}
	err = AddBBBPToExcel("kk.xlsx", BBBPData)
	if err != nil {
		t.Fatal(err)
	}
}
