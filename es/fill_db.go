package es

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/vadimlarionov/expert-system/model"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Fill() (err error) {
	o := orm.NewOrm()
	if err = o.Using("default"); err != nil {
		fmt.Printf("Can't use defauld alias: %s", err)
		return err
	}

	if err = fillAttributes("data/attributes.csv", o); err != nil {
		return err
	}

	if err = fillParameters("data/parameters.csv", o); err != nil {
		return err
	}

	return nil
}

func fillFromCSV(fileName string, fillFunc func(row []string) error) error {
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

		if err = fillFunc(row); err != nil {
			fmt.Printf("Can't fill row: %s", err)
			return err
		}
	}

	return nil
}

func fillAttributes(fileName string, o orm.Ormer) (err error) {
	fmt.Printf("-----\nInsert attributes and it's values\n-----\n")
	fillFunc := func(row []string) error {
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
		return nil
	}

	return fillFromCSV(fileName, fillFunc)
}

func fillParameters(fileName string, o orm.Ormer) (err error) {
	fmt.Printf("-----\nInsert parameters and it's values\n-----\n")
	fillFunc := func(row []string) error {
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
		return nil
	}

	return fillFromCSV(fileName, fillFunc)
}

func fillObjects(fileName string, o orm.Ormer) (err error) {
	_, err = parseJson(fileName)
	if err != nil {
		return err
	}

	return nil
}

func parseJson(fileName string) (jsonData *map[string]interface{}, err error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Can't read file %s: %s", fileName, err)
		return nil, err
	}

	if err = json.Unmarshal(bytes, &jsonData); err != nil {
		fmt.Printf("Can't parse json: %s", err)
		return nil, err
	}
	return jsonData, nil
}
