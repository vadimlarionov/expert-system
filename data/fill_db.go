package data

import (
	"encoding/csv"
	"encoding/json"
	"errors"
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

	if err = fillObjects("data/objects.json", o); err != nil {
		fmt.Printf("Can't fill objects: %s", err)
		return err
	}

	if err = fillQuestionsAndAnswers("data/questions.json", o); err != nil {
		fmt.Printf("Can't fill questions: %s", err)
		return err
	}

	if err = fillConditionalsAttribute("data/conditionals_attribute.json", o); err != nil {
		fmt.Printf("Can't fill conditionals (attribute): %s\n", err)
		return err
	}

	fmt.Printf("------------------------------------------------------\n")
	return nil
}

func fillTablesFromCSV(fileName string, fillFunc func(row []string) error) error {
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

	return fillTablesFromCSV(fileName, fillFunc)
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

	return fillTablesFromCSV(fileName, fillFunc)
}

func fillObjects(fileName string, o orm.Ormer) (err error) {
	objects, err := parseJsonObjects(fileName)
	if err != nil {
		return err
	}

	attrCount, err := o.QueryTable("attribute").Count()
	if err != nil {
		fmt.Printf("Can't count attributes from table: %s\n", err)
		return err
	}

	fmt.Printf("-----\nInsert objects and it's attribute values\n-----\n")
	for _, obj := range objects {
		fmt.Printf("Create object \"%s\"", obj.Name)
		modelObj := model.Object{Name: obj.Name}
		_, err := o.Insert(&modelObj)
		if err != nil {
			fmt.Printf("Can't insert object %s\n", obj.Name)
			return err
		}

		values := []*model.AttributeValue{}
		for attr, val := range obj.Attributes {
			attribute := model.Attribute{Text: attr}
			err = o.Read(&attribute, "Text")
			if err != nil {
				fmt.Printf("Can't read attribute %s: %s\n", attr, err)
				return err
			}

			value := &model.AttributeValue{Text: val, Attribute: &attribute}
			err = o.Read(value, "Text", "Attribute")
			if err != nil {
				fmt.Printf("Can't read attribute value \"%s\": %s\n", val, err)
				return err
			}
			values = append(values, value)
		}

		if attrCount != int64(len(values)) {
			fmt.Printf("Fill all attributes")
			return errors.New("Len of object attributes not equals number of attributes")
		}

		m2m := o.QueryM2M(&modelObj, "AttributeValues")
		for _, value := range values {
			fmt.Printf("Add attribute value \"%s\" to object \"%s\"\n", value.Text, modelObj.Name)
			_, err := m2m.Add(value)
			if err != nil {
				fmt.Printf("Can't insert attibute value \"%v\"to object \"%v\"\n", value, modelObj)
				return err
			}
		}
	}

	return nil
}

type object struct {
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
}

// TODO You can use factory pattern
func parseJsonObjects(fileName string) (objects []object, err error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Can't read file %s: %s\n", fileName, err)
		return nil, err
	}

	if err = json.Unmarshal(bytes, &objects); err != nil {
		fmt.Printf("Can't parse json: %s\n", err)
		return nil, err
	}

	return objects, nil
}

func parseJsonQuestions(fileName string) (questions []question, err error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Can't read file %s: %s\n", fileName, err)
		return nil, err
	}

	if err = json.Unmarshal(bytes, &questions); err != nil {
		fmt.Printf("Can't parse json: %s\n", err)
		return nil, err
	}

	return questions, nil
}

type answer struct {
	Answer         string `json:"answer"`
	ParameterValue string `json:"parameter_value"`
	NextQuestion   int    `json:"next_question"`
}

type question struct {
	Parameter      string   `json:"parameter"`
	Question       string   `json:"question"`
	QuestionNumber int      `json:"question_number"`
	IsSelect       bool     `json:"is_select"`
	Answers        []answer `json:"answers"`
}

func fillQuestionsAndAnswers(fileName string, o orm.Ormer) (err error) {
	questions, err := parseJsonQuestions(fileName)
	if err != nil {
		return err
	}

	for _, q := range questions {
		p := model.Parameter{Name: q.Parameter}
		if err = o.Read(&p, "Name"); err != nil {
			fmt.Printf("Can't read \"%s\" parameter: %s\n", q.Parameter, err)
			return err
		}

		qModel := model.Question{Text: q.Question, IsSelect: q.IsSelect,
			Number: q.QuestionNumber, Parameter: &p}
		_, err = o.Insert(&qModel)
		if err != nil {
			fmt.Printf("Can't insert question: %s\n", err)
			return err
		}
		fmt.Printf("Add question \"%s\" for parameter \"%s\"\n", qModel.Text, qModel.Parameter.Name)

		if err = fillAnswers(q.Answers, &qModel, o); err != nil {
			return err
		}
	}

	return nil
}

func fillAnswers(answers []answer, q *model.Question, o orm.Ormer) (err error) {
	if answers == nil || len(answers) == 0 {
		return nil
	}

	for _, answ := range answers {
		paramVal := model.ParameterValue{Parameter: q.Parameter, Value: answ.ParameterValue}
		if err = o.Read(&paramVal, "Parameter", "Value"); err != nil {
			fmt.Printf("Can't read parameter value \"%s\" for \"%s\"\n",
				answ.ParameterValue, q.Parameter)
		}

		a := model.Answer{Question: q, Text: answ.Answer,
			NextQuestionNumber: answ.NextQuestion, Value: &paramVal}
		_, err := o.Insert(&a)
		if err != nil {
			fmt.Printf("Can't insert answer: %s\n", err)
			return err
		}
		fmt.Printf("Add answer \"%s\" for question \"%s\"\n", a.Text, q.Text)
	}

	return nil
}

type conditionalJson struct {
	IsAnd  bool                             `json:"is_and"`
	Items  []conditionalItemJson            `json:"items"`
	Result []conditionalAttributeResultJson `json:"result"`
}

type conditionalItemJson struct {
	Parameter string `json:"parameter"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

type conditionalAttributeResultJson struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

func parseJsonCondAttr(fileName string) (conditionals []conditionalJson, err error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Can't read file %s: %s\n", fileName, err)
		return nil, err
	}

	if err = json.Unmarshal(bytes, &conditionals); err != nil {
		fmt.Errorf("Can't parse json: %s\n", err)
		return nil, err
	}

	return conditionals, nil
}

func fillConditionalsAttribute(fileName string, o orm.Ormer) (err error) {
	conditionalsJson, err := parseJsonCondAttr(fileName)
	if err != nil {
		return err
	}

	for _, condJson := range conditionalsJson {
		c := model.Conditional{IsAnd: condJson.IsAnd}
		_, err = o.Insert(&c)
		if err != nil {
			fmt.Printf("Can't insert condition: %s\n", err)
			return err
		}

		if err = fillCondItemsAttr(&c, condJson.Items, o); err != nil {
			return err
		}

		if err = fillCondResultAttr(&c, condJson.Result, o); err != nil {
			return err
		}
		fmt.Printf("Insert conditional %d\n", c.Id)
	}

	return nil
}

func fillCondItemsAttr(c *model.Conditional, items []conditionalItemJson, o orm.Ormer) (err error) {
	for _, itemJson := range items {
		p := model.Parameter{Name: itemJson.Parameter}
		err = o.Read(&p, "Name")
		if err != nil {
			fmt.Printf("Can't read parameter \"%s\": %s\n", itemJson.Parameter, err)
			return err
		}

		item := model.ConditionalItem{Condition: c, Parameter: &p,
			Operation: itemJson.Operation, Value: itemJson.Value}

		_, err := o.Insert(&item)
		if err != nil {
			fmt.Printf("Can't insert conditional item \"%v\": %s\n", item, err)
			return err
		}
	}
	return nil
}

func fillCondResultAttr(c *model.Conditional, result []conditionalAttributeResultJson, o orm.Ormer) (err error) {
	for _, res := range result {
		attr := model.Attribute{Text: res.Attribute}
		err = o.Read(&attr, "Text")
		if err != nil {
			fmt.Printf("Can't read attribute \"%s\": %s", res.Attribute, err)
			return err
		}

		attrVal := model.AttributeValue{Attribute: &attr, Text: res.Value}
		err = o.Read(&attrVal, "Attribute", "Text")
		if err != nil {
			fmt.Printf("Can't read attribute value \"%s\": %s\n", res.Value, err)
			return err
		}

		condAttrRes := model.ConditionalAttributeResult{Conditional: c,
			Attribute: &attr, AttributeValue: &attrVal}
		_, err = o.Insert(&condAttrRes)
		if err != nil {
			fmt.Printf("Can't insert conditional attribute result: %s\n", err)
			return err
		}
	}
	return nil
}
