package conf

func IsDebug() bool {
	return true
}

/**SessionSecret is the secret used to sign your cookies
 * TODO : Change the value before deployment.
 * secret is used to sign your cookies
 */
func SessionSecret() string {
	return "8M3W5A0CGYFx_bzjMzAFqLZ8esI3F0_CveBbgDZLd0hc2ManB3il2Cw9IPcY7Fr1"
}

/**SessionKeyPage is the random seed for key name for Page Id.
 */
func SessionKeyPage() string {
	return "PkNMgRN4kx_uxrmduaVK1AyL8L7aCxhVDHmSPWHpp9v6UD-BJGEMPMbRPQaa9Dc1"
}

/**SessionKeyUser is the random seed for key name for User Id.
 */
func SessionKeyUser() string {
	return "FXtlPiHcOtUHJ5Z7u_sw5CvdO23LAbvHhNeUmzqC59N5gmffoHki4-mIdbUb89Qw"
}

/**SessionKeySID is the key for session id
 */
func SessionKeySID() string {
	return "ogIbzGvEnFY2XfuQsLXv6a38dF49tvMuum4R27abuY5xAv0xF2Hc3SXkGoIcbWqd"
}

/**SessionValidHours defines how long a session could be idle for.
 */
func SessionValidHours() float64 {
	return 24 * 3
}

func GoogleSignInAppId() string {
	return ""
}
