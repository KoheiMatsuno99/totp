package main

import (
	"testing"
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
	type fields struct {
		Secret    string
		Digits    int
		Period    int64
		Algorithm string
	}
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "RFC6238テストベクター1",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args:    args{timestamp: 59},
			want:    "287082",
			wantErr: false,
		},
		{
			name: "RFC6238テストベクター2", 
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args:    args{timestamp: 1111111109},
			want:    "081804",
			wantErr: false,
		},
		{
			name: "無効なシークレット",
			fields: fields{
				Secret:    "INVALID@SECRET!",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args:    args{timestamp: 59},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TOTP{
				Secret:    tt.fields.Secret,
				Digits:    tt.fields.Digits,
				Period:    tt.fields.Period,
				Algorithm: tt.fields.Algorithm,
			}
			got, err := tr.GenerateCode(tt.args.timestamp)
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
	type fields struct {
		Secret    string
		Digits    int
		Period    int64
		Algorithm string
	}
	type args struct {
		code      string
		timestamp int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "有効なコード",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args: args{code: "287082", timestamp: 59},
			want: true,
		},
		{
			name: "無効なコード",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args: args{code: "123456", timestamp: 59},
			want: false,
		},
		{
			name: "空のコード",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args: args{code: "", timestamp: 59},
			want: false,
		},
		{
			name: "時間窓での検証（30秒前）",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args: args{code: "287082", timestamp: 59 - 30},
			want: true,
		},
		{
			name: "時間窓での検証（30秒後）",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args: args{code: "287082", timestamp: 59 + 30},
			want: true,
		},
		{
			name: "時間窓外（60秒前）",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args: args{code: "287082", timestamp: 59 - 60},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TOTP{
				Secret:    tt.fields.Secret,
				Digits:    tt.fields.Digits,
				Period:    tt.fields.Period,
				Algorithm: tt.fields.Algorithm,
			}
			if got := tr.Verify(tt.args.code, tt.args.timestamp); got != tt.want {
				t.Errorf("TOTP.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTOTP_GenerateQRCode(t *testing.T) {
	type fields struct {
		Secret    string
		Digits    int
		Period    int64
		Algorithm string
	}
	type args struct {
		issuer   string
		account  string
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "正常なQRコード生成",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args:    args{issuer: "TestApp", account: "user@example.com", filename: "test_qr.png"},
			wantErr: false,
		},
		{
			name: "無効なファイルパス",
			fields: fields{
				Secret:    "JBSWY3DPEHPK3PXP",
				Digits:    6,
				Period:    30,
				Algorithm: "SHA1",
			},
			args:    args{issuer: "TestApp", account: "user@example.com", filename: "/invalid/path/test_qr.png"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TOTP{
				Secret:    tt.fields.Secret,
				Digits:    tt.fields.Digits,
				Period:    tt.fields.Period,
				Algorithm: tt.fields.Algorithm,
			}
			if err := tr.GenerateQRCode(tt.args.issuer, tt.args.account, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("TOTP.GenerateQRCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
