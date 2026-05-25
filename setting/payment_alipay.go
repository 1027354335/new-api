package setting

import "github.com/QuantumNous/new-api/common"

var AlipayEnabled = false
var AlipayAppId = ""
var AlipayPrivateKey = ""
var AlipayPublicKey = ""
var AlipaySandbox = false
var AlipayCallbackUrl = ""
var AlipayNotifyUrl = ""
var AlipayUsdToCnyRate = 7.2
var AlipayBridgeEnabled = common.GetEnvOrDefaultBool("ALIPAY_BRIDGE_ENABLED", false)
var AlipayBridgeCreateUrl = common.GetEnvOrDefaultString("ALIPAY_BRIDGE_CREATE_URL", "")
var AlipayBridgeSecret = common.GetEnvOrDefaultString("ALIPAY_BRIDGE_SECRET", "")
