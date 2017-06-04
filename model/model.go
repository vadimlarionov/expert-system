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
	IsSelect  bool
	Parameter *Parameter `orm:"rel(one)"`
	Answers   []*Answer  `orm:"reverse(many)"`
}

type Answer struct {
	Id           uint
	Question     *Question `orm:"rel(fk)"`
	NextQuestion *Question `orm:"rel(fk)"`
	Value        string
}

type Condition struct {
	Id    uint
	IsAnd bool
	Items []*ConditionItem `orm:"reverse(many)"`
}

type ConditionItem struct {
	Id        uint
	Condition *Condition `orm:"rel(fk)"`
	Parameter *Parameter `orm:"rel(fk)"`
	Operation string     `orm:"size(10)"`
	Value     string
}

type ConditionResult struct {
	Id               uint
	AttributesResult []*ConditionAttributeResult   `orm:"reverse(many)"`
	ParametersResult []*ConditionalParameterResult `orm:"reverse(many)"`
}

type ConditionAttributeResult struct {
	Id              uint
	ConditionResult *ConditionResult `orm:"rel(fk)"`
	Attribute       *Attribute       `orm:"rel(fk)"`
}

type ConditionalParameterResult struct {
	Id              uint
	ConditionResult *ConditionResult `orm:"rel(fk)"`
	Parameter       *Parameter       `orm:"rel(fk)"`
}
