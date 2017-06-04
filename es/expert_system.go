package es

import (
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/vadimlarionov/expert-system/model"
	"io"
	"os"
	"strconv"
	"strings"
)

func Fill() {
	o := orm.NewOrm()
	o.Using("default")

	fillAttributes("data/attributes.csv", o)
	fillParameters("data/parameters.csv", o)
}

func fillAttributes(fileName string, o orm.Ormer) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Can't open %s: %s\n", fileName, err)
		return err
	}
	defer f.Close()

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
			if _, err = o.Insert(&attrVal); err != nil {
				fmt.Printf("Can't insert attribute value: %s", err)
				return nil
			}
			fmt.Printf("Value{id:%d text:\"%s\"}\n", attrVal.Id, attrVal.Text)
		}
	}
	fmt.Printf("-----\n")

	return nil
}

func fillParameters(fileName string, o orm.Ormer) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Can't open %s: %s\n", fileName, err)
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// Skip the header
	_, err = r.Read()
	if err != nil {
		fmt.Printf("Can't skip header: %s\n", err)
		return err
	}

	fmt.Printf("-----\nInsert parameters and it's values\n-----\n")
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

		isSelect, err := strconv.ParseBool(row[1])
		if err != nil {
			fmt.Printf("Can't convert %s to bool", row[1])
			return err
		}

		parameter := model.Parameter{Name: row[0], IsSelect: isSelect}
		_, err = o.Insert(&parameter)
		if err != nil {
			fmt.Printf("Can't insert parameter record: %s\n", err)
			return err
		}
		fmt.Printf("Parameter{id:%d is_select:%t name:\"%s\"}\n",
			parameter.Id, parameter.IsSelect, parameter.Name)

		values := strings.Split(row[2], ",")
		for _, value := range values {
			paramVal := model.ParameterValue{Parameter: &parameter, Value: value}
			if _, err = o.Insert(&paramVal); err != nil {
				fmt.Printf("Can't insert parameter value: %s", err)
				return err
			}
			fmt.Printf("Value{id:%d value:\"%s\"}\n", paramVal.Id, paramVal.Value)
		}
	}
	fmt.Printf("-----\n")

	return nil
}
