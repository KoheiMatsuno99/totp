package main

const LoginTemplate = `<!DOCTYPE html>
<html>
<head>
	<title>TOTP Login</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 400px; margin: 100px auto; padding: 20px; }
		form { background: #f5f5f5; padding: 30px; border-radius: 8px; }
		input { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; border-radius: 4px; }
		button { width: 100%; padding: 12px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
		button:hover { background: #0056b3; }
		.error { color: red; margin: 10px 0; }
	</style>
</head>
<body>
	<h2>TOTP ログイン</h2>
	<form method="POST" action="/verify">
		<input type="email" name="email" placeholder="メールアドレス" required>
		<input type="text" name="code" placeholder="6桁のTOTPコード" pattern="[0-9]{6}" required>
		<button type="submit">ログイン</button>
	</form>
	{{if .Error}}
		<div class="error">{{.Error}}</div>
	{{end}}
</body>
</html>`

const SetupTemplate = `<!DOCTYPE html>
<html>
<head>
	<title>TOTP セットアップ</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; text-align: center; }
		.setup { background: #e3f2fd; color: #0d47a1; padding: 30px; border-radius: 8px; border: 1px solid #bbdefb; }
		.qr-section { margin: 20px 0; padding: 20px; background: white; border-radius: 5px; }
		.secret { font-family: monospace; background: #f5f5f5; padding: 10px; border-radius: 4px; word-break: break-all; }
		button { padding: 10px 20px; background: #2196f3; color: white; border: none; border-radius: 4px; cursor: pointer; margin: 10px; }
		button:hover { background: #1976d2; }
		.instructions { text-align: left; background: #fff3e0; padding: 15px; border-radius: 5px; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="setup">
		<h2>🔐 TOTP セットアップ</h2>
		<p><strong>ユーザー:</strong> {{.Email}}</p>
		<div class="instructions">
			<h3>セットアップ手順:</h3>
			<ol>
				<li>Google Authenticator、Authy等のTOTPアプリをインストール</li>
				<li>下記のQRコードをスキャン、または手動でシークレットキーを入力</li>
				<li>アプリで生成された6桁のコードでログインテスト</li>
			</ol>
		</div>
		<div class="qr-section">
			<h3>QRコード:</h3>
			<div style="font-family: monospace; font-size: 12px; word-break: break-all;">{{.QRUrl}}</div>
			<br>
			<h3>シークレットキー:</h3>
			<div class="secret">{{.Secret}}</div>
		</div>
		<button onclick="location.href='/'">ログイン画面に戻る</button>
	</div>
</body>
</html>`

const SuccessTemplate = `<!DOCTYPE html>
<html>
<head>
	<title>ログイン成功</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 500px; margin: 100px auto; padding: 20px; text-align: center; }
		.success { background: #d4edda; color: #155724; padding: 20px; border-radius: 8px; border: 1px solid #c3e6cb; }
		.user-info { margin: 20px 0; padding: 15px; background: #f8f9fa; border-radius: 5px; }
		button { padding: 10px 20px; background: #28a745; color: white; border: none; border-radius: 4px; cursor: pointer; }
		button:hover { background: #218838; }
	</style>
</head>
<body>
	<div class="success">
		<h2>🎉 ログイン成功！</h2>
		<div class="user-info">
			<strong>ユーザー:</strong> {{.Email}}
		</div>
		<p>TOTP認証が正常に完了しました。</p>
		<button onclick="location.href='/'">再度ログイン</button>
	</div>
</body>
</html>`