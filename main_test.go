package main

import (
	"fmt"
	"ledis/htable"
	"ledis/models"
	"ledis/pocket"
	"testing"
)

func testHashTable(t *testing.T) {
	ht := htable.NewHashTable()
	err := ht.Set("name", "gaojian")
	if err != nil {
		t.Error(err)
	}
	err = ht.Set("age", 100)
	if err != nil {
		t.Error(err)
	}
	err = ht.Set("ledis", "haha")
	if err != nil {
		t.Error(err)
	}
	v, err := ht.Get("name")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("v = ", v)
	v, err = ht.Get("age")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("v = ", v)
	v, err = ht.Get("ledis")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("v = ", v)
	fmt.Println(ht)
	//m, err := ht.ToMap()
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(m)
	ht2 := htable.NewHashTable()
	err = ht2.Set("user", ht)
	if err != nil {
		t.Error(err)
	}
	v, err = ht2.Get("user")
	if err != nil {
		t.Error(err)
	}
	//m, err = ht2.ToMap()
	//if err != nil {
	//	t.Error(err)
	//}
	fmt.Println(ht2)
}

func testModel(t *testing.T) {
	dbName := "ledis"
	tbName := "user"
	model := models.NewModel(dbName, tbName, nil)
	err := pocket.AddModel(model)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Begin to add data")
	data := make(map[string]interface{})
	data["name1"] = "gaojian"
	data["age"] = 100
	err = pocket.AddData(model, data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Finish to add data")
	record, err := pocket.GetData(model, "age", 100)
	if err != nil {
		t.Error(err)
	}
	fmt.Println((*htable.HashTable)(record))
}

func TestAll(t *testing.T) {
	t.Run("testModel", testModel)
}
