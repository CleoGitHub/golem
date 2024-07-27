package coredomaindefinition

type ValidationRule string

const (
	ValidationRuleRequired ValidationRule = "required"
	ValidationRuleEmail    ValidationRule = "email"
	ValidationRuleGT       ValidationRule = "gt"
	ValidationRuleGTE      ValidationRule = "gte"
	ValidationRuleLT       ValidationRule = "lt"
	ValidationRuleLTE      ValidationRule = "lte"
	ValidationRuleUUID     ValidationRule = "uuid"
	ValidationRuleHexColor ValidationRule = "hexcolor"
	ValidationRuleUnique   ValidationRule = "unique"
	ValidationRuleUniqueIn ValidationRule = "uniqueIn"
)

type Validation struct {
	Rule  ValidationRule
	Value interface{}
}
