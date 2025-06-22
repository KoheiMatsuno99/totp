package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"totp/ctxtime"
)

func main() {
	secret, err := GenerateSecret()
	if err != nil {
		log.Fatal("シークレット生成エラー:", err)
	}

	totp := NewTOTP(secret)
	err = totp.GenerateQRCode("TestApp", "user@example.com", "qrcode.png")
	if err != nil {
		log.Printf("QRコード生成エラー: %v", err)
	}

	fmt.Println("\n=== TOTP コード生成 ===")
	for i := range 3 {
		ctx := context.Background()
		now := ctxtime.Now(ctx)
		code, err := totp.GenerateCode(&now)
		if err != nil {
			log.Printf("コード生成エラー: %v", err)
			continue
		}

		fmt.Printf("現在のTOTPコード: %s\n", code)

		now2 := ctxtime.Now(ctx)
		isValid := totp.Verify(code, &now2)
		fmt.Printf("検証結果: %t\n", isValid)

		if i < 2 {
			fmt.Println("30秒待機中...")
			time.Sleep(time.Duration(totp.Period) * time.Second)
		}
	}
}
