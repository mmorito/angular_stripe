package main

import (
	"fmt"
	// "net/http"
	"os"
	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/coupon"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/sub"
	"github.com/stripe/stripe-go/transfer"
)

var PUBLIC_KEY = "sk_test_tEcSylz1n1P7KihpdShXNBac"

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("引数足りん！！\n")
		return
	}
	token := os.Args[1]
	stripe.Key = PUBLIC_KEY
	fmt.Printf(token + "\n")
	// シンプルなお支払い
	simplePayment(token)

	// カード情報を登録してお支払い
	if cur, err := registerCustomerInfo(token); err == nil {
		paymentByCustomerInfo(token, cur)
	}

	// クーポン作成
	createCoupon()

	// クーポンを使ってお支払い
	if cur, err := registerCustomerInfo(token); err == nil {
		paymentByCustomerInfoWithCoupon(cur)
	}

	// プランを作成
	createPlan()

	// プランを利用して定期購読
	createSubscribe()

	// アカウント追加
	createCustomAccount()

	// アカウントを認証
	acceptAccount()

	// 売上を送金
	transferSales(token)
}

// カード情報を顧客として保存
func registerCustomerInfo(token string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String("morito@mnes.org"),
	}
	params.SetSource(token)
	cus, err := customer.New(params)

	return cus, err
}

// 保存した顧客IDをもとに支払いを実行
func paymentByCustomerInfo(cur *stripe.Customer) {
	params := &stripe.ChargeParams{
		Amount:       stripe.Int64(1500),
		Currency:     stripe.String(string(stripe.CurrencyJPY)),
		Description:  stripe.String("これはテストです"),
		Customer:     &cur.ID,
		ReceiptEmail: stripe.String("morito@mnes.org"),
	}
	_, _ = charge.New(params)
}

// 普通のお支払い
func simplePayment(token string) {
	params := &stripe.ChargeParams{
		Amount:              stripe.Int64(1000),
		Currency:            stripe.String(string(stripe.CurrencyJPY)),
		Description:         stripe.String("これはテストです"),
		StatementDescriptor: stripe.String("statement"),
		ReceiptEmail:        stripe.String("morito@mnes.org"),
	}
	params.SetSource(token)
	_, _ = charge.New(params)
}

// クーポンを登録する
func createCoupon() {
	params := &stripe.CouponParams{
		PercentOff: stripe.Float64(25),
		Duration:   stripe.String(string(stripe.CouponDurationOnce)),
		ID:         stripe.String("25OFF"),
	}
	_, _ = coupon.New(params)
}

// クーポンを使ってお支払い
func paymentByCustomerInfoWithCoupon(cur *stripe.Customer) {
	var amount float64 = 1000
	if c, err := coupon.Get("25OFF", nil); err == nil {
		off := c.PercentOff / 100
		discount := amount * off
		amount = amount - discount

		params := &stripe.ChargeParams{
			Amount:       stripe.Int64(int64(amount)),
			Currency:     stripe.String(string(stripe.CurrencyJPY)),
			Description:  stripe.String("クーポンを使いました"),
			Customer:     &cur.ID,
			ReceiptEmail: stripe.String("morito@mnes.org"),
		}
		_, _ = charge.New(params)
	}
}

// プランを追加
func createPlan() {
	params := &stripe.PlanParams{
		Amount:   stripe.Int64(30000),
		Interval: stripe.String("month"),
		Product: &stripe.PlanProductParams{
			Name: stripe.String("Standard Plan"),
		},
		Currency: stripe.String(string(stripe.CurrencyJPY)),
	}
	_, _ = plan.New(params)
}

// プランを利用して定期購読
func createSubscribe() {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String("cus_Dqq5YFZLdPTtF5"),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String("plan_DqtzR4mPLkEOgq"),
			},
		},
	}
	_, _ = sub.New(params)
}

// アカウントの作成
func createCustomAccount() {
	params := &stripe.AccountParams{
		Country: stripe.String("JP"),
		Type:    stripe.String(string(stripe.AccountTypeCustom)),
		Email:   stripe.String("morito@gmail.com"),
	}
	acct, _ := account.New(params)
	fmt.Println(acct.Keys.Secret)
}

// アカウントの利用規約の同意
func acceptAccount() {
	params := &stripe.AccountParams{
		TOSAcceptance: &stripe.TOSAcceptanceParams{
			Date: stripe.Int64(time.Now().Unix()),
			IP:   stripe.String("192.168.x.x"), // Assumes you're not using a proxy
		},
	}
	_, _ = account.Update("acct_1DPAhhDSowBmK9ov", params)
}

// お支払いと同時に送金
func transferSales(token string) {
	// お支払い
	params := &stripe.ChargeParams{
		Amount:        stripe.Int64(3000),
		Currency:      stripe.String(string(stripe.CurrencyJPY)),
		TransferGroup: stripe.String("transfer002"),
	}
	params.SetSource(token)
	_, _ = charge.New(params)

	// 1人目に送金
	transferParams := &stripe.TransferParams{
		Amount:        stripe.Int64(1000),
		Currency:      stripe.String(string(stripe.CurrencyJPY)),
		Destination:   stripe.String("acct_1DPAhhDSowBmK9ov"),
		TransferGroup: stripe.String("transfer002"),
	}
	_, _ = transfer.New(transferParams)

	// 2人目に送金
	secondTransferParams := &stripe.TransferParams{
		Amount:        stripe.Int64(1000),
		Currency:      stripe.String(string(stripe.CurrencyJPY)),
		Destination:   stripe.String("acct_1DPAsqARF9i7A4af"),
		TransferGroup: stripe.String("transfer002"),
	}
	_, _ = transfer.New(secondTransferParams)
}
