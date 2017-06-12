package es

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/vadimlarionov/expert-system/model"
	"github.com/vadimlarionov/expert-system/utils"
)

func StartQuest() {
	o := orm.NewOrm()

	quest := model.Quest{Username: "Vadim"}
	_, err := o.Insert(&quest)
	if err != nil {
		fmt.Printf("Can't start quest: %s\n", err)
		return
	}

	//attributesMap := make(map[string]string)
	parametersMap := make(map[int]string)

	questionNum := 1
	for questionNum != -1 {
		q := nextQuestion(o, questionNum)
		if q == nil {
			break
		}

		if _, err = o.LoadRelated(q, "Parameter"); err != nil {
			fmt.Printf("Can't load related parameter for question \"%s\": %s",
				q.Text, err)
			return
		}

		fmt.Printf("%s\n", q.Text)
		if q.IsSelect {
			_, err := o.LoadRelated(q, "Answers")
			if err != nil {
				fmt.Printf("Can't load related answers for question \"%s\": %s", q.Text, err)
				return
			}
			for i, answ := range q.Answers {
				fmt.Printf("%d: %s\n", i+1, answ.Text)
			}
			fmt.Printf("Ваш ответ: ")
			var userAnswer int
			fmt.Scanf("%d", &userAnswer)

			answer := q.Answers[userAnswer-1]
			if _, err = o.LoadRelated(answer, "Value"); err != nil {
				fmt.Printf("Can't load related parameter value for answer \"%s\": %s\n",
					answer.Text, err)
				return
			}

			err = writeQuestParameter(&quest, q.Parameter, answer.Value.Value, o)
			if err != nil {
				return
			}

			questionNum = answer.NextQuestionNumber
		} else {
			fmt.Printf("Ваш ответ: ")
			var userAnswer string
			fmt.Scanf("%s", &userAnswer)
			err = writeQuestParameter(&quest, q.Parameter, userAnswer, o)

			questionNum++
		}
	}

	fmt.Printf("%v", parametersMap)

}

func writeQuestParameter(quest *model.Quest, parameter *model.Parameter, value string, o orm.Ormer) (err error) {
	questParam := model.QuestParameter{Quest: quest, Parameter: parameter, Value: value}
	_, err = o.Insert(&questParam)
	if err != nil {
		fmt.Printf("Can't insert quest parameter \"%+v\": %s", questParam, err)
		return err
	}
	return nil
}

func nextQuestion(o orm.Ormer, expectedQuestion int) *model.Question {
	if expectedQuestion < 0 {
		return nil
	}

	q := model.Question{Number: expectedQuestion}
	err := o.Read(&q, "Number")
	if err != nil {
		fmt.Printf("Can't find question with number %d: %s\n", expectedQuestion, err)
		return nil
	}
	return &q
}

func checkConditionals(quest *model.Quest, parametersMap map[uint]string, o orm.Ormer) {
	conditionals, err := loadConditions(o)
	if err != nil {
		return
	}

	for _, conditional := range conditionals {
		result := false
		for _, item := range conditional.Items {
			if val, ok := parametersMap[item.Parameter.Id]; ok {
				if item.Parameter.IsSelect {
					result = utils.CompareStrings(val, item.Operation, item.Value)
				} else {
					result = utils.CompareInts(val, item.Operation, item.Value)
				}
			} else {
				result = false
			}

			if (conditional.IsAnd && !result) || (!conditional.IsAnd && result) {
				break
			}
		}

		if result {
			// Выполняем что надо сделать
		}
	}
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

	// Load items
	for _, conditional := range conditionals {
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

	}

	return conditionals, nil
}
