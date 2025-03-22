package main

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"grule-protobuf-dsl/dsl"
	"grule-protobuf-dsl/grl"
)

type Customer struct {
	Age                     int
	Gender                  string
	Location                string
	DeviceType              string
	IsLoyaltyProgramMember  bool
	TotalSpent              float32
	AvgOrderValue           float32
	LastPurchaseDaysAgo     int
	LastCategoryPurchased   string
	PreferredCategories     []string
	CartTotal               float32
	CartContainsCategories  []string
	BrowsingCategories      []string
	PurchaseCount30d        int
	ReturnRatePercent       float32
	HasRedeemedCouponBefore bool
	SignupDaysAgo           int
}

func (c Customer) HasCategory(field []string, categories ...string) bool {
	for _, input := range categories {
		for _, existing := range field {
			if input == existing {
				return true
			}
		}
	}
	return false
}

type Offer struct {
	ApplyDiscountPercent float32
	ApplyFlatDiscount    float32
	ShowPromotionId      string
	FreeShipping         bool
	AssignCoupon         string
	PromoMessage         string
	AddLoyaltyPoints     int
}

type RuleContext struct {
	Customer Customer
	Offer    Offer
}

func loadAllRulesFromDir(dir string) ([]*dsl.EcommerceOfferRule, error) {
	rules := []*dsl.EcommerceOfferRule{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var rule dsl.EcommerceOfferRule
		err = protojson.UnmarshalOptions{
			DiscardUnknown: true,
			AllowPartial:   true,
		}.Unmarshal(data, &rule)
		rules = append(rules, &rule)
		return nil
	})
	return rules, err
}

func setupRuleEngine(grlRules []string) (*ast.KnowledgeBase, error) {
	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)

	for _, ruleStr := range grlRules {
		err := ruleBuilder.BuildRuleFromResource("EcommerceOffersRuleEngine", "0.0.1", pkg.NewBytesResource([]byte(ruleStr)))
		if err != nil {
			return nil, err
		}
	}

	kb, _ := lib.NewKnowledgeBaseInstance("EcommerceOffersRuleEngine", "0.0.1")
	return kb, nil
}

func exampleOne() (ast.IDataContext, *RuleContext, error) {
	// Prepare context
	ruleCtx := &RuleContext{
		Customer: Customer{
			Age:                     30,
			CartTotal:               1500.0,
			PreferredCategories:     []string{"Electronics", "Home"},
			IsLoyaltyProgramMember:  true,
			TotalSpent:              50000,
			SignupDaysAgo:           300,
			HasRedeemedCouponBefore: false,
		},
	}
	dc := ast.NewDataContext()
	err := dc.Add("Customer", &ruleCtx.Customer)
	if err != nil {
		return nil, nil, err
	}
	err = dc.Add("Offer", &ruleCtx.Offer)
	if err != nil {
		return nil, nil, err
	}
	return dc, ruleCtx, nil
}

func exampleTwo() (ast.IDataContext, *RuleContext, error) {
	// Prepare context
	ruleCtx := &RuleContext{
		Customer: Customer{
			Age:                     21,
			BrowsingCategories:      []string{"Electronics", "Home"},
			TotalSpent:              500,
			SignupDaysAgo:           30,
			HasRedeemedCouponBefore: false,
		},
	}
	dc := ast.NewDataContext()
	err := dc.Add("Customer", &ruleCtx.Customer)
	if err != nil {
		return nil, nil, err
	}
	err = dc.Add("Offer", &ruleCtx.Offer)
	if err != nil {
		return nil, nil, err
	}
	return dc, ruleCtx, nil
}

func main() {
	// Step 1: Load rules from disk
	rules, err := loadAllRulesFromDir("rules")
	if err != nil {
		panic(err)
	}

	// Step 2: Convert to GRL format
	var grlRules []string
	for _, r := range rules {
		entity, err := grl.EcommerceOfferRuleToGRuleEntity(r)
		if err != nil {
			panic(err)
		}
		grlRules = append(grlRules, grl.ToGRL(entity))
		fmt.Println("Loaded GRule:\n", grl.ToGRL(entity))
	}

	// Step 3: Load into engine and create context
	kb, err := setupRuleEngine(grlRules)
	if err != nil {
		panic(err)
	}

	// Step 4: Evaluate
	dc, ruleCtx, err := exampleOne()
	e := engine.NewGruleEngine()
	matchingRules, err := e.FetchMatchingRules(dc, kb)
	if err != nil {
		return
	}
	for _, rule := range matchingRules {
		fmt.Println("Matching Rule: ", rule)
	}
	if err := e.Execute(dc, kb); err != nil {
		panic(err)
	}

	// Step 5: Show result
	fmt.Printf("Final Offer Applied: %+v\n", ruleCtx.Offer)

	// Some More
	dc, ruleCtx, err = exampleTwo()
	matchingRules, err = e.FetchMatchingRules(dc, kb)
	if err != nil {
		return
	}
	for _, rule := range matchingRules {
		fmt.Println("Matching Rule: ", rule)
	}
	if err := e.Execute(dc, kb); err != nil {
		panic(err)
	}

	// Step 5: Show result
	fmt.Printf("Final Offer Applied: %+v\n", ruleCtx.Offer)
}
