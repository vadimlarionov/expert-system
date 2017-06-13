package model

type Attribute struct {
	Id     uint
	Text   string            `orm:"unique"`
	Values []*AttributeValue `orm:"reverse(many)"`
}

type AttributeValue struct {
	Id        uint
	Attribute *Attribute `orm:"rel(fk)"`
	Text      string
}

type Object struct {
	Id              uint
	Name            string            `orm:"unique"`
	AttributeValues []*AttributeValue `orm:"rel(m2m)"`
}

type Parameter struct {
	Id       uint
	Name     string `orm:"unique"`
	IsSelect bool
	Values   []*ParameterValue `orm:"reverse(many)"`
}

type ParameterValue struct {
	Id        uint
	Parameter *Parameter `orm:"rel(fk)"`
	Value     string
}

type Question struct {
	Id        uint
	Text      string
	IsSelect  bool
	Number    int        `orm:"unique"`
	Parameter *Parameter `orm:"rel(one)"`
	Answers   []*Answer  `orm:"reverse(many)"`
}

type Answer struct {
	Id                 uint
	Question           *Question `orm:"rel(fk)"`
	NextQuestionNumber int
	Text               string
	Value              *ParameterValue `orm:"rel(fk)"`
}

type Conditional struct {
	Id               uint
	IsAnd            bool
	Items            []*ConditionalItem            `orm:"reverse(many)"`
	AttributeResults []*ConditionalAttributeResult `orm:"reverse(many)"`
}

type ConditionalItem struct {
	Id        uint
	Condition *Conditional `orm:"rel(fk)"`
	Parameter *Parameter   `orm:"rel(fk)"`
	Operation string       `orm:"size(10)"`
	Value     string
}

type ConditionalParameterResult struct {
	Id          uint
	Conditional *Conditional `orm:"rel(fk)"`
	Parameter   *Parameter   `orm:"rel(fk)"`
	Value       string
}

type ConditionalAttributeResult struct {
	Id             uint
	Conditional    *Conditional    `orm:"rel(fk)"`
	Attribute      *Attribute      `orm:"rel(fk)"`
	AttributeValue *AttributeValue `orm:"rel(fk)"`
}

type Quest struct {
	Id       uint
	Username string
}

type QuestAttribute struct {
	Id             uint
	Quest          *Quest          `orm:"rel(fk)"`
	Attribute      *Attribute      `orm:"rel(fk)"`
	AttributeValue *AttributeValue `orm:"rel(fk)"`
}

type QuestParameter struct {
	Id        uint
	Quest     *Quest     `orm:"rel(fk)"`
	Parameter *Parameter `orm:"rel(fk)"`
	Value     string
}

type QuestQuestions struct {
	Id         uint
	Quest      *Quest    `orm:"rel(fk)"`
	Question   *Question `orm:"rel(fk)"`
	IsAnswered bool
}
