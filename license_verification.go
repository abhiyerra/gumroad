// https://help.gumroad.com/11165-Digging-Deeper/license-keys
package gumroad

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type LicenseResponse struct {
	Success  bool `json:"success"`
	Uses     int  `json:"uses"`
	Purchase struct {
		ID          string `json:"id"`
		ProductName string `json:"product_name"`
		CreatedAt   string `json:"created_at"`
		FullName    string `json:"full_name"`
		Variants    string `json:"variants"`
		// purchase was refunded, non-subscription product only
		Refunded bool `json:"refunded"`
		// purchase was refunded, non-subscription product only
		Chargebacked bool `json:"chargebacked"`

		// subscription product only
		// subscription was cancelled,
		SubscriptionCancelledAt *string `json:"subscription_cancelled_at"`
		// we were unable to charge the subscriber's card
		SubscriptionFailedAt *string `json:"subscription_failea_dat"`

		CustomFields []string `json:"custom_fields"`
		Email        string   `json:"email"`
	}
}

// curl
//   \ -d "product_permalink=QMGY"
//   \ -d "license_key=YOUR_CUSTOMERS_LICENSE_KEY"
//   \ -X POST

// productPermalink (the unique permalink of the product)
// licensekey (the license key provided by your customer)
// incrementUsesCount ("true"/"false", optional, default: "true")
func VerifyLicense(productPermalink, licenseKey string, incrementUsesCount bool) error {
	var (
		response LicenseResponse
	)

	resp, err := http.PostForm("https://api.gumroad.com/v2/licenses/verify",
		url.Values{
			"product_permalink":    {productPermalink},
			"license_key":          {licenseKey},
			"increment_uses_count": {strconv.FormatBool(incrementUsesCount)},
		})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Purchase.Refunded {
		return errors.New("gumroad: product was refunded")
	}

	if response.Purchase.Chargebacked {
		return errors.New("gumroad: product was chargebacked")
	}

	if response.Purchase.SubscriptionCancelledAt != nil {
		return errors.New("gumroad: subscription cancelled")
	}

	if response.Purchase.SubscriptionFailedAt != nil {
		return errors.New("gumroad: subscription failed")
	}

	return nil
}
