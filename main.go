package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	secret, err := GenerateSecret()
	if err != nil {
		log.Fatal("シークレット生成エラー:", err)
	}

	fmt.Println("生成されたシークレット:", secret)

	totp := NewTOTP(secret)

	fmt.Println("\n=== TOTP設定情報 ===")
	fmt.Println("Secret:", totp.Secret)
	fmt.Println("Digits:", totp.Digits)
	fmt.Println("Period:", totp.Period, "秒")

	err = totp.GenerateQRCode("TestApp", "user@example.com", "qrcode.png")
	if err != nil {
		log.Printf("QRコード生成エラー: %v", err)
	} else {
		fmt.Println("\nQRコードをqrcode.pngに保存しました")
	}

	fmt.Println("\n=== TOTP コード生成テスト ===")
	for i := range 3 {
		now := time.Now()
		code, err := totp.GenerateCode(&now)
		if err != nil {
			log.Printf("コード生成エラー: %v", err)
			continue
		}

		fmt.Printf("現在のTOTPコード: %s\n", code)

		now2 := time.Now()
		isValid := totp.Verify(code, &now2)
		fmt.Printf("検証結果: %t\n", isValid)

		if i < 2 {
			fmt.Println("30秒待機中...")
			time.Sleep(30 * time.Second)
		}
	}

	fmt.Println("\n=== 無効なコードのテスト ===")
	invalidCode := "123456"
	now3 := time.Now()
	isValid := totp.Verify(invalidCode, &now3)
	fmt.Printf("無効なコード '%s' の検証結果: %t\n", invalidCode, isValid)
}
