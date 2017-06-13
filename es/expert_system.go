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

	expertSystem, err := initExpertSystem(o)
	parametersMap := make(map[uint]string)
	attributesMap := make(map[string]string)

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

			parametersMap[q.Parameter.Id] = answer.Value.Value
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

		parametersMapChanged := true
		for parametersMapChanged {
			parametersMapChanged = checkConditionalsParam(parametersMap, expertSystem.conditionals)
		}

		attributesMap = checkConditionalsAttr(parametersMap, expertSystem.conditionals)

		expertSystem.printObjects(attributesMap)
	}

	expertSystem.printObjects(attributesMap)
}

func (expertSystem *expertSystemType) printObjects(attributesMap map[string]string) {
	ratingMap := make(map[string]int)
	for _, obj := range expertSystem.objects {
		ratingMap[obj.Name] = 0
		for _, attrVal := range obj.AttributeValues {
			if val, ok := attributesMap[attrVal.Attribute.Text]; ok {
				if val == attrVal.Text {
					ratingMap[obj.Name] += 1
				}
			}
		}
	}

	utils.PrintObjectsWithRating(ratingMap)
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

func checkConditionalsAttr(parametersMap map[uint]string, conditionals []*model.Conditional) map[string]string {
	attributesMap := make(map[string]string)
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
			for _, attrResult := range conditional.AttributeResults {
				attributesMap[attrResult.Attribute.Text] = attrResult.AttributeValue.Text
			}
		}
	}
	return attributesMap
}

func checkConditionalsParam(parametersMap map[uint]string, conditionals []*model.Conditional) (changed bool) {
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
			for _, paramResult := range conditional.ParameterResults {
				if _, exists := parametersMap[paramResult.Parameter.Id]; !exists {
					parametersMap[paramResult.Parameter.Id] = paramResult.Value
					changed = true
				}
			}
		}
	}
	return changed
}
