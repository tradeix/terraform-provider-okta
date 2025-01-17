package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppAutoLoginApplication_crud(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appAutoLogin)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appAutoLogin)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appAutoLogin, createDoesAppExist(okta.NewAutoLoginApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "sign_on_url", "https://example.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "sign_on_redirect_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "reveal_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "EDIT_USERNAME_AND_PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template", "user.firstName"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_push_status", "DONT_PUSH"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttr(resourceName, "admin_note", "admin note"),
					resource.TestCheckResourceAttr(resourceName, "enduser_note", "enduser note"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "sign_on_url", "https://exampleupdate.com/login.html"),
					resource.TestCheckResourceAttr(resourceName, "sign_on_redirect_url", "https://exampleupdate.com"),
					resource.TestCheckResourceAttr(resourceName, "reveal_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "shared_password", "sharedpassword"),
					resource.TestCheckResourceAttr(resourceName, "shared_username", "sharedusername"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "SHARED_USERNAME_AND_PASSWORD"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template", "user.firstName"),
					resource.TestCheckResourceAttr(resourceName, "user_name_template_push_status", "PUSH"),
					resource.TestCheckResourceAttrSet(resourceName, "logo_url"),
					resource.TestCheckResourceAttr(resourceName, "admin_note", "admin note updated"),
					resource.TestCheckResourceAttr(resourceName, "enduser_note", "enduser note updated"),
				),
			},
		},
	})
}

func TestAccAppAutoLoginApplication_timeouts(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appAutoLogin)
	resourceName := fmt.Sprintf("%s.test", appAutoLogin)
	config := `
resource "okta_app_auto_login" "test" {
  label       = "testAcc_replace_with_uuid"
  sign_on_url = "https://example.com/login.html"
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appAutoLogin, createDoesAppExist(okta.NewAutoLoginApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}
