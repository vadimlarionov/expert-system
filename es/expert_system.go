package es

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/vadimlarionov/expert-system/model"
	"github.com/vadimlarionov/expert-system/utils"
	"strconv"
	"strings"
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

			answerIndex := readAnswerWithChoice(expertSystem, q.Answers)
			if answerIndex >= 0 {
				answer := q.Answers[answerIndex]
				if _, err = o.LoadRelated(answer, "Value"); err != nil {
					fmt.Printf("Can't load related parameter value for answer \"%s\": %s\n",
						answer.Text, err)
					return
				}

				expertSystem.parametersMapState[q.Parameter.Id] = answer.Value.Value

				questionNum = answer.NextQuestionNumber
			} else {
				questionNum++
			}

		} else {
			userAnswer, skip := readFreeAnswer(expertSystem)
			if !skip {
				expertSystem.parametersMapState[q.Parameter.Id] = strconv.Itoa(userAnswer)
			}
			questionNum++
		}

		parametersMapChanged := true
		for parametersMapChanged {
			parametersMapChanged = checkConditionalsParam(expertSystem.parametersMapState, expertSystem.conditionals)
		}

		expertSystem.attributesMapState = checkConditionalsAttr(expertSystem.parametersMapState, expertSystem.conditionals)
	}

	expertSystem.printObjects()
}

func (expertSystem *expertSystemType) printObjects() {
	ratingMap := make(map[string]int)
	for _, obj := range expertSystem.objects {
		ratingMap[obj.Name] = 0
		for _, attrVal := range obj.AttributeValues {
			if val, ok := expertSystem.attributesMapState[attrVal.Attribute.Text]; ok {
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

func readAnswerWithChoice(expertSystem *expertSystemType, answers []*model.Answer) (answerIndex int) {
	for {
		for i, answer := range answers {
			fmt.Printf("%d: %s\n", i+1, answer.Text)
		}
		fmt.Printf("Ваш ответ: ")

		var userAnswer string
		fmt.Scanf("%s", &userAnswer)

		if userAnswer == "" {
			return -1
		} else if strings.HasPrefix(userAnswer, "/") {
			printSystemMessage(expertSystem, userAnswer)
		} else {
			answerIndex, err := strconv.Atoi(userAnswer)
			if err != nil {
				fmt.Println("Необходимо выбрать один из предложенных вариантов")
				continue
			}

			if answerIndex > 0 && answerIndex <= len(answers) {
				return answerIndex - 1
			} else {
				fmt.Println("Необходимо выбрать один из предложенных вариантов")
			}
		}
	}
}

func readFreeAnswer(expertSystem *expertSystemType) (answer int, skip bool) {
	for {
		var userAnswer string
		fmt.Scanf("%s", &userAnswer)

		if userAnswer == "" {
			return 0, true
		} else if strings.HasPrefix(userAnswer, "/") {
			printSystemMessage(expertSystem, userAnswer)
		} else {
			answer, err := strconv.Atoi(userAnswer)
			if err != nil {
				fmt.Println("Необходимо ввести число")
				continue
			}

			return answer, false
		}
	}
}

func printSystemMessage(expertSystem *expertSystemType, cmd string) {
	switch cmd {
	case "/rating":
		expertSystem.printObjects()
	default:
		fmt.Printf("Unexpected command \"%s\"\n", cmd)
	}
}
