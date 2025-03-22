package grl

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"grule-protobuf-dsl/dsl"
)

var grlFieldToInputEnum = map[string]dsl.EcommerceOfferRule_Condition_InputField{
	"Customer.Age":                     dsl.EcommerceOfferRule_Condition_AGE,
	"Customer.Gender":                  dsl.EcommerceOfferRule_Condition_GENDER,
	"Customer.Location":                dsl.EcommerceOfferRule_Condition_LOCATION,
	"Customer.DeviceType":              dsl.EcommerceOfferRule_Condition_DEVICE_TYPE,
	"Customer.IsLoyaltyProgramMember":  dsl.EcommerceOfferRule_Condition_IS_LOYALTY_PROGRAM_MEMBER,
	"Customer.TotalSpent":              dsl.EcommerceOfferRule_Condition_TOTAL_LIFETIME_SPENT,
	"Customer.AvgOrderValue":           dsl.EcommerceOfferRule_Condition_AVG_ORDER_VALUE,
	"Customer.LastPurchaseDaysAgo":     dsl.EcommerceOfferRule_Condition_LAST_PURCHASE_DAYS_AGO,
	"Customer.LastCategoryPurchased":   dsl.EcommerceOfferRule_Condition_LAST_CATEGORY_PURCHASED,
	"Customer.PreferredCategories":     dsl.EcommerceOfferRule_Condition_PREFERRED_CATEGORIES,
	"Customer.CartTotal":               dsl.EcommerceOfferRule_Condition_CART_TOTAL,
	"Customer.CartContainsCategories":  dsl.EcommerceOfferRule_Condition_CART_CONTAINS_CATEGORIES,
	"Customer.BrowsingCategories":      dsl.EcommerceOfferRule_Condition_BROWSING_CATEGORIES,
	"Customer.PurchaseCount30d":        dsl.EcommerceOfferRule_Condition_PURCHASE_COUNT_LAST_30_DAYS,
	"Customer.ReturnRatePercent":       dsl.EcommerceOfferRule_Condition_RETURN_RATE_PERCENT,
	"Customer.HasRedeemedCouponBefore": dsl.EcommerceOfferRule_Condition_HAS_COUPON_REDEEMED_BEFORE,
	"Customer.SignupDaysAgo":           dsl.EcommerceOfferRule_Condition_SIGNUP_DAYS_AGO,
}

var outputFieldToGRLName = map[dsl.EcommerceOfferRule_Action_OutputField]string{
	dsl.EcommerceOfferRule_Action_APPLY_DISCOUNT_PERCENT: "Offer.ApplyDiscountPercent",
	dsl.EcommerceOfferRule_Action_APPLY_FLAT_DISCOUNT:    "Offer.ApplyFlatDiscount",
	dsl.EcommerceOfferRule_Action_SHOW_PROMOTION_ID:      "Offer.ShowPromotionId",
	dsl.EcommerceOfferRule_Action_FREE_SHIPPING:          "Offer.FreeShipping",
	dsl.EcommerceOfferRule_Action_ASSIGN_COUPON_CODE:     "Offer.AssignCoupon",
	dsl.EcommerceOfferRule_Action_PROMO_MESSAGE:          "Offer.PromoMessage",
	dsl.EcommerceOfferRule_Action_ADD_LOYALTY_POINTS:     "Offer.AddLoyaltyPoints",
}

// ParseGRLToRuleEntity parses a GRL string into an EcommerceOfferRule proto
func ParseGRLToRuleEntity(grl string) (*dsl.EcommerceOfferRule, error) {
	grl = strings.TrimSpace(grl)
	if !strings.HasPrefix(grl, "rule") {
		return nil, errors.New("invalid GRL: must start with 'rule'")
	}

	nameRegex := regexp.MustCompile(`(?m)^rule\s+(\S+)\s+\"(.*?)\"\s+salience\s+(\d+)\s+\{`)
	whenRegex := regexp.MustCompile(`(?s)when\s+(.*?)then`)
	thenRegex := regexp.MustCompile(`(?s)then\s+(.*?)\}`)

	nameMatch := nameRegex.FindStringSubmatch(grl)
	if len(nameMatch) != 4 {
		return nil, errors.New("could not extract rule name/description/salience")
	}

	name := nameMatch[1]
	description := nameMatch[2]
	salience, _ := strconv.Atoi(nameMatch[3])

	whenMatch := whenRegex.FindStringSubmatch(grl)
	thenMatch := thenRegex.FindStringSubmatch(grl)
	if len(whenMatch) < 2 || len(thenMatch) < 2 {
		return nil, errors.New("could not extract when/then clauses")
	}

	whenClause := strings.TrimSpace(whenMatch[1])
	thenClause := strings.TrimSpace(thenMatch[1])

	rule := &dsl.EcommerceOfferRule{
		Name:        name,
		Description: description,
		Salience:    uint32(salience),
	}

	// Parse WHEN clause
	condition := &dsl.EcommerceOfferRule_Condition{
		ExpressionJoinOperator: dsl.GRuleJoinOperator_AND, // Default
	}
	expressions := strings.Split(whenClause, "&&")
	for _, expr := range expressions {
		expr = strings.Trim(expr, " ()")
		for grlField, inputEnum := range grlFieldToInputEnum {
			if strings.Contains(expr, grlField) {
				operator := detectOperator(expr)
				val := extractValue(expr, operator)
				condition.Expressions = append(condition.Expressions, &dsl.EcommerceOfferRule_Condition_Expression{
					Input:    inputEnum,
					Operator: operator,
					Value:    parseValue(val),
				})
				break
			}
		}
	}
	if len(condition.Expressions) > 0 {
		rule.Conditions = append(rule.Conditions, condition)
	}

	// Parse THEN clause
	thenLines := strings.Split(thenClause, ";")
	for _, line := range thenLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		left := strings.TrimSpace(parts[0])
		right := strings.TrimSpace(parts[1])

		for output, grlName := range outputFieldToGRLName {
			if grlName == left {
				rule.Actions = append(rule.Actions, &dsl.EcommerceOfferRule_Action{
					Output: output,
					Value:  parseValue(right),
				})
				break
			}
		}
	}

	return rule, nil
}

func detectOperator(expr string) dsl.GRuleExpressionOperator {
	switch {
	case strings.Contains(expr, " >= "):
		return dsl.GRuleExpressionOperator_GREATER_THAN_EQUALS
	case strings.Contains(expr, " <= "):
		return dsl.GRuleExpressionOperator_LESS_THAN_EQUALS
	case strings.Contains(expr, " > "):
		return dsl.GRuleExpressionOperator_GREATER_THAN
	case strings.Contains(expr, " < "):
		return dsl.GRuleExpressionOperator_LESS_THAN
	case strings.Contains(expr, " == "):
		return dsl.GRuleExpressionOperator_EQUALS
	case strings.Contains(expr, " != "):
		return dsl.GRuleExpressionOperator_NOT_EQUALS
	default:
		return dsl.GRuleExpressionOperator_EXPRESSION_OPERATOR_UNSPECIFIED
	}
}

func extractValue(expr string, op dsl.GRuleExpressionOperator) string {
	parts := strings.Split(expr, getGrlOperator(op))
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func parseValue(val string) *dsl.RuleValue {
	val = strings.Trim(val, "\" ")
	if b, err := strconv.ParseBool(val); err == nil {
		return &dsl.RuleValue{Value: &dsl.RuleValue_BoolVal{BoolVal: b}}
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return &dsl.RuleValue{Value: &dsl.RuleValue_FloatVal{FloatVal: float32(f)}}
	}
	if i, err := strconv.Atoi(val); err == nil {
		return &dsl.RuleValue{Value: &dsl.RuleValue_IntVal{IntVal: int32(i)}}
	}
	return &dsl.RuleValue{Value: &dsl.RuleValue_StringVal{StringVal: val}}
}

func getGrlOperator(op dsl.GRuleExpressionOperator) string {
	switch op {
	case dsl.GRuleExpressionOperator_LESS_THAN:
		return " < "
	case dsl.GRuleExpressionOperator_LESS_THAN_EQUALS:
		return " <= "
	case dsl.GRuleExpressionOperator_GREATER_THAN:
		return " > "
	case dsl.GRuleExpressionOperator_GREATER_THAN_EQUALS:
		return " >= "
	case dsl.GRuleExpressionOperator_EQUALS:
		return " == "
	case dsl.GRuleExpressionOperator_NOT_EQUALS:
		return " != "
	default:
		return ""
	}
}
