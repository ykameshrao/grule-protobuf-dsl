package grl_test

import (
	"testing"

	"grule-protobuf-dsl/dsl"
	"grule-protobuf-dsl/grl"

	"github.com/stretchr/testify/assert"
)

func TestParseGRLToRuleEntity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *dsl.EcommerceOfferRule
		wantErr  bool
	}{
		{
			name: "Simple rule with one condition and one action",
			input: `rule ApplyDiscount "Apply 10% discount if cart total > 500" salience 10 {
	when
		(Customer.CartTotal > 500)
	then
		Offer.ApplyDiscountPercent = 10.00;
}`,
			expected: &dsl.EcommerceOfferRule{
				Name:        "ApplyDiscount",
				Description: "Apply 10% discount if cart total > 500",
				Salience:    10,
				Conditions: []*dsl.EcommerceOfferRule_Condition{
					{
						ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
						Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
							{
								Input:    dsl.EcommerceOfferRule_Condition_CART_TOTAL,
								Operator: dsl.GRuleExpressionOperator_GREATER_THAN,
								Value: &dsl.RuleValue{
									Value: &dsl.RuleValue_FloatVal{FloatVal: 500.00},
								},
							},
						},
					},
				},
				Actions: []*dsl.EcommerceOfferRule_Action{
					{
						Output: dsl.EcommerceOfferRule_Action_APPLY_DISCOUNT_PERCENT,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_FloatVal{FloatVal: 10.00},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple conditions with AND and multiple actions",
			input: `rule OfferWithMultipleConditions "Show promotion if loyal and cart value > 300" salience 15 {
	when
		(Customer.IsLoyaltyProgramMember == true) && (Customer.CartTotal > 300)
	then
		Offer.ShowPromotionId = "LOYAL300";
		Offer.FreeShipping = true;
}`,
			expected: &dsl.EcommerceOfferRule{
				Name:        "OfferWithMultipleConditions",
				Description: "Show promotion if loyal and cart value > 300",
				Salience:    15,
				Conditions: []*dsl.EcommerceOfferRule_Condition{
					{
						ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
						Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
							{
								Input:    dsl.EcommerceOfferRule_Condition_IS_LOYALTY_PROGRAM_MEMBER,
								Operator: dsl.GRuleExpressionOperator_EQUALS,
								Value: &dsl.RuleValue{
									Value: &dsl.RuleValue_BoolVal{BoolVal: true},
								},
							},
							{
								Input:    dsl.EcommerceOfferRule_Condition_CART_TOTAL,
								Operator: dsl.GRuleExpressionOperator_GREATER_THAN,
								Value: &dsl.RuleValue{
									Value: &dsl.RuleValue_FloatVal{FloatVal: 300},
								},
							},
						},
					},
				},
				Actions: []*dsl.EcommerceOfferRule_Action{
					{
						Output: dsl.EcommerceOfferRule_Action_SHOW_PROMOTION_ID,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_StringVal{StringVal: "LOYAL300"},
						},
					},
					{
						Output: dsl.EcommerceOfferRule_Action_FREE_SHIPPING,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_BoolVal{BoolVal: true},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "String match condition and assign coupon",
			input: `rule AssignCoupon "Give coupon if last category is electronics" salience 8 {
	when
		(Customer.LastCategoryPurchased == "electronics")
	then
		Offer.AssignCoupon = "ELEC10";
}`,
			expected: &dsl.EcommerceOfferRule{
				Name:        "AssignCoupon",
				Description: "Give coupon if last category is electronics",
				Salience:    8,
				Conditions: []*dsl.EcommerceOfferRule_Condition{
					{
						ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
						Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
							{
								Input:    dsl.EcommerceOfferRule_Condition_LAST_CATEGORY_PURCHASED,
								Operator: dsl.GRuleExpressionOperator_EQUALS,
								Value: &dsl.RuleValue{
									Value: &dsl.RuleValue_StringVal{StringVal: "electronics"},
								},
							},
						},
					},
				},
				Actions: []*dsl.EcommerceOfferRule_Action{
					{
						Output: dsl.EcommerceOfferRule_Action_ASSIGN_COUPON_CODE,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_StringVal{StringVal: "ELEC10"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Integer comparison and loyalty points action",
			input: `rule AddPoints "Reward frequent buyers" salience 12 {
	when
		(Customer.PurchaseCount30d > 5)
	then
		Offer.AddLoyaltyPoints = 100;
}`,
			expected: &dsl.EcommerceOfferRule{
				Name:        "AddPoints",
				Description: "Reward frequent buyers",
				Salience:    12,
				Conditions: []*dsl.EcommerceOfferRule_Condition{
					{
						ExpressionJoinOperator: dsl.GRuleJoinOperator_AND,
						Expressions: []*dsl.EcommerceOfferRule_Condition_Expression{
							{
								Input:    dsl.EcommerceOfferRule_Condition_PURCHASE_COUNT_LAST_30_DAYS,
								Operator: dsl.GRuleExpressionOperator_GREATER_THAN,
								Value: &dsl.RuleValue{
									Value: &dsl.RuleValue_IntVal{IntVal: 100},
								},
							},
						},
					},
				},
				Actions: []*dsl.EcommerceOfferRule_Action{
					{
						Output: dsl.EcommerceOfferRule_Action_ADD_LOYALTY_POINTS,
						Value: &dsl.RuleValue{
							Value: &dsl.RuleValue_IntVal{IntVal: 100},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Missing when clause",
			input:   `rule MissingWhen "No when" salience 1 {then Offer.ShowPromotionId = "PROMO123";}`,
			wantErr: true,
		},
		{
			name:    "Missing then clause",
			input:   `rule MissingThen "No then" salience 1 {when Customer.Age > 18}`,
			wantErr: true,
		},
		{
			name:    "Invalid syntax",
			input:   `rule InvalidSyntax "Broken" salience x {}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rule, err := grl.ParseGRLToRuleEntity(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Name, rule.Name)
				assert.Equal(t, tt.expected.Description, rule.Description)
				assert.Equal(t, tt.expected.Salience, rule.Salience)
				assert.Len(t, rule.Conditions, len(tt.expected.Conditions))
				assert.Len(t, rule.Actions, len(tt.expected.Actions))
			}
		})
	}
}
