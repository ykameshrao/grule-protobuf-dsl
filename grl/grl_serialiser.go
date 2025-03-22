package grl

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"grule-protobuf-dsl/dsl"
)

// EcommerceOfferRuleToGRuleEntity converts an EcommerceOfferRule to a GRuleEntity
func EcommerceOfferRuleToGRuleEntity(rule *dsl.EcommerceOfferRule) (*GRuleEntity, error) {
	if len(rule.Conditions) == 0 {
		return nil, fmt.Errorf("no conditions defined")
	}
	if len(rule.Actions) == 0 {
		return nil, fmt.Errorf("no actions defined")
	}

	// Parse conditions to GRL 'when' clause
	conditions := make([]string, 0)
	for _, cond := range rule.Conditions {
		expressions := make([]string, 0)
		for _, expr := range cond.Expressions {
			val, err := getRuleValue(expr.Value)
			if err != nil {
				return nil, err
			}
			op := getEnumGrlOperator(expr.Operator)
			field := getEnumGrlFieldName(expr.Input)
			if expr.Operator == dsl.GRuleExpressionOperator_STRING_IN_FUNCTION {
				if val == "" {
					return nil, fmt.Errorf("STRING_IN_FUNCTION used with empty list for field %s", field)
				}
				exprStr := strings.Replace(fmt.Sprintf("%s%s", field, op), ":replace", val, 1)
				expressions = append(expressions, fmt.Sprintf("( %s )", exprStr))
			} else {
				expressions = append(expressions, fmt.Sprintf("( %s%s%s )", field, op, val))
			}
		}
		joined := strings.Join(expressions, getEnumGrlOperator(cond.ExpressionJoinOperator))
		conditions = append(conditions, joined)
	}
	when := strings.Join(conditions, getEnumGrlOperator(rule.ConditionJoinOperator))

	// Parse actions to GRL 'then' clause
	then := make([]string, 0, len(rule.Actions))
	for _, action := range rule.Actions {
		val, err := getRuleValue(action.Value)
		if err != nil {
			return nil, err
		}
		actionStr := fmt.Sprintf("%s = %s;", getEnumGrlFieldName(action.Output), val)
		then = append(then, actionStr)
	}

	return &GRuleEntity{
		Name:        rule.Name,
		Description: rule.Description,
		Salience:    strconv.Itoa(int(rule.Salience)),
		When:        when,
		Then:        then,
	}, nil
}

func getRuleValue(val *dsl.RuleValue) (string, error) {
	switch v := val.Value.(type) {
	case *dsl.RuleValue_StringVal:
		return strconv.Quote(v.StringVal), nil
	case *dsl.RuleValue_BoolVal:
		return strconv.FormatBool(v.BoolVal), nil
	case *dsl.RuleValue_IntVal:
		return strconv.Itoa(int(v.IntVal)), nil
	case *dsl.RuleValue_FloatVal:
		return fmt.Sprintf("%.2f", v.FloatVal), nil
	case *dsl.RuleValue_StringListCommaConcatenated:
		elements := strings.Split(v.StringListCommaConcatenated, ",")
		quoted := make([]string, 0, len(elements))
		for _, e := range elements {
			e = strings.TrimSpace(e)
			if e != "" {
				quoted = append(quoted, strconv.Quote(e))
			}
		}
		return strings.Join(quoted, ", "), nil
	default:
		return "", fmt.Errorf("unsupported rule value type")
	}
}

func getEnumGrlFieldName(enum interface{ protoreflect.Enum }) string {
	fieldName := proto.GetExtension(enum.Descriptor().Values().ByNumber(enum.Number()).Options(), dsl.E_GrlFieldName)
	return fieldName.(string)
}

func getEnumGrlOperator(enum interface{ protoreflect.Enum }) string {
	fieldName := proto.GetExtension(enum.Descriptor().Values().ByNumber(enum.Number()).Options(), dsl.E_GrlOperator)
	return fieldName.(string)
}

func getEnumGrlFieldType(enum interface{ protoreflect.Enum }) string {
	fieldName := proto.GetExtension(enum.Descriptor().Values().ByNumber(enum.Number()).Options(), dsl.E_GrlFieldType)
	return fieldName.(string)
}

// ToGRL converts a GRuleEntity to a GRL string
func ToGRL(grule *GRuleEntity) string {
	return fmt.Sprintf(`rule %s "%s" salience %s {
	when
		%s
	then
		%s
}`,
		grule.Name, grule.Description, grule.Salience, grule.When, strings.Join(grule.Then, "\n\t\t"))
}

// ToMultipleGRLs converts a slice of GRuleEntity to a GRL string
func ToMultipleGRLs(rules []*GRuleEntity) string {
	var sb strings.Builder
	for _, rule := range rules {
		sb.WriteString(ToGRL(rule))
		sb.WriteString("\n")
	}
	return sb.String()
}
