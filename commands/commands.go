package commands

import (
	"context"
	"fmt"
	"github.com/platformsh/platformify/vendorization"
	"github.com/spf13/viper"
)

// Execute executes the ify command and sets flags appropriately.
func Execute(assets *vendorization.VendorAssets) error {
	cmd := NewPlatformifyCmd(assets)
	cmd.PersistentFlags().BoolP(
		"no-interaction",
		"",
		false,
		fmt.Sprintf(
			"Do not ask any interactive questions; accept default values. Equivalent to using the environment variable: %s_CLI_NO_INTERACTION=1",
			assets.NIPrefix,
		),
	)
	viper.BindPFlag("no-interaction", cmd.PersistentFlags().Lookup("no-interaction"))
	validateCmd := NewValidateCommand(assets)
	cmd.AddCommand(validateCmd)
	return cmd.ExecuteContext(vendorization.WithVendorAssets(context.Background(), assets))
}
