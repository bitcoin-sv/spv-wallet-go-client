package examples

const (
	// A string representing the alias part of the Paymail address.
	Alias = "test"
	// Domain specifies the domain portion of the Paymail address.
	Domain = "example.com"
	// Paymail constructs the complete Paymail address by concatenating
	// Alias and Domain with an "@" separator.
	Paymail = Alias + "@" + Domain
)

const (
	// AdminXPriv is used to authenticate as an admin in the spv-wallet.
	// NOTE: The provided key is a default key that matches SPV Wallet's default configuration.
	AdminXPriv string = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"

	// UserXPriv is used to authenticate as a user in the spv-wallet.
	// You can generate a new key pair using "generate_keys" example.
	UserXPriv string = ""

	// UserXPub is a public part of the UserXPriv key. It is used to create a new user in the spv-wallet by the admin.
	UserXPub string = ""
)
