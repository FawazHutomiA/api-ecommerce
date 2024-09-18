package helper

const (
	REGISTER_USER_EMAIL_HTML = `
	<!DOCTYPE html>
<html>
<head>
    <style>
      @import url('https://fonts.googleapis.com/css2?family=Montserrat&display=swap');
    </style>
  <style>
    
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f4;
      
    }
    .call-center {
      color: #8e44ad;
      text-decoration: underline;
    }
    .container {
      width: 80%;
      margin: 0 auto;
      padding: 20px;
      background-color: #fff;
      border-radius: 5px;
      box-shadow: 0px 0px 5px rgba(0,0,0,0.1);
    }
    .header {
      text-align: center;
      margin-bottom: 30px;
    }
    .content {
      font-size: 16px;
      font-weight: 500;
      line-height: 26px;
      letter-spacing: 0.02em;
      text-align: left;
    }
    .footer {
      text-align: center;
      margin-top: 30px;
      font-size: 14px;
      color: #888;
    }
    .button-link {
      display: inline-block;
      background-color: #8e44ad;
      color: #fff;
      padding: 10px 20px;
      border-radius: 5px;
      text-decoration: none;
    }
    .button-link:hover {
      background-color: #7a378b;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <img src="{example_main_logo}" alt="Example logo" style="height: 60px; margin: 0 auto; display: block" />
      <h1>Permintaan Verifikasi Akun</h1>
    </div>
    <div class="content" style="color: inherit; text-decoration: none;">
      <p>Silakan tekan tombol dibawah ini untuk verifikasi akun anda:</p>
      <p><a href="{$1}" style="color: white; text-decoration: none;" class="button-link"><strong>Verifikasi Akun</strong></a></p>
      <p>Jika Anda membutuhkan bantuan, silakan hubungi <a href="mailto:noreply@example.id" class="call-center">Call Center</a> kami.</p>
      <p>Salam hangat,</p>
      <p>Tim Example</p>
      <img src="{example_main_logo}" alt="Example logo" style="height: 60px; margin: 0 auto; display: block" />
    </div>
    <div class="footer">
      <p>example.id</p>
      <p>JI. Sidomukti No.21, Sukaluyu, Kec. Cibeunying Kaler, Kota Bandung, Jawa Barat 40123</p>
    </div>
  </div>
</body>
</html>
	`
)
