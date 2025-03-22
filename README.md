# Ecommerce Offer Rules - GRule Engine powered DSL System

This module defines an example rule-based Domain Specific Language (DSL) powered by GRule for managing conditional e-commerce offers. It includes:

- **Protobuf-based schema** (`EcommerceOfferRule`) to define offer logic.
- **Go serializer** to convert the proto rule into [GRL](https://github.com/hyperjumptech/grule-rule-engine) compatible format.
- **JSON rule format** for easy authoring.

---

## 📦 Protobuf DSL Overview

**Proto File**: `proto/ecommerce_offer_rules.proto`

This defines the `EcommerceOfferRule` with:

- `Condition` blocks (based on customer profile, cart, behavior, etc.)
- `Action` blocks (what to apply: discount, coupon, etc.)

It uses custom field annotations (`grl_field_name`, `grl_operator`) to generate GRL rules.

---

## 🛠️ Compile Protobuf
### 1. Install Protocol Buffers Compiler

**macOS:**
```bash
brew install protobuf
```

**Ubuntu/Linux:**
```bash
sudo apt install -y protobuf-compiler
```

**Check version:**
```bash
protoc --version
```

### 2. Install Go Plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

(For gRPC)
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Ensure the plugins are in your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Ensure you have `protoc` and Go plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Then run:

```bash
protoc \
  --proto_path=proto \
  --go_out=proto \
  ecommerce_offer_rules.proto
```

## 📄 Sample JSON Rule

File: `sample_rule.json`

```json
{
  "name": "loyalty_discount_above_2000",
  "description": "Apply 10% discount if user is a loyalty member and cart total > 2000",
  "salience": 10,
  "conditions": [
    {
      "expressions": [
        {
          "input": "IS_LOYALTY_PROGRAM_MEMBER",
          "operator": "EQUALS",
          "value": { "boolVal": true }
        },
        {
          "input": "CART_TOTAL",
          "operator": "GREATER_THAN",
          "value": { "floatVal": 2000.0 }
        }
      ],
      "expressionJoinOperator": "AND"
    }
  ],
  "conditionJoinOperator": "AND",
  "actions": [
    {
      "output": "APPLY_DISCOUNT_PERCENT",
      "value": { "floatVal": 10.0 }
    },
    {
      "output": "PROMO_MESSAGE",
      "value": { "stringVal": "Congrats! You've unlocked a loyalty discount." }
    }
  ]
}
```

---

## ✅ Resulting GRL Output

```grl
rule loyalty_discount_above_2000 "Apply 10% discount if user is a loyalty member and cart total > 2000" salience 10 {
	when
		( Customer.IsLoyaltyProgramMember == true ) && ( Customer.CartTotal > 2000.00 )
	then
		Offer.ApplyDiscountPercent = 10.00;
		Offer.PromoMessage = "Congrats! You've unlocked a loyalty discount.";
}
```
