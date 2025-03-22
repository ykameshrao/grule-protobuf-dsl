package grl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"grule-protobuf-dsl/dsl"
	"grule-protobuf-dsl/grl"
)

func TestEcommerceOfferRuleToGRuleEntity_SimpleRule(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "SimpleDiscountRule",
		Description: "Apply discount if cart total > 1000",
		Salience:    10,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_CART_TOTAL,
						Operator: dsl.GRuleExpressionOperator_GREATER_THAN,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_FloatVal{FloatVal: 1000.00},
						},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_JOIN_OPERATOR_UNSPECIFIED,
			},
		},
		ConditionJoinOperator: dsl.GRuleJoinOperator_JOIN_OPERATOR_UNSPECIFIED,
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_APPLY_DISCOUNT_PERCENT,
				Value: &dsl.RuleValue{
					Value: &dsl.RuleValue_FloatVal{FloatVal: 10.0},
				},
			},
		},
	}

	entity, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.NoError(t, err)
	assert.Equal(t, "SimpleDiscountRule", entity.Name)
	assert.Equal(t, "Apply discount if cart total > 1000", entity.Description)
	assert.Equal(t, "10", entity.Salience)
	assert.Contains(t, entity.When, "Customer.CartTotal > 1000.00")
	assert.Contains(t, entity.Then[0], "Offer.ApplyDiscountPercent = 10.00;")
}

func TestEcommerceOfferRuleToGRuleEntity_MissingConditions(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "MissingConditions",
		Description: "Should error when no conditions",
		Salience:    1,
		Conditions:  []*dsl.EcommerceOfferRule_Condition{},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_APPLY_FLAT_DISCOUNT,
				Value:  &dsl.RuleValue{Value: &dsl.RuleValue_FloatVal{FloatVal: 50}},
			},
		},
	}

	_, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.Error(t, err)
	assert.Equal(t, "no conditions defined", err.Error())
}

func TestEcommerceOfferRuleToGRuleEntity_MissingActions(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "MissingActions",
		Description: "Should error when no actions",
		Salience:    1,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_CART_TOTAL,
						Operator: dsl.GRuleExpressionOperator_GREATER_THAN,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_FloatVal{FloatVal: 500.00},
						},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
			},
		},
		Actions: nil,
	}

	_, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.Error(t, err)
	assert.Equal(t, "no actions defined", err.Error())
}

func TestEcommerceOfferRuleToGRuleEntity_StringListValue(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "CategoryMatch",
		Description: "Applies if preferred categories match",
		Salience:    5,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_PREFERRED_CATEGORIES,
						Operator: dsl.GRuleExpressionOperator_HAS_CATEGORY_FUNCTION,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_StringListCommaConcatenated{StringListCommaConcatenated: "Electronics, Home"},
						},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_JOIN_OPERATOR_UNSPECIFIED,
			},
		},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_PROMO_MESSAGE,
				Value: &dsl.RuleValue{
					Value: &dsl.RuleValue_StringVal{StringVal: "Special Deal!"},
				},
			},
		},
	}

	entity, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.NoError(t, err)
	assert.Contains(t, entity.When, "Customer.HasCategory(Customer.PreferredCategories, \"Electronics\", \"Home\")")
	assert.Contains(t, entity.Then[0], "Offer.PromoMessage = \"Special Deal!\";")
}

func TestEcommerceOfferRuleToGRuleEntity_MultipleConditionsAndActions(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "MultiConditionMultiAction",
		Description: "Apply discount and show promo if age > 25 and is loyalty member",
		Salience:    15,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_AGE,
						Operator: dsl.GRuleExpressionOperator_GREATER_THAN,
						Value:    &dsl.RuleValue{Value: &dsl.RuleValue_IntVal{IntVal: 25}},
					},
					{
						Input:    dsl.EcommerceOfferRule_Condition_IS_LOYALTY_PROGRAM_MEMBER,
						Operator: dsl.GRuleExpressionOperator_EQUALS,
						Value:    &dsl.RuleValue{Value: &dsl.RuleValue_BoolVal{BoolVal: true}},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
			},
		},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_APPLY_FLAT_DISCOUNT,
				Value:  &dsl.RuleValue{Value: &dsl.RuleValue_FloatVal{FloatVal: 200.0}},
			},
			{
				Output: dsl.EcommerceOfferRule_Action_SHOW_PROMOTION_ID,
				Value:  &dsl.RuleValue{Value: &dsl.RuleValue_StringVal{StringVal: "LOYAL25"}},
			},
		},
	}

	entity, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.NoError(t, err)
	assert.Contains(t, entity.When, "Customer.Age > 25")
	assert.Contains(t, entity.When, "Customer.IsLoyaltyProgramMember == true")
	assert.Contains(t, entity.Then[0], "Offer.ApplyFlatDiscount = 200.00;")
	assert.Contains(t, entity.Then[1], `Offer.ShowPromotionId = "LOYAL25";`)
}

func TestEcommerceOfferRuleToGRuleEntity_UsingBoolCondition(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "FreeShippingLoyalty",
		Description: "Free shipping for loyalty members",
		Salience:    20,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_IS_LOYALTY_PROGRAM_MEMBER,
						Operator: dsl.GRuleExpressionOperator_EQUALS,
						Value:    &dsl.RuleValue{Value: &dsl.RuleValue_BoolVal{BoolVal: true}},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_JOIN_OPERATOR_UNSPECIFIED,
			},
		},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_FREE_SHIPPING,
				Value:  &dsl.RuleValue{Value: &dsl.RuleValue_BoolVal{BoolVal: true}},
			},
		},
	}

	entity, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.NoError(t, err)
	assert.Contains(t, entity.When, "Customer.IsLoyaltyProgramMember == true")
	assert.Contains(t, entity.Then[0], "Offer.FreeShipping = true;")
}

func TestEcommerceOfferRuleToGRuleEntity_ComplexStringList(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "CategoryAndCouponTest",
		Description: "Apply coupon if browsing categories match",
		Salience:    12,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_BROWSING_CATEGORIES,
						Operator: dsl.GRuleExpressionOperator_HAS_CATEGORY_FUNCTION,
						Value:    &dsl.RuleValue{Value: &dsl.RuleValue_StringListCommaConcatenated{StringListCommaConcatenated: "Books, Gadgets"}},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_JOIN_OPERATOR_UNSPECIFIED,
			},
		},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_ASSIGN_COUPON_CODE,
				Value:  &dsl.RuleValue{Value: &dsl.RuleValue_StringVal{StringVal: "BOOKS10"}},
			},
		},
	}

	entity, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.NoError(t, err)
	assert.Contains(t, entity.When, `Customer.HasCategory(Customer.BrowsingCategories, "Books", "Gadgets")`)
	assert.Contains(t, entity.Then[0], `Offer.AssignCoupon = "BOOKS10";`)
}

func TestEcommerceOfferRuleToGRuleEntity_MultipleStringInFunctions(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "MultiCategoryMatch",
		Description: "Match on both preferred and browsing categories",
		Salience:    8,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_PREFERRED_CATEGORIES,
						Operator: dsl.GRuleExpressionOperator_HAS_CATEGORY_FUNCTION,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_StringListCommaConcatenated{StringListCommaConcatenated: "Electronics, Furniture"},
						},
					},
					{
						Input:    dsl.EcommerceOfferRule_Condition_BROWSING_CATEGORIES,
						Operator: dsl.GRuleExpressionOperator_HAS_CATEGORY_FUNCTION,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_StringListCommaConcatenated{StringListCommaConcatenated: "Books, Toys"},
						},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
			},
		},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_PROMO_MESSAGE,
				Value: &dsl.RuleValue{
					Value: &dsl.RuleValue_StringVal{StringVal: "Combo Promo!"},
				},
			},
		},
	}

	entity, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.NoError(t, err)
	assert.Contains(t, entity.When, `Customer.HasCategory(Customer.PreferredCategories, "Electronics", "Furniture")`)
	assert.Contains(t, entity.When, `Customer.HasCategory(Customer.BrowsingCategories, "Books", "Toys")`)
	assert.Contains(t, entity.Then[0], `Offer.PromoMessage = "Combo Promo!";`)
}

func TestEcommerceOfferRuleToGRuleEntity_StringInFunctionMixedWithOther(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "CategoryAndCartCombo",
		Description: "Apply if cart total is high and categories match",
		Salience:    11,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_CART_TOTAL,
						Operator: dsl.GRuleExpressionOperator_GREATER_THAN_EQUALS,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_FloatVal{FloatVal: 1500},
						},
					},
					{
						Input:    dsl.EcommerceOfferRule_Condition_CART_CONTAINS_CATEGORIES,
						Operator: dsl.GRuleExpressionOperator_HAS_CATEGORY_FUNCTION,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_StringListCommaConcatenated{StringListCommaConcatenated: "Appliances"},
						},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
			},
		},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_APPLY_FLAT_DISCOUNT,
				Value: &dsl.RuleValue{
					Value: &dsl.RuleValue_FloatVal{FloatVal: 300},
				},
			},
		},
	}

	entity, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.NoError(t, err)
	assert.Contains(t, entity.When, `Customer.CartTotal >= 1500.00`)
	assert.Contains(t, entity.When, `Customer.HasCategory(Customer.CartContainsCategories, "Appliances")`)
	assert.Contains(t, entity.Then[0], `Offer.ApplyFlatDiscount = 300.00;`)
}

func TestEcommerceOfferRuleToGRuleEntity_EmptyStringList(t *testing.T) {
	rule := &dsl.EcommerceOfferRule{
		Name:        "EmptyCategoryCheck",
		Description: "Should error or skip if no categories",
		Salience:    4,
		Conditions: []*dsl.EcommerceOfferRule_Condition{
			{
				Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
					{
						Input:    dsl.EcommerceOfferRule_Condition_CART_CONTAINS_CATEGORIES,
						Operator: dsl.GRuleExpressionOperator_HAS_CATEGORY_FUNCTION,
						Value:    &dsl.RuleValue{Value: &dsl.RuleValue_StringListCommaConcatenated{StringListCommaConcatenated: ""}},
					},
				},
				ExpressionJoinOperator: dsl.GRuleJoinOperator_JOIN_OPERATOR_UNSPECIFIED,
			},
		},
		Actions: []*dsl.EcommerceOfferRule_Action{
			{
				Output: dsl.EcommerceOfferRule_Action_SHOW_PROMOTION_ID,
				Value:  &dsl.RuleValue{Value: &dsl.RuleValue_StringVal{StringVal: "EMPTYCAT"}},
			},
		},
	}

	_, err := grl.EcommerceOfferRuleToGRuleEntity(rule)
	assert.Error(t, err)
	assert.Equal(t, `HAS_CATEGORY_FUNCTION used with empty list for field Customer.CartContainsCategories`, err.Error())
}
