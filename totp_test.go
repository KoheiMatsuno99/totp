package main

import (
	"context"
	"testing"
	"time"
	"totp/ctxtime/ctxtimetest"
)

func TestGenerateSecret(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"基本的なシークレット生成"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret, err := GenerateSecret()
			if err != nil {
				t.Fatalf("シークレット生成でエラーが発生しました: %v", err)
			}

			if len(secret) == 0 {
				t.Fatal("シークレットが空です")
			}

			if len(secret) != 32 {
				t.Errorf("シークレットの長さが期待値と異なります。期待値: 32, 実際: %d", len(secret))
			}

			secret2, err := GenerateSecret()
			if err != nil {
				t.Fatalf("2つ目のシークレット生成でエラーが発生しました: %v", err)
			}

			if secret == secret2 {
				t.Error("連続して生成されたシークレットが同じです（ランダム性に問題がある可能性）")
			}
		})
	}
}

func TestTOTP_Verify(t *testing.T) {
	tests := []struct {
		name      string
		Secret    string
		code      string
		fixedTime time.Time
		want      bool
	}{
		{
			name:      "有効なコード",
			Secret:    "JBSWY3DPEHPK3PXP",
			code:      "742275",
			fixedTime: time.Unix(1234567890, 0),
			want:      true,
		},
		{
			name:      "無効なコード",
			Secret:    "JBSWY3DPEHPK3PXP",
			code:      "123456",
			fixedTime: time.Unix(1234567890, 0),
			want:      false,
		},
		{
			name:      "空のコード",
			Secret:    "JBSWY3DPEHPK3PXP",
			code:      "",
			fixedTime: time.Unix(1234567890, 0),
			want:      false,
		},
		{
			name:      "時間窓での検証（30秒前）",
			Secret:    "JBSWY3DPEHPK3PXP",
			code:      "742275",
			fixedTime: time.Unix(1234567890-30, 0),
			want:      true,
		},
		{
			name:      "時間窓での検証（30秒後）",
			Secret:    "JBSWY3DPEHPK3PXP",
			code:      "742275",
			fixedTime: time.Unix(1234567890+30, 0),
			want:      true,
		},
		{
			name:      "時間窓外（60秒前）",
			Secret:    "JBSWY3DPEHPK3PXP",
			code:      "742275",
			fixedTime: time.Unix(1234567890-60, 0),
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ctxtimetest.WithFixedNow(t, context.Background(), tt.fixedTime)
			totp := NewTOTP(tt.Secret)
			if got := totp.Verify(ctx, tt.code); got != tt.want {
				t.Errorf("TOTP.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
