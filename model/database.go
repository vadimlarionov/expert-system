package model

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Init(username, password, dbName string, createTables bool) (err error) {
	err = orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		fmt.Printf("Can't register driver: %s\n", err)
		return err
	}

	dataSource := fmt.Sprintf("%s:%s@/%s?charset=utf8", username, password, dbName)
	err = orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		fmt.Printf("Can't register database: %s\n", err)
		return err
	}

	registerModels()
	if createTables {
		err = orm.RunSyncdb("default", true, false)
		if err != nil {
			fmt.Printf("Can't sync database: %s\n", err)
			return err
		}
	}
	return nil
}

func registerModels() {
	orm.RegisterModel(
		new(Attribute),
		new(AttributeValue),
		new(Object),
		new(Parameter),
		new(ParameterValue),
		new(Question),
		new(Answer),
		new(Conditional),
		new(ConditionalItem),
		new(ConditionalAttributeResult),
		new(ConditionalParameterResult),

		new(Quest),
		new(QuestQuestions),
		new(QuestAttribute),
		new(QuestParameter),
	)
}
