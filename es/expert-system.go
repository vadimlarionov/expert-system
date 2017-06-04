package es

import (
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/vadimlarionov/expert-system/model"
	"io"
	"os"
	"strings"
)

func Fill() {
	o := orm.NewOrm()
	o.Using("default")

	fillAttributes(o)
}

func fillAttributes(o orm.Ormer) (err error) {
	fileName := "data/attributes.csv"
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Can't open %s: %s\n", fileName, err)
		return err
	}

	r := csv.NewReader(f)
	// Skip the header
	_, err = r.Read()
	if err != nil {
		fmt.Printf("Can't skip header: %s\n", err)
		return err
	}

	fmt.Printf("-----\nInsert attributes and it's values\n-----\n")
	for {
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Can't read file: %s\n", err)
				return err
			}
		}

		attribute := model.Attribute{Text: row[0]}
		_, err = o.Insert(&attribute)
		if err != nil {
			fmt.Printf("Can't insert attribute record: %s\n", err)
			return err
		}
		fmt.Printf("Attribute{id:%d text:\"%s\"}\n", attribute.Id, attribute.Text)

		values := strings.Split(row[1], ",")
		for _, value := range values {
			attrVal := model.AttributeValue{Attribute: &attribute, Text: value}
			o.Insert(&attrVal)
			fmt.Printf("Value{id:%d text:\"%s\"}\n", attrVal.Id, attrVal.Text)
		}
	}
	fmt.Printf("-----\n")

	return nil
}
