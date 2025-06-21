package main

import (
	"testing"
	"time"
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

func TestTOTP_GenerateCode(t *testing.T) {
	type args struct {
		timestamp *time.Time
	}
	tests := []struct {
		name    string
		Secret  string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "固定時刻でのコード生成",
			Secret:  "JBSWY3DPEHPK3PXP",
			args:    args{timestamp: func() *time.Time { t := time.Unix(1234567890, 0); return &t }()},
			want:    "005924",
			wantErr: false,
		},
		{
			name:    "別の固定時刻でのコード生成",
			Secret:  "JBSWY3DPEHPK3PXP",
			args:    args{timestamp: func() *time.Time { t := time.Unix(1111111111, 0); return &t }()},
			want:    "050471",
			wantErr: false,
		},
		{
			name:    "無効なシークレット",
			Secret:  "INVALID@SECRET!",
			args:    args{timestamp: func() *time.Time { t := time.Unix(1234567890, 0); return &t }()},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totp := NewTOTP(tt.Secret)
			got, err := totp.GenerateCode(tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("TOTP.GenerateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TOTP.GenerateCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTOTP_Verify(t *testing.T) {
	type args struct {
		code      string
		timestamp *time.Time
	}
	tests := []struct {
		name   string
		Secret string
		args   args
		want   bool
	}{
		{
			name:   "有効なコード",
			Secret: "JBSWY3DPEHPK3PXP",
			args:   args{code: "005924", timestamp: func() *time.Time { t := time.Unix(1234567890, 0); return &t }()},
			want:   true,
		},
		{
			name:   "無効なコード",
			Secret: "JBSWY3DPEHPK3PXP",
			args:   args{code: "123456", timestamp: func() *time.Time { t := time.Unix(1234567890, 0); return &t }()},
			want:   false,
		},
		{
			name:   "空のコード",
			Secret: "JBSWY3DPEHPK3PXP",
			args:   args{code: "", timestamp: func() *time.Time { t := time.Unix(1234567890, 0); return &t }()},
			want:   false,
		},
		{
			name:   "時間窓での検証（30秒前）",
			Secret: "JBSWY3DPEHPK3PXP",
			args:   args{code: "005924", timestamp: func() *time.Time { t := time.Unix(1234567890-30, 0); return &t }()},
			want:   true,
		},
		{
			name:   "時間窓での検証（30秒後）",
			Secret: "JBSWY3DPEHPK3PXP",
			args:   args{code: "005924", timestamp: func() *time.Time { t := time.Unix(1234567890+30, 0); return &t }()},
			want:   true,
		},
		{
			name:   "時間窓外（60秒前）",
			Secret: "JBSWY3DPEHPK3PXP",
			args:   args{code: "005924", timestamp: func() *time.Time { t := time.Unix(1234567890-60, 0); return &t }()},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totp := NewTOTP(tt.Secret)
			if got := totp.Verify(tt.args.code, tt.args.timestamp); got != tt.want {
				t.Errorf("TOTP.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTOTP_GenerateQRCode(t *testing.T) {
	type args struct {
		issuer   string
		account  string
		filename string
	}
	tests := []struct {
		name    string
		Secret  string
		args    args
		wantErr bool
	}{
		{
			name:    "正常なQRコード生成",
			Secret:  "JBSWY3DPEHPK3PXP",
			args:    args{issuer: "TestApp", account: "user@example.com", filename: "test_qr.png"},
			wantErr: false,
		},
		{
			name:    "無効なファイルパス",
			Secret:  "JBSWY3DPEHPK3PXP",
			args:    args{issuer: "TestApp", account: "user@example.com", filename: "/invalid/path/test_qr.png"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totp := NewTOTP(tt.Secret)
			if err := totp.GenerateQRCode(tt.args.issuer, tt.args.account, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("TOTP.GenerateQRCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
