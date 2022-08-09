package staticemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticEmail(t *testing.T) {
	var testData = []struct {
		name    string
		key     string
		contain string
	}{
		{
			name:    "Test EMAIL_MERCHANT_UPGRADE_TEMPLATE_ID ",
			key:     "EMAIL_MERCHANT_UPGRADE_TEMPLATE_ID",
			contain: "Terima kasih telah mengajukan upgrade",
		},
		{
			name:    "Test EMAIL_MERCHANT_REGISTRATION_TEMPLATE_ID",
			key:     "EMAIL_MERCHANT_REGISTRATION_TEMPLATE_ID",
			contain: "Terima kasih telah mendaftar ke Marketplace",
		},
		{
			name:    "Test EMAIL_MERCHANT_UPGRADE_APPROVAL",
			key:     "EMAIL_MERCHANT_UPGRADE_APPROVAL",
			contain: "Pengajuan upgrade toko telah disetujui",
		},
		{
			name:    "Test EMAIL_SUCCESS_FORGOT_PASSWORD_TEMPLATE_ID",
			key:     "EMAIL_SUCCESS_FORGOT_PASSWORD_TEMPLATE_ID",
			contain: "Kata sandi baru untuk akun kamu",
		},
		{
			name:    "Test EMAIL_FORGOT_PASSWORD_TEMPLATE_ID",
			key:     "EMAIL_FORGOT_PASSWORD_TEMPLATE_ID",
			contain: "Silahkan klik tombol besar di bawah untuk mengganti kata sandi Anda",
		},
		{
			name:    "Test EMAIL_SUCCESS_REGISTRATION_TEMPLATE_ID",
			key:     "EMAIL_SUCCESS_REGISTRATION_TEMPLATE_ID",
			contain: "Terima kasih kami ucapkan, atas pendaftaran yang telah dilakukan",
		},
		{
			name:    "Test EMAIL_PERSONAL_REGISTRATION_TEMPLATE_ID",
			key:     "EMAIL_PERSONAL_REGISTRATION_TEMPLATE_ID",
			contain: "kami butuh mengkonfirmasikan alamat email Anda",
		},
		{
			name:    "Test EMAIL_ADD_MEMBER_TEMPLATE_ID",
			key:     "EMAIL_ADD_MEMBER_TEMPLATE_ID",
			contain: "silakan klik tombol di bawah ini untuk pembuatan password login",
		},
		{
			name:    "Test EMAIL_MERCHANT_ACTIVATION",
			key:     "EMAIL_MERCHANT_ACTIVATION",
			contain: "Selamat! ##MERCHANTNAME## telah aktif di Marketplace",
		},
		{
			name:    "Test EMAIL_MERCHANT_UPGRADE_REJECT",
			key:     "EMAIL_MERCHANT_UPGRADE_REJECT",
			contain: "Terima kasih telah melakukan <em>upgrade</em>",
		},
		{
			name:    "Test EMAIL_MERCHANT_REJECT",
			key:     "EMAIL_MERCHANT_REJECT",
			contain: "Terima kasih telah mendaftar menjadi Merchant di Marketplace",
		},
		{
			name:    "Test random",
			key:     "random",
			contain: "",
		},
	}

	for _, tc := range testData {
		emailContent := GetFallbackEmailContent(tc.key)
		assert.Contains(t, emailContent, tc.contain)
	}
}
