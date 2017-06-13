package es

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/vadimlarionov/expert-system/model"
)

type expertSystemType struct {
	objects      []*model.Object
	conditionals []*model.Conditional
}

func initExpertSystem(o orm.Ormer) (expertSystem *expertSystemType, err error) {
	expertSystem = &expertSystemType{}

	objects, err := loadObjects(o)
	if err != nil {
		fmt.Printf("Can't load objects: %s\n", err)
		return nil, err
	}
	expertSystem.objects = objects

	conditionals, err := loadConditions(o)
	if err != nil {
		fmt.Printf("Can't load conditionals: %s\n", err)
		return nil, err
	}
	expertSystem.conditionals = conditionals

	return expertSystem, nil
}

func loadObjects(o orm.Ormer) (objects []*model.Object, err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		fmt.Printf("Can't create query builder: %s\n", err)
		return nil, err
	}

	query := qb.Select("*").From("object").String()
	_, err = o.Raw(query).QueryRows(&objects)
	if err != nil {
		fmt.Printf("Can't execute query: %s\n", err)
		return nil, err
	}

	for _, obj := range objects {
		_, err := o.LoadRelated(obj, "AttributeValues")
		if err != nil {
			fmt.Printf("Can't load related attribute values: %s\n", err)
			return nil, err
		}

		for _, attrVal := range obj.AttributeValues {
			_, err := o.LoadRelated(attrVal, "Attribute")
			if err != nil {
				fmt.Printf("Can't load related attribute: %s\n", err)
				return nil, err
			}
		}
	}

	return objects, nil
}

func loadConditions(o orm.Ormer) (conditionals []*model.Conditional, err error) {
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		fmt.Printf("Can't create query builder: %s\n", err)
		return nil, err
	}

	query := qb.Select("*").From("conditional").String()
	_, err = o.Raw(query).QueryRows(&conditionals)
	if err != nil {
		fmt.Printf("Can't execute query: %s\n", err)
		return nil, err
	}

	for _, conditional := range conditionals {
		// Load items
		_, err = o.LoadRelated(conditional, "Items")
		if err != nil {
			fmt.Printf("Can't load related items: %s\n", err)
			return nil, err
		}

		for _, item := range conditional.Items {
			_, err := o.LoadRelated(item, "Parameter")
			if err != nil {
				fmt.Printf("Can't load related parameters: %s\n", err)
				return nil, err
			}
		}

		// Load attribute results
		_, err = o.LoadRelated(conditional, "AttributeResults")
		if err != nil {
			fmt.Printf("Can't load related attribute results: %s\n", err)
			return nil, err
		}

		for _, attrResult := range conditional.AttributeResults {
			_, err := o.LoadRelated(attrResult, "Attribute")
			if err != nil {
				fmt.Printf("Can't load related attributes: %s\n", err)
				return nil, err
			}

			_, err = o.LoadRelated(attrResult, "AttributeValue")
			if err != nil {
				fmt.Printf("Can't load related attribute value: %s\n", err)
				return nil, err
			}
		}

		// Load parameter results
		_, err = o.LoadRelated(conditional, "ParameterResults")
		if err != nil {
			fmt.Printf("Can't load related parameter results: %s\n", err)
			return nil, err
		}

		for _, paramResult := range conditional.ParameterResults {
			_, err := o.LoadRelated(paramResult, "Parameter")
			if err != nil {
				fmt.Printf("Can't load related parameters: %s\n", err)
				return nil, err
			}
		}
	}

	return conditionals, nil
}
