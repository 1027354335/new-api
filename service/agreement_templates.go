package service

// Global variable containing standard templates for 6 languages
var agreementTemplates = map[string]AgreementTemplate{
	"zh": {
		Title: "AI Token 购买及使用协议",
		MetaKeys: map[string]string{
			"version":     "版本号：V1.1",
			"effective":   "生效日期：{{date}}",
			"provider":    "服务提供方（平台）：OSS Energietechnik GmbH",
			"address":     "注册地址：Adam-Opel-Straβe 16-18, 60386 Frankfurt am Main",
			"email":       "联系邮箱：info@oss-energietechnik.de",
			"user":        "用户：{{user}}",
			"order":       "支付订单：{{order}} ({{amount}})",
			"signMethod":  "签署方式：用户勾选同意本协议并点击 “同意协议并支付”“确认购买” 等同类按钮，即视为以电子方式签署本协议，协议即时生效。",
			"applicable":  "适用区域： 境内（中国大陆）、境外用户均可适用； 境内用户优先适用中华人民共和国法律；境外用户同时遵守其所在地强制性消费者保护、数据保护、税务、电子商务、出口管制与数字服务相关法律。",
			"declaration": "重要声明：本协议为商业合同文本，非法律意见；正式上线前，建议结合用户所在地法域、支付渠道、税务及数据合规要求完成最终校验。",
		},
		Sections: []string{
			"第一条 定义\n" +
				"AI Token：指用户向平台购买的、用于调用平台人工智能模型、算法能力、API 服务、智能体服务、文本生成、图像生成、语音处理、数据分析或其他 AI 功能的服务调用额度 / 消耗凭证，不属于货币、虚拟货币、电子货币、证券、金融产品、储值卡或可转让资产。\n" +
				"服务：指平台向用户提供的 AI 模型调用、API 接入、在线生成、数据处理、账户管理、用量统计、技术支持等相关服务。\n" +
				"账户：指用户在平台注册、登录或由平台分配的账户、API Key、密钥、组织 ID、项目 ID 或其他身份标识。\n" +
				"输入内容：指用户在使用服务时提交、上传、输入、传输或调用的文本、图片、音频、视频、代码、文件、数据、提示词、接口参数等内容。\n" +
				"输出内容：指服务根据用户输入内容生成或返回的文本、图片、音频、视频、代码、分析结果、模型回复或其他结果。\n" +
				"消耗：指用户因调用 AI 服务、发起请求、生成内容、处理数据或使用相关功能而扣减 AI Token 的行为。\n" +
				"合规定性：本 AI Token 仅为平台专属 AI 服务调用额度，无流通性、法偿性，不得兑换法定货币、虚拟货币或用于交易炒作；境内适用央行等监管规定，境外同时符合用户所在地反虚拟货币与数字资产监管规则。",

			"第二条 协议签署与电子确认\n" +
				"本协议以电子数据形式展示、确认、签署、留存，与书面协议具有同等法律效力。\n" +
				"用户勾选同意本协议并点击支付、确认购买、提交订单、完成付款或实际使用 AI Token 的行为，即构成对本协议的有效确认和签署。\n" +
				"平台有权记录并保存用户确认本协议时的账户信息、订单编号、支付记录、IP 地址、设备信息、浏览器信息、操作时间、协议版本号、勾选记录、点击记录、日志记录等，用于证明协议订立、履行、争议解决和合规审计。\n" +
				"用户不得以未签署纸质合同、未加盖实体印章、未进行线下签字为由否认本协议的效力。\n" +
				"如用户为企业、机构或其他组织，代表该主体进行购买或使用服务的人员确认其已获得充分授权；未经授权而代表他人购买或使用服务的，由实际操作人承担相应责任。\n" +
				"本协议中退款、责任限制等免除 / 限制平台责任的条款，已通过正文加粗、支付页面独立弹窗完成显著提示，用户确认已充分理解并自愿接受。",

			"第三条 Token 购买、计价与交付\n" +
				"用户可根据平台页面展示的套餐、价格、币种、数量、有效期、适用模型、消耗规则和其他说明购买 AI Token。\n" +
				"AI Token 的价格、兑换比例、消耗规则、支持模型、上下文长度、并发限制、功能范围等，以用户下单时平台展示的信息为准；平台公示调整的，按公示内容执行。\n" +
				"除非平台另有明确说明，AI Token 仅限平台内专属服务使用，不属于金融资产，不可转让、交易、兑换其他资产。\n" +
				"平台在收到用户支付款项并确认订单后，将在 24 小时内为用户账户充值对应 AI Token；逾期未到账的，用户有权无条件申请全额退款。\n" +
				"因支付渠道、银行、第三方支付机构、外汇结算、反欺诈审查、税务审查或合规审查导致的到账延迟、手续费扣除等，平台不承担因第三方原因造成的责任，但应在合理范围内协助用户查询。\n" +
				"用户应确保订单信息、账户信息、发票信息、付款主体信息真实、准确、完整；因用户信息错误导致充值失败、发票错误、账户归属争议或其他损失的，由用户自行承担责任。",

			"第四条 Token 使用规则\n" +
				"用户购买的 AI Token 仅可在平台指定范围内使用，具体可用于哪些模型、服务、API、功能或产品，以平台页面、控制台或订单说明为准。\n" +
				"AI Token 按照平台公布的计量方式消耗，包括但不限于输入字符数、输出字符数、token 数、图片张数、音频时长、视频时长、请求次数、模型类型、计算资源、存储资源、插件调用、工具调用等。\n" +
				"不同模型、不同功能、不同地区节点、不同服务等级的 Token 消耗标准可能不同，用户应在使用前查看相关计费说明。\n" +
				"AI Token 一经消耗，即视为服务已经交付或部分交付；除本协议另有约定或法律强制规定外，已消耗的 AI Token 不支持恢复、退还、转让或折现。\n" +
				"用户应妥善保管账户、密码、API Key、访问密钥、验证码及其他身份凭证；通过用户账户或密钥发起的调用，均视为用户本人或其授权人员的行为。\n" +
				"如用户发现账户或密钥泄露、被盗用或存在异常调用，应立即通知平台并采取重置密码、禁用密钥、关闭接口等措施；平台在收到通知前已发生的 Token 消耗，原则上由用户承担，但平台存在故意或重大过失的除外。\n" +
				"用户不得转让、出售、出租、出借、提现、折现或在平台外交易 AI Token，不得将其用于平台外交易或融资活动。",

			"第五条 有效期、过期与续费\n" +
				"AI Token 的有效期以平台在购买页面、订单页面、套餐说明或用户账户后台中明确展示的信息为准；未明确展示的，默认有效期为自充值到账之日起 12 个月。\n" +
				"有效期届满后，未使用的 AI Token 将自动失效，平台不再提供使用、退还、兑换或延期服务，法律强制规定另有要求的除外。\n" +
				"平台可根据运营安排提供延期、续费、升级、套餐转换等服务，具体规则以平台届时展示为准。\n" +
				"用户购买新套餐后，新旧 Token 的使用顺序以平台系统规则为准；如无特别说明，平台可优先消耗即将过期的 Token。",

			"第六条 退款与撤销（境内外通用 + 法域适配）\n" +
				"核心规则：AI Token 为即时交付、即时消耗的数字服务，不适用七日无理由退货；支付到账后原则上不予无理由退款，本条 2、3、4 款及法律强制规定除外。\n" +
				"境内用户退款规则：（1）充值到账后，不适用退款规则，不予退款；（2）已消耗的 AI Token、赠送 Token、活动 Token、试用 Token、优惠券抵扣部分、因违反本协议被限制使用的 Token、用户自身操作失误（如误充、误调用）导致消耗的 Token，不予退款、提现、转让或折现。\n" +
				"境外用户退款规则（含欧盟、英国等法域）：（1）用户所在地有法定撤回权的，用户勾选确认“立即履行、充值即丧失撤回权”后生效；（2）未确认立即履行且未消耗 Token 的，按用户所在地法定撤回期执行；已消耗部分不予退款。\n" +
				"平台责任退款：因平台原因导致用户连续 30 个自然日无法使用已购买且未过期 Token 对应的核心服务，平台应在收到用户有效通知后 15 个工作日内，提供补偿等值 Token、延长有效期或按未使用部分退款的补救措施；逾期未补救的，用户有权要求按未使用 Token 比例退款。",

			"第七条 用户合规义务\n" +
				"用户承诺遵守中华人民共和国法律及用户所在地法律法规、监管要求、行业规范、出口管制、经济制裁、数据跨境传输规则及公共秩序、善良风俗。\n" +
				"用户不得生成、传播或协助生成违法、侵权、欺诈、虚假、仇恨、骚扰、暴力、色情、恐怖主义、极端主义、自残、自杀、毒品、武器、恶意代码、网络攻击、非法金融活动等违规内容。\n" +
				"用户不得侵犯他人知识产权、商业秘密、隐私权、个人信息权益、肖像权、名誉权或其他合法权益。\n" +
				"用户不得未经授权抓取、复制、训练、反向工程、破解、绕过平台技术限制、安全机制、访问控制或计费系统。\n" +
				"用户不得将服务用于自动化垃圾信息、刷量、虚假评论、钓鱼、诈骗、冒充他人、绕过内容审核或批量生成违法违规内容。\n" +
				"用户不得向平台提交其无权处理的数据、个人信息、敏感个人信息、保密信息、国家秘密、商业秘密或受出口管制限制的数据。\n" +
				"平台有权根据法律法规、监管要求、平台规则、风控策略或第三方模型供应商要求，对用户输入内容、输出内容、调用行为、账户行为进行必要的安全审查、风控识别和合规处理。",

			"第八条 AI 输出内容与风险提示\n" +
				"用户理解并同意，AI 输出内容由算法模型根据输入内容、模型参数、训练数据、上下文和系统配置生成，可能存在不准确、不完整、过时、虚构、偏差、不适用或不可解释的情况。\n" +
				"输出内容不构成法律、财务、医疗、投资、工程、安全、合规或其他专业意见；用户如需依赖输出内容作出重要决策，应自行进行人工复核，并咨询具备资质的专业人士。\n" +
				"用户应对其输入内容、使用方式、输出内容的审查、发布、传播、应用及后果承担全部责任。\n" +
				"平台不保证 AI 输出内容完全准确、服务连续可用、满足用户特定目的，亦不保证输出内容不会与其他用户生成内容相同或相似。\n" +
				"对于用户将输出内容用于商业用途、公开发布、广告宣传、投标文件、法律文件、医疗建议、投资建议、工程设计、自动化决策或其他高风险场景的，用户应自行承担合规审查和使用风险。",

			"第九条 知识产权\n" +
				"平台及其关联方、授权方对平台系统、软件、模型、算法、接口、文档、页面设计、商标、标识、数据库、技术方案、商业模式及相关知识产权享有完整权利。\n" +
				"用户购买 AI Token 并不取得平台模型、算法、软件、源代码、系统架构、商标或其他知识产权的所有权。\n" +
				"用户对其合法提供的输入内容保留其依法享有的权利。\n" +
				"在用户遵守本协议且已支付相应费用的前提下，平台在其权利范围内允许用户使用服务生成的输出内容；但输出内容是否可独占、是否侵犯第三方权利、是否可登记为知识产权、是否可商用，应由用户结合具体内容和适用法律自行判断，侵权风险由用户自行承担。\n" +
				"用户授权平台在提供服务、计费、风控、安全审查、故障排查、服务优化、合规审计和履行法律义务所必要的范围内，处理输入内容、输出内容和使用数据。\n" +
				"如平台使用第三方基础模型、云服务、支付服务或其他第三方服务，相关权利限制、内容政策、使用限制可能同时适用于用户。",

			"第十条 数据保护与隐私（境内外双合规）\n" +
				"境内用户：适用《中华人民共和国个人信息保护法》，遵循最小必要原则；平台仅收集提供服务所必需的信息：①注册信息（手机号、邮箱）；②支付信息（支付宝订单号、支付金额）；③设备信息（IP 地址、浏览器类型）；④用量信息（Token 充值 / 消耗记录），不收集身份证号、银行卡号等敏感信息。\n" +
				"境外用户：适用其所在地数据保护法律（如 GDPR）；平台按用户所在地法律要求处理个人信息，不跨境传输数据至未达数据保护标准的地区，必要时可与用户签署数据处理协议。\n" +
				"用户个人信息留存期限：自账户注销或服务终止之日起 1 年，到期后自动匿名化或删除；订单记录、税务记录、Token 充值 / 消耗记录留存 3 年（适配税务与合规审计要求）。\n" +
				"如用户向平台输入涉及他人的个人信息、敏感个人信息、商业秘密、保密信息或受监管数据，用户应确保其具有合法处理基础，并已履行必要的告知、同意、授权、评估或备案义务。\n" +
				"平台可在必要范围内委托合规第三方服务商处理数据，包括云服务商、模型服务商、支付机构、短信服务商、发票服务商、客户支持系统、数据安全服务商等，并要求相关第三方采取合理的数据保护措施。",

			"第十一条 税费、发票与跨境支付\n" +
				"用户应按照订单页面显示的金额支付费用。相关费用是否含税，以平台页面、订单或发票说明为准。\n" +
				"因用户所在地、付款方式、支付渠道、收单机构、银行、外汇结算、增值税、消费税、销售税、预扣税、数字服务税或其他税费产生的额外费用，由双方根据适用法律和订单约定承担。涉及外币支付的，实际扣款金额可能因汇率、银行手续费、支付渠道费用或结算时间不同而发生差异，该等差异由用户自行承担，平台另有承诺的除外。\n" +
				"用户需要发票、收据或税务凭证的，应按平台要求提供真实、准确、完整的开票信息；因用户提供信息错误、付款主体与开票主体不一致导致无法开票或需更正发票的，用户应配合处理并承担相应损失。",

			"第十二条 服务变更、中断与维护\n" +
				"平台可根据业务发展、技术升级、模型供应、成本变化、法律法规、监管要求或第三方服务变化，对服务内容、模型类型、计费规则、Token 消耗比例、使用限制、功能范围进行调整；但如前述变更导致用户已购买且未使用的 Token 可使用范围显著缩减、单位服务成本显著上升（Token 消耗比例提高超过 20%），用户有权在变更生效前申请按比例退还未使用 Token 对应费用；前述变更不影响已购买未消耗 Token 在原服务范围内的使用权益。\n" +
				"平台进行重大变更时，应通过网站公告、站内通知、邮件、短信或其他合理方式通知用户。\n" +
				"因系统维护、升级、网络故障、云服务故障、第三方模型服务故障、支付渠道故障、安全事件、不可抗力或监管要求导致服务中断、延迟、不可用的，平台应尽合理努力恢复服务，不承担间接损失。\n" +
				"平台不对因用户网络、设备、接口配置、密钥管理、代码错误、调用参数错误、违反使用规则或第三方原因导致的使用失败承担责任。",

			"第十三条 账户限制与违约处理\n" +
				"如用户违反本协议或适用法律法规，平台有权视情况采取提示整改、要求删除违规内容、暂停相关调用功能等措施。\n" +
				"平台可限制部分功能、限制并发、降低额度、暂停 API Key、冻结账户、暂停服务或终止服务。\n" +
				"平台可扣除、冻结或作废违规获得、违规使用的 Token；对于因违规产生的订单，平台有权拒绝退款、取消优惠并追回损失。\n" +
				"在法律允许范围内，平台可向监管机关、司法机关、权利人或受害方提供必要信息，并采取法律允许的其他措施。\n" +
				"如用户的违约行为导致平台、关联方、合作方、其他用户或第三方遭受损失、索赔、处罚、调查、诉讼、律师费、公证费、鉴定费、差旅费等，用户应承担全部赔偿责任。",

			"第十四条 责任限制\n" +
				"在法律允许的最大范围内，平台对服务的责任以用户就产生争议的订单实际支付且未消耗、未退款的金额为上限。\n" +
				"平台不对间接损失、利润损失、商誉损失、业务中断、数据丢失、替代服务采购成本、第三方索赔、预期收益损失或惩罚性赔偿承担责任，法律强制规定另有要求的除外。\n" +
				"对于免费 Token、试用 Token、赠送 Token、活动 Token，平台可按“现状”提供，不承诺可用性、稳定性、持续性或赔偿责任，法律强制规定另有要求的除外。\n" +
				"本条不限制因平台故意、重大过失、欺诈、人身损害或法律不得限制责任的情形所产生的责任。",

			"第十五条 不可抗力\n" +
				"因自然灾害、战争、暴乱、恐怖袭击、政府行为、法律法规变化、监管要求、网络攻击、黑客事件、病毒、基础电信故障、云服务故障、能源中断、支付渠道故障、第三方模型服务中断、国际制裁、出口管制、疫情或其他不可预见、不可避免、不可克服的不可抗力事件，导致一方无法履行或迟延履行本协议的，该方在受影响范围内可部分或全部免除责任。\n" +
				"受不可抗力影响的一方应在合理期限内通知另一方，并尽合理努力减少损失。",

			"第十六条 协议变更\n" +
				"平台可根据法律法规、监管要求、业务调整或服务变化更新本协议。\n" +
				"更新后的协议将通过网站公告、支付页面、站内通知、邮件或其他合理方式向用户展示。\n" +
				"如用户在协议更新后继续购买、充值或使用 AI Token，视为用户接受更新后的协议。\n" +
				"如协议变更对用户既有已购买且未消耗 Token 的核心权益产生重大不利影响（实质性减少 Token 数量、显著缩短有效期、完全禁止原承诺模型的使用），平台应提供合理通知，并可选择提供同等价值 of 替代方案或按比例退还未使用 Token 对应费用，法律强制规定另有要求的除外。",

			"第十七条 适用法律与争议解决（境内外分设）\n" +
				"境内用户：本协议 of 订立、效力、履行、解释和争议解决，适用中华人民共和国法律（不含冲突法规则）；因本协议产生或与本协议有关的争议，双方应先友好协商解决；协商不成的，任一方应当提交服务方所在地有管辖权的人民法院诉讼解决。\n" +
				"境外用户：本协议适用服务方所在地法律（同时遵守用户所在地强制性法律）；因本协议产生或与本协议有关的争议，双方应先友好协商解决；协商不成的，提交中国国际经济贸易仲裁委员会（CIETAC），按申请仲裁时该会现行有效的仲裁规则进行仲裁；仲裁地为中国上海，仲裁语言为中文 / 英文，仲裁裁决为终局裁决，对双方均有约束力。\n" +
				"消费者用户所在地有强制性管辖规定的，从其规定。",

			"第十八条 通知\n" +
				"平台可通过网页公告、站内消息、电子邮件、短信、电话、API 控制台通知、账单通知或用户预留联系方式向用户发送通知。\n" +
				"用户应确保其联系方式真实、准确、有效；如联系方式变更，用户应及时更新。\n" +
				"通知送达规则：网页公告发布时送达；站内消息、电子邮件、短信发送成功时送达；快递或挂号信签收时送达，或因用户原因拒收、无法送达时视为送达。",

			"第十九条 其他\n" +
				"本协议构成双方关于 AI Token 购买及使用事项的完整协议，取代此前双方口头或书面达成的任何约定。\n" +
				"本协议任何条款被认定为无效、违法或不可执行的，不影响其他条款的效力。\n" +
				"未经平台书面同意，用户不得转让本协议项下权利义务。\n" +
				"平台可在业务重组、合并、分立、资产转让、关联公司承接或服务迁移时转让本协议项下权利义务，但应确保不实质性降低用户既有权益。\n" +
				"本协议标题仅为阅读方便，不影响条款解释。\n" +
				"本协议可提供中文、英文或其他语言版本；不同语言版本存在冲突的，以中文版本为准，除非用户所在地强制性法律另有要求。\n" +
				"用户在支付前应自行下载、保存或打印本协议；平台亦可在用户账户中提供协议版本、订单记录或交易记录查询。",
		},
	},
	"en": {
		Title: "AI Token Purchase and Use Agreement",
		MetaKeys: map[string]string{
			"version":     "Version: V1.1",
			"effective":   "Effective Date: {{date}}",
			"provider":    "Service Provider (Platform): OSS Energietechnik GmbH",
			"address":     "Registered Address: Adam-Opel-Straβe 16-18, 60386 Frankfurt am Main",
			"email":       "Contact Email: info@oss-energietechnik.de",
			"user":        "User: {{user}}",
			"order":       "Payment Order: {{order}} ({{amount}})",
			"signMethod":  "Signing Method: The user's checking of this agreement and clicking buttons like \"Agree and Pay\" or \"Confirm Purchase\" constitutes valid electronic signing, and the agreement takes effect immediately.",
			"applicable":  "Applicable Jurisdiction: Applicable to both domestic (Mainland China) and international users. Domestic users are preferentially subject to the laws of the People's Republic of China; international users shall also comply with mandatory laws regarding consumer protection, data protection, taxation, e-commerce, export control, and digital services in their respective locations.",
			"declaration": "Important Declaration: This agreement is a commercial contract, not a legal opinion; before going online officially, it is recommended to verify it against the user's jurisdiction, payment channels, taxation, and data compliance requirements.",
		},
		Sections: []string{
			"Article 1 Definitions\n" +
				"AI Token: Refers to the service call quota / consumption voucher purchased by the user from the platform to invoke the platform's AI models, algorithm capabilities, API services, agent services, text generation, image generation, voice processing, data analysis, or other AI functions. It does not belong to currency, virtual currency, electronic money, securities, financial products, stored value cards, or transferable assets.\n" +
				"Service: Refers to AI model calling, API integration, online generation, data processing, account management, usage statistics, technical support, and other related services provided by the platform to users.\n" +
				"Account: Refers to the account, API Key, secret key, organization ID, project ID, or other identifiers registered, logged in, or allocated by the platform for the user.\n" +
				"Input Content: Refers to the text, images, audio, video, code, files, data, prompts, interface parameters, and other content submitted, uploaded, input, transmitted, or called by the user when using the services.\n" +
				"Output Content: Refers to the text, images, audio, video, code, results of analysis, model replies, or other results generated or returned by the service based on the user's input content.\n" +
				"Consumption: Refers to the act of deducting AI Tokens due to the user calling AI services, initiating requests, generating content, processing data, or using related functions.\n" +
				"Compliance: This AI Token is strictly a platform-exclusive AI service call quota, with no circulation or legal tender status, and shall not be exchanged for fiat currency, virtual currency, or used for transaction speculation; domestic users are subject to regulatory provisions of the central bank, etc., while international users must also comply with anti-virtual currency and digital asset regulatory rules in their respective locations.",

			"Article 2 Agreement Signing and Electronic Confirmation\n" +
				"This agreement is displayed, confirmed, signed, and retained in the form of electronic data, and has the same legal effect as a written agreement.\n" +
				"The user's actions of checking and agreeing to this agreement and clicking payment, confirming purchase, submitting an order, completing payment, or actually using AI Tokens constitute valid confirmation and signing of this agreement.\n" +
				"The platform has the right to record and save the account information, order number, payment record, IP address, device information, browser information, operation time, agreement version number, checking record, click record, log records, etc., when the user confirms this agreement, which are used to prove the establishment, performance, dispute resolution, and compliance audit of the agreement.\n" +
				"The user shall not deny the effectiveness of this agreement on the grounds that no paper contract has been signed, no physical seal has been stamped, or no offline signature has been made.\n" +
				"If the user is an enterprise, institution, or other organization, the personnel purchasing or using the service on behalf of the subject confirm that they have obtained full authorization; anyone purchasing or using the service on behalf of others without authorization shall bear the corresponding responsibility themselves.\n" +
				"The terms in this agreement that exempt or limit the platform's liability, such as refunds and limitation of liability, have been prominently highlighted in bold text and independent pop-up windows on the payment page. The user confirms that they have fully understood and voluntarily accepted them.",

			"Article 3 Token Purchase, Pricing, and Delivery\n" +
				"Users may purchase AI Tokens in accordance with the packages, prices, currency, quantities, validity periods, applicable models, consumption rules, and other descriptions displayed on the platform page.\n" +
				"The price, exchange ratio, consumption rules, supported models, context length, concurrency limits, and scope of functions of AI Tokens shall be subject to the information displayed on the platform at the time the user places the order; if the platform announces adjustments, they shall be executed in accordance with the announced content.\n" +
				"Unless otherwise clearly specified by the platform, AI Tokens are limited to platform-exclusive services only, do not constitute financial assets, and shall not be transferred, traded, or exchanged for other assets.\n" +
				"After receiving the user's payment and confirming the order, the platform will credit the corresponding AI Tokens to the user's account within 24 hours; if the tokens fail to arrive within the time limit, the user has the right to unconditionally apply for a full refund.\n" +
				"The platform shall not bear any responsibility for delays in crediting accounts, deduction of transaction fees, etc., caused by payment channels, banks, third-party payment institutions, foreign exchange settlements, anti-fraud reviews, tax audits, or compliance reviews. However, the platform shall assist the user in inquiry within a reasonable scope.\n" +
				"Users shall ensure that the order information, account information, invoice information, and payment entity information are true, accurate, and complete; if incorrect user information causes recharge failure, invoice errors, account ownership disputes, or other losses, the user shall bear the responsibility themselves.",

			"Article 4 Token Usage Rules\n" +
				"The AI Tokens purchased by the user may only be used within the range designated by the platform. The specific models, services, APIs, functions, or products they can be used for shall be subject to the descriptions on the platform page, console, or order.\n" +
				"AI Tokens are consumed according to the metering methods published by the platform, including but not limited to the number of input characters, output characters, tokens, images, audio duration, video duration, request counts, model types, computing resources, storage resources, plugin calls, tool calls, etc.\n" +
				"Token consumption standards may vary for different models, functions, regional nodes, and service levels. Users should review the relevant billing instructions before use.\n" +
				"Once AI Tokens are consumed, the service is deemed to have been delivered or partially delivered; unless otherwise agreed in this agreement or mandated by law, consumed AI Tokens do not support recovery, refund, transfer, or cash exchange.\n" +
				"Users shall properly keep accounts, passwords, API Keys, access keys, verification codes, and other identity credentials; any calls initiated through the user's account or secret keys shall be deemed as actions of the user themselves or their authorized personnel.\n" +
				"If the user finds that the account or secret key has been leaked, stolen, or has abnormal calls, the user shall immediately notify the platform and take measures such as resetting the password, disabling the key, and closing the interface; the Token consumption that occurred before the platform received the notification shall, in principle, be borne by the user, except where the platform committed intentional misconduct or gross negligence.\n" +
				"Users shall not transfer, sell, rent, lend, withdraw, cash out, or trade AI Tokens outside the platform, and shall not use them for transactions or financing activities outside the platform.",

			"Article 5 Validity, Expiration, and Renewal\n" +
				"The validity period of AI Tokens shall be subject to the information explicitly displayed by the platform on the purchase page, order page, package description, or the user account dashboard; if not explicitly displayed, the default validity period is 12 months from the date the recharge credits arrive in the account.\n" +
				"Upon expiration of the validity period, unused AI Tokens will automatically become invalid, and the platform will no longer provide usage, refund, exchange, or extension services, except where mandatory laws require otherwise.\n" +
				"The platform may provide services such as extension, renewal, upgrade, and package conversion according to operational arrangements. The specific rules shall be subject to what is displayed on the platform at that time.\n" +
				"After the user purchases a new package, the order of usage for new and old Tokens shall be subject to the system rules of the platform; unless there is a special explanation, the platform may prioritize the consumption of Tokens that are closest to expiration.",

			"Article 6 Refund and Cancellation\n" +
				"Core Rules: AI Tokens are digital services delivered and consumed instantly, and the 7-day unconditional return policy does not apply; in principle, no unconditional refunds will be granted after the payment arrives in the account, except for paragraphs 2, 3, and 4 of this Article and mandatory legal provisions.\n" +
				"Domestic User Refund Rules: (1) After recharge arrives in the account, refund rules do not apply and no refund will be granted; (2) Consumed AI Tokens, promotional Tokens, activity Tokens, trial Tokens, coupon-deducted portions, Tokens restricted due to violation of this agreement, and Tokens consumed due to the user's own operational mistakes (such as accidental recharge or accidental calls) are non-refundable, non-withdrawable, non-transferable, and non-exchangeable for cash.\n" +
				"International User Refund Rules (including jurisdictions like EU, UK): (1) Where the user's location provides a statutory right of withdrawal, it takes effect after the user checks and confirms \"immediate performance, withdrawal right is lost upon recharge\"; (2) If immediate performance is not confirmed and no Tokens have been consumed, the statutory withdrawal period of the user's location applies; consumed portions are non-refundable.\n" +
				"Platform Fault Refund: Where platform reasons cause the user to be unable to use the core service corresponding to the purchased and unexpired Tokens for 30 consecutive calendar days, the platform shall, within 15 working days after receiving valid notice from the user, provide remedy measures such as compensating equivalent Tokens, extending the validity period, or refunding the unused portion; if the platform fails to remedy within the time limit, the user has the right to demand a refund in proportion to the unused Tokens.",

			"Article 7 User Compliance Obligations\n" +
				"The user promises to comply with the laws of the People's Republic of China and the laws and regulations, regulatory requirements, industry standards, export controls, economic sanctions, cross-border data transfer rules, and public order and good customs of the user's location.\n" +
				"Users shall not generate, disseminate, or assist in generating content that is illegal, infringing, fraudulent, false, hateful, harassing, violent, pornographic, terrorist, extremist, self-harm, suicide, drug-related, weapon-related, malicious code, cyberattacks, illegal financial activities, or other violating content.\n" +
				"Users shall not infringe on other people's intellectual property rights, trade secrets, privacy rights, personal information rights, portrait rights, reputation rights, or other legitimate rights and interests.\n" +
				"Users shall not, without authorization, scrape, copy, train, reverse engineer, crack, or bypass the platform's technical limits, security mechanisms, access controls, or billing systems.\n" +
				"Users shall not use the services for automated spam, click farming, fake reviews, phishing, fraud, impersonating others, bypassing content review, or batch generating illegal and violating content.\n" +
				"Users shall not submit to the platform data, personal information, sensitive personal information, confidential information, state secrets, trade secrets, or data restricted by export controls that they have no right to process.\n" +
				"The platform has the right to conduct necessary security reviews, risk control identification, and compliance processing on user input content, output content, call behavior, and account behavior in accordance with laws and regulations, regulatory requirements, platform rules, risk control policies, or requirements of third-party model providers.",

			"Article 8 AI Output Content and Risk Disclaimer\n" +
				"The user understands and agrees that AI output content is generated by algorithm models based on input content, model parameters, training data, context, and system configurations, and may contain inaccuracies, incompleteness, obsolescence, fabrication, biases, inapplicability, or unexplainability.\n" +
				"Output content does not constitute legal, financial, medical, investment, engineering, safety, compliance, or other professional advice; if the user needs to rely on the output content to make important decisions, they shall conduct manual review on their own and consult qualified professionals.\n" +
				"The user shall bear full responsibility for their input content, usage methods, and the review, publication, dissemination, application, and consequences of the output content.\n" +
				"The platform does not guarantee that the AI output content is completely accurate, that the service is continuously available, or that it meets the user's specific purposes; nor does it guarantee that the output content will not be identical or similar to content generated by other users.\n" +
				"Where the user uses the output content for commercial purposes, public release, advertisements, bidding documents, legal documents, medical advice, investment advice, engineering designs, automated decision-making, or other high-risk scenarios, the user shall bear the compliance review and usage risks themselves.",

			"Article 9 Intellectual Property Rights\n" +
				"The platform and its affiliates and licensors enjoy complete rights to the platform system, software, models, algorithms, interfaces, documents, page designs, trademarks, logos, databases, technical solutions, business models, and related intellectual property rights.\n" +
				"The user's purchase of AI Tokens does not grant ownership of the platform's models, algorithms, software, source code, system architecture, trademarks, or other intellectual property rights.\n" +
				"The user retains the rights legally enjoyed by them over the input content legally provided by them.\n" +
				"On the premise that the user complies with this agreement and has paid the corresponding fees, the platform permits the user to use the output content generated by the service within the scope of its rights. However, whether the output content is exclusive, whether it infringes third-party rights, whether it can be registered as intellectual property, or whether it can be commercialized shall be judged by the user themselves in combination with the specific content and applicable laws, and the risk of infringement shall be borne by the user themselves.\n" +
				"The user authorizes the platform to process input content, output content, and usage data within the scope necessary for providing services, billing, risk control, security reviews, troubleshooting, service optimization, compliance audits, and fulfilling legal obligations.\n" +
				"If the platform uses third-party base models, cloud services, payment services, or other third-party services, the relevant rights limitations, content policies, and usage restrictions may apply to the user simultaneously.",

			"Article 10 Data Protection and Privacy\n" +
				"Domestic Users: Governed by the \"Personal Information Protection Law of the People's Republic of China\", adhering to the principle of minimization and necessity. The platform only collects information necessary to provide the service: (1) Registration info (phone number, email); (2) Payment info (Alipay order number, payment amount); (3) Device info (IP address, browser type); (4) Usage info (Token recharge / consumption records). It does not collect sensitive information such as ID numbers or credit card numbers.\n" +
				"International Users: Governed by the data protection laws of their respective locations (such as GDPR). The platform processes personal information in accordance with the legal requirements of the user's location, does not transmit data across borders to regions that do not meet data protection standards, and may sign a data processing agreement with the user when necessary.\n" +
				"User Personal Information Retention Period: 1 year from the date of account deletion or service termination, and will be automatically anonymized or deleted upon expiration; order records, tax records, and Token recharge / consumption records are retained for 3 years (to accommodate tax and compliance audit requirements).\n" +
				"If the user inputs personal information, sensitive personal information, trade secrets, confidential information, or regulated data involving others to the platform, the user shall ensure that they have a legal basis for processing and have fulfilled the necessary notification, consent, authorization, assessment, or filing obligations.\n" +
				"La platform may, within the necessary scope, entrust compliant third-party service providers to process data, including cloud providers, model providers, payment institutions, SMS service providers, invoicing service providers, customer support systems, data security service providers, etc., and require relevant third parties to take reasonable data protection measures.",

			"Article 11 Taxes, Invoices, and Cross-Border Payments\n" +
				"Users shall pay fees in accordance with the amount displayed on the order page. Whether the relevant fees include tax shall be subject to the platform page, order, or invoice descriptions.\n" +
				"Extra fees generated due to the user's location, payment method, payment channel, acquiring institution, bank, foreign exchange settlement, VAT, consumption tax, sales tax, withholding tax, digital service tax, or other taxes and fees shall be borne by the parties in accordance with applicable laws and order agreements. For payments involving foreign currencies, the actual deducted amount may differ due to exchange rates, bank transaction fees, payment channel fees, or settlement times; such differences shall be borne by the user themselves, unless otherwise promised by the platform.\n" +
				"If the user needs invoices, receipts, or tax certificates, they shall provide true, accurate, and complete invoicing information in accordance with the platform's requirements; if invoicing fails or invoices need to be corrected due to incorrect information provided by the user or inconsistency between the paying entity and the invoiced entity, the user shall cooperate with the processing and bear the corresponding losses.",

			"Article 12 Service Modification, Interruption, and Maintenance\n" +
				"The platform may adjust service contents, model types, billing rules, Token consumption ratios, usage limits, and function scopes in accordance with business development, technical upgrades, model supply, cost changes, laws and regulations, regulatory requirements, or changes in third-party services. However, if the aforementioned changes cause a significant reduction in the usable scope of Tokens purchased but unused by the user, or a significant rise in the unit service cost (Token consumption ratio increases by more than 20%), the user has the right to apply for a pro-rata refund of the fee corresponding to the unused Tokens before the changes take effect; the aforementioned changes do not affect the usage rights of purchased and unconsumed Tokens within the original service scope.\n" +
				"When making major changes, the platform shall notify the user through website announcements, system notifications, emails, SMS, or other reasonable methods.\n" +
				"Where services are interrupted, delayed, or unavailable due to system maintenance, upgrades, network failures, cloud service failures, third-party model service failures, payment channel failures, security incidents, force majeure, or regulatory requirements, the platform shall make reasonable efforts to restore services, and shall not bear liability for indirect losses.\n" +
				"The platform shall not bear responsibility for usage failures caused by the user's network, equipment, interface configurations, key management, code errors, call parameter errors, violation of usage rules, or third-party reasons.",

			"Article 13 Account Restrictions and Breach Treatment\n" +
				"If the user violates this agreement or applicable laws and regulations, the platform has the right to take measures such as prompting rectification, demanding the deletion of violating content, and suspending related call functions depending on the situation.\n" +
				"The platform may restrict partial functions, limit concurrency, lower quotas, suspend API Keys, freeze accounts, suspend services, or terminate services.\n" +
				"The platform may deduct, freeze, or invalidate Tokens obtained or used in violation of regulations; for orders generated due to violations, the platform has the right to refuse refunds, cancel discounts, and recover losses.\n" +
				"Within the scope permitted by law, the platform may provide necessary information to regulatory authorities, judicial organs, rights holders, or victims, and take other measures permitted by law.\n" +
				"If the user's breach causes the platform, affiliates, partners, other users, or third parties to suffer losses, claims, penalties, investigations, lawsuits, attorney fees, notary fees, appraisal fees, travel expenses, etc., the user shall bear full compensation liability.",

			"Article 14 Limitation of Liability\n" +
				"To the maximum extent permitted by law, the platform's liability for services is capped at the amount actually paid but unconsumed and unrefunded by the user for the disputed order.\n" +
				"The platform does not bear responsibility for indirect losses, loss of profits, loss of goodwill, business interruption, data loss, procurement costs of substitute services, third-party claims, loss of expected revenue, or punitive damages, except where mandatory laws require otherwise.\n" +
				"Regarding free Tokens, trial Tokens, promotional Tokens, and activity Tokens, the platform provides them \"as is\" and does not promise availability, stability, continuity, or compensation liability, except where mandatory laws require otherwise.\n" +
				"This Article does not limit liability arising from the platform's intentional misconduct, gross negligence, fraud, personal injury, or situations where liability cannot be limited by law.",

			"Article 15 Force Majeure\n" +
				"If a party is unable to perform or delays in performing this agreement due to natural disasters, war, riots, terrorist attacks, government actions, changes in laws and regulations, regulatory requirements, cyberattacks, hacker incidents, viruses, basic telecommunications failures, cloud service failures, energy outages, payment channel failures, third-party model service interruptions, international sanctions, export controls, epidemics, or other unforeseen, unavoidable, and insurmountable force majeure events, the party may be partially or completely exempted from liability within the affected scope.\n" +
				"The party affected by force majeure shall notify the other party within a reasonable period and make reasonable efforts to mitigate losses.",

			"Article 16 Agreement Modification\n" +
				"The platform may update this agreement in accordance with laws and regulations, regulatory requirements, business adjustments, or service changes.\n" +
				"The updated agreement will be displayed to the user through website announcements, payment pages, system notifications, emails, or other reasonable methods.\n" +
				"If the user continues to purchase, recharge, or use AI Tokens after the agreement is updated, the user is deemed to have accepted the updated agreement.\n" +
				"If the agreement modification has a major adverse effect on the user's existing core rights and interests of purchased and unconsumed Tokens (substantially reducing the Token quantity, significantly shortening the validity period, or completely prohibiting the use of models originally promised), the platform shall provide reasonable notice, and may choose to provide an alternative solution of equal value or refund the fee corresponding to the unused Tokens on a pro-rata basis, except where mandatory laws require otherwise.",

			"Article 17 Applicable Law and Dispute Resolution\n" +
				"Domestic Users: The establishment, effectiveness, performance, interpretation, and dispute resolution of this agreement shall apply the laws of the People's Republic of China (excluding conflict of laws rules); disputes arising from or in connection with this agreement shall first be resolved through friendly negotiation; if negotiation fails, either party shall submit the dispute to the People's Court with jurisdiction in the location of the service provider for litigation.\n" +
				"International Users: This agreement is governed by the laws of the location of the service provider (while also complying with mandatory laws in the user's location); disputes arising from or in connection with this agreement shall first be resolved through friendly negotiation; if negotiation fails, they shall be submitted to the China International Economic and Trade Arbitration Commission (CIETAC) for arbitration in accordance with the arbitration rules in force at the time of applying for arbitration; the place of arbitration shall be Shanghai, China, the language of arbitration shall be Chinese / English, and the arbitral award shall be final and binding on both parties.\n" +
				"Where the location of the consumer user has mandatory provisions on jurisdiction, those provisions shall prevail.",

			"Article 18 Notices\n" +
				"The platform may send notices to users through website announcements, system messages, emails, SMS, phone calls, API console notifications, billing notifications, or contact information reserved by the user.\n" +
				"Users shall ensure that their contact information is true, accurate, and effective; if the contact information changes, the user shall update it in a timely manner.\n" +
				"Notice Delivery Rules: Deemed delivered when website announcements are published; deemed delivered when system messages, emails, or SMS are sent successfully; deemed delivered upon signature when sent by express delivery or registered mail, or deemed delivered when refused or unable to be delivered due to the user's reasons.",

			"Article 19 Miscellaneous\n" +
				"This agreement constitutes the entire agreement between the parties regarding the purchase and use of AI Tokens, replacing any prior oral or written agreements reached between the parties.\n" +
				"If any provision of this agreement is determined to be invalid, illegal, or unenforceable, the validity of the remaining provisions shall not be affected.\n" +
				"Without the prior written consent of the platform, the user shall not transfer their rights and obligations under this agreement.\n" +
				"The platform may transfer its rights and obligations under this agreement in the event of business reorganization, merger, division, asset transfer, assumption by affiliates, or service migration, but shall ensure that the existing rights and interests of users are not substantially reduced.\n" +
				"The headings in this agreement are for convenience of reading only and do not affect the interpretation of the clauses.\n" +
				"This agreement may be provided in Chinese, English, or other language versions; if there is a conflict between different language versions, the Chinese version shall prevail, unless mandatory laws in the user's location require otherwise.\n" +
				"Users shall download, save, or print this agreement on their own before payment; the platform may also provide queries for agreement versions, order records, or transaction records in the user's account.",
		},
	},
	"fr": {
		Title: "Contrat d'Achat et d'Utilisation de Jetons d'IA",
		MetaKeys: map[string]string{
			"version":     "Version : V1.1",
			"effective":   "Date d'effet : {{date}}",
			"provider":    "Prestataire de services (Plateforme) : OSS Energietechnik GmbH",
			"address":     "Adresse enregistrée : Adam-Opel-Straβe 16-18, 60386 Frankfurt am Main",
			"email":       "E-mail de contact : info@oss-energietechnik.de",
			"user":        "Utilisateur : {{user}}",
			"order":       "Commande de paiement : {{order}} ({{amount}})",
			"signMethod":  "Méthode de signature : Le fait pour l'utilisateur de cocher cet accord et de cliquer sur des boutons similaires à \"Accepter et Payer\" ou \"Confirmer l'achat\" constitue une signature électronique valide, et le contrat prend effet immédiatement.",
			"applicable":  "Zone d'application : Applicable aux utilisateurs nationaux (Chine continentale) et internationaux. Le droit de la République populaire de Chine s'applique en priorité aux utilisateurs nationaux ; les utilisateurs internationaux doivent également se conformer aux lois obligatoires de leur lieu de résidence concernant la protection des consommateurs, la protection des données, la fiscalité, le commerce électronique, le contrôle des exportations et les services numériques.",
			"declaration": "Déclaration importante : Ce contrat est un document commercial et non un avis juridique ; avant la mise en ligne officielle, il est recommandé de procéder à une validation finale tenant compte de la juridiction de l'utilisateur, des canaux de paiement, des exigences fiscales et de conformité des données.",
		},
		Sections: []string{
			"Article 1 Définitions\n" +
				"Jetons d'IA (AI Token) : Désigne le quota d'appel de service / bon de consommation acheté par l'utilisateur auprès de la plateforme pour invoquer les modèles d'IA, les capacités algorithmiques, les services API, les services d'agent, la génération de texte, la génération d'images, le traitement de la voix, l'analyse de données ou d'autres fonctions d'IA de la plateforme. Il ne s'agit pas de monnaie, de monnaie virtuelle, de monnaie électronique, de valeurs mobilières, de produits financiers, de cartes à valeur stockée ou d'actifs transférables.\n" +
				"Service : Désigne l'appel de modèle d'IA, l'intégration d'API, la génération en ligne, le traitement de données, la gestion de compte, les statistiques d'utilisation, le support technique et d'autres services connexes fournis par la plateforme aux utilisateurs.\n" +
				"Compte : Désigne le compte, la clé API (API Key), la clé secrète, l'identifiant d'organisation, l'identifiant de projet ou tout autre identifiant enregistré, connecté ou attribué par la plateforme pour l'utilisateur.\n" +
				"Contenu d'entrée (Input) : Désigne le texte, les images, l'audio, la vidéo, le code, les fichiers, les données, les invites (prompts), les paramètres d'interface ou tout autre contenu soumis, téléchargé, saisi, transmis ou appelé par l'utilisateur lors de l'utilisation des services.\n" +
				"Contenu de sortie (Output) : Désigne le texte, les images, l'audio, la vidéo, le code, les résultats d'analyse, les réponses du modèle ou tout autre résultat généré ou renvoyé par le service sur la base du contenu d'entrée de l'utilisateur.\n" +
				"Consommation : Désigne l'action de déduire des jetons d'IA (AI Token) en raison de l'appel de services d'IA, du lancement de requêtes, de la génération de contenu, du traitement de données ou de l'utilisation de fonctions connexes par l'utilisateur.\n" +
				"Conformité : Ce jeton d'IA (AI Token) est uniquement un quota d'appel de service d'IA exclusif à la plateforme, sans circulation ni cours légal, et ne peut être échangé contre de la monnaie fiduciaire, de la monnaie virtuelle ou utilisé pour la spéculation commerciale ; le droit national applique les réglementations de la banque centrale, etc., et l'international est conforme aux règles de réglementation des monnaies virtuelles et des actifs numériques du lieu de résidence de l'utilisateur.",

			"Article 2 Signature et validation\n" +
				"Cet accord est présenté, confirmé, signé et conservé sous forme de données électroniques, et a la même force juridique qu'un accord écrit.\n" +
				"L'action de l'utilisateur consistant à cocher son accord, à cliquer sur payer, à confirmer l'achat, à soumettre la commande, à effectuer le paiement ou à utiliser réellement les jetons d'IA constitue une confirmation et une signature valides du présent contrat.\n" +
				"La plateforme a le droit d'enregistrer et de conserver les informations de compte de l'utilisateur, le numéro de commande, l'enregistrement de paiement, l'adresse IP, les informations sur l'appareil, les informations sur le navigateur, l'heure de l'opération, le numéro de version de l'accord, les enregistrements de coche, de clic, de journal, etc., lors de la confirmation du présent contrat, afin de prouver la conclusion du contrat, son exécution, le règlement des différends et les audits de conformité.\n" +
				"L'utilisateur ne peut nier la validité du présent accord sous prétexte qu'aucun contrat papier n'a été signé, qu'aucun sceau physique n'a été apposé ou qu'aucune signature hors ligne n'a été effectuée.\n" +
				"Si l'utilisateur est une entreprise, une institution ou une autre organisation, la personne effectuant l'achat ou utilisant le service au nom de cette entité confirme qu'elle a obtenu une autorisation complète ; l'achat ou l'utilisation du service au nom d'un tiers sans autorisation engage la responsabilité personnelle de la personne concernée.\n" +
				"Les clauses du présent contrat qui exonèrent ou limitent la responsabilité de la plateforme, telles que le remboursement et la limitation de responsabilité, ont été mises en évidence en gras dans le texte et via des fenêtres pop-up indépendantes sur la page de paiement. L'utilisateur confirme qu'il les a pleinement comprises et acceptées volontairement.",

			"Article 3 Achat, tarification et livraison\n" +
				"Les utilisateurs peuvent acheter des jetons d'IA (AI Token) conformément aux forfaits, prix, devises, quantités, périodes de validité, modèles applicables, règles de consommation et autres descriptions affichées sur la page de la plateforme.\n" +
				"Le prix, le taux de conversion, les règles de consommation, les modèles pris en charge, la longueur du contexte, les limites de concurrence et l'étendue des fonctions des jetons d'IA sont soumis aux informations affichées sur la plateforme au moment où l'utilisateur passe la commande ; en cas d'ajustements annoncés par la plateforme, ils seront appliqués conformément au contenu annoncé.\n" +
				"Sauf mention contraire expresse de la plateforme, les jetons d'IA (AI Token) sont limités à l'utilisation au sein des services exclusifs de la plateforme, ne constituent pas des actifs financiers et ne peuvent pas être transférés, échangés ou convertis en d'autres actifs.\n" +
				"Après réception du paiement de l'utilisateur et confirmation de la commande, la plateforme créditera les jetons d'IA (AI Token) correspondants sur le compte de l'utilisateur dans les 24 heures ; en cas de non-crédit dans ce délai, l'utilisateur a le droit de demander un remboursement intégral sans condition.\n" +
				"La plateforme n'assume aucune responsabilité pour les retards de crédit, les déductions de frais de transaction, etc., causés par les canaux de paiement, les banques, les institutions de paiement tierces, les règlements de change, les contrôles anti-fraude, les audits fiscaux ou de conformité. Toutefois, elle assistera l'utilisateur dans ses démarches d'information dans des limites raisonnables.\n" +
				"Les utilisateurs doivent s'assurer que les informations de commande, de compte, de facturation et de l'entité de paiement sont réelles, exactes et complètes ; si des informations utilisateur erronées entraînent un échec de recharge, des erreurs de facturation, des litiges de propriété de compte ou d'autres pertes, l'utilisateur en assumera l'entière responsabilité.",

			"Article 4 Règles d'utilisation des jetons\n" +
				"Les jetons d'IA (AI Token) achetés par l'utilisateur ne peuvent être utilisés que dans la limite spécifiée par la plateforme ; les modèles, services, API, fonctions ou produits spécifiques pour lesquels ils peuvent être utilisés sont soumis aux descriptions figurant sur la page de la plateforme, la console ou la commande.\n" +
				"Les jetons d'IA (AI Token) sont consommés conformément aux méthodes de mesure publiées par la plateforme, y compris, mais sans s'y limiter, le nombre de caractères d'entrée, de caractères de sortie, de jetons, le nombre d'images, la durée de l'audio, la durée de la vidéo, le nombre de requêtes, le type de modèle, les ressources de calcul, les ressources de stockage, les appels de plug-ins, les appels d'outils, etc.\n" +
				"Les normes de consommation de jetons peuvent varier selon les modèles, les fonctions, les nœuds régionaux et les niveaux de service. Les utilisateurs doivent consulter les instructions de facturation correspondantes avant toute utilisation.\n" +
				"Dès que les jetons d'IA (AI Token) consommés, le service est considéré comme livré ou partiellement livré ; sauf convention contraire dans le présent accord ou disposition légale obligatoire, les jetons d'IA consommés ne peuvent pas être restaurés, remboursés, transférés ou convertis en espèces.\n" +
				"Les utilisateurs doivent conserver en toute sécurité leurs comptes, mots de passe, clés API (API Key), clés d'accès, codes de vérification et autres identifiants ; tous les appels initiés via le compte ou la clé de l'utilisateur sont considérés comme des actions de l'utilisateur lui-même ou de son personnel autorisé.\n" +
				"Si l'utilisateur constate que son compte ou sa clé a été divulgué, usurpé ou fait l'objet d'appels anormaux, il doit immédiatement en informer la plateforme et prendre des mesures telles que la réinitialisation du mot de passe, la désactivation de la clé ou la fermeture de l'interface ; la consommation de jetons survenue avant que la plateforme ne reçoive la notification est en principe à la charge de l'utilisateur, sauf faute intentionnelle ou négligence grave de la plateforme.\n" +
				"L'utilisateur ne doit pas transférer, vendre, louer, prêter, retirer, convertir en espèces ou échanger des jetons d'IA (AI Token) en dehors de la plateforme, ni les utiliser pour des transactions ou des activités de financement en dehors de la plateforme.",

			"Article 5 Validité, expiration et renouvellement\n" +
				"La période de validité des jetons d'IA (AI Token) est soumise aux informations explicitement affichées par la plateforme sur la page d'achat, la page de commande, la description du forfait ou l'espace client ; à défaut d'affichage explicite, la période de validité par défaut est de 12 mois à compter de la date de réception de la recharge.\n" +
				"À l'expiration de la période de validité, les jetons d'IA (AI Token) non utilisés expireront automatiquement, et la plateforme ne fournira plus de services d'utilisation, de remboursement, d'échange ou de prolongation, sauf disposition contraire de la loi impérative.\n" +
				"La plateforme peut proposer des services de prolongation, de renouvellement, de mise à niveau, de changement de forfait, etc., selon ses plans opérationnels. Les règles spécifiques seront soumises à ce qui sera affiché sur la plateforme à ce moment-là.\n" +
				"Après l'achat d'un nouveau forfait par l'utilisateur, l'ordre d'utilisation des jetons anciens et nouveaux est soumis aux règles système de la plateforme ; sauf indication contraire, la plateforme peut consommer en priorité les jetons arrivant bientôt à expiration.",

			"Article 6 Politique de remboursement et de rétractation\n" +
				"Règle de base : Les jetons d'IA (AI Token) sont des services numériques livrés et consommés instantanément, et ne bénéficient pas du droit de rétractation de sept jours sans motif ; aucun remboursement sans motif ne sera accordé après réception du paiement, sauf pour les paragraphes 2, 3 et 4 de cet article et les dispositions légales obligatoires.\n" +
				"Règles de remboursement pour les utilisateurs nationaux : (1) Aucun droit de remboursement ne s'applique après la recharge ; (2) Les jetons d'IA consommés, offerts, promotionnels, d'essai, les portions déduites par coupon, les jetons restreints pour violation du contrat, et les jetons consommés par erreur de manipulation de l'utilisateur (ex. recharge ou appel erroné) ne peuvent faire l'objet de remboursement, retrait, transfert ou conversion en espèces.\n" +
				"Règles de remboursement pour les utilisateurs internationaux (y compris l'UE, le Royaume-Uni, etc.) : (1) Si l'utilisateur dispose d'un droit de rétractation légal dans son pays, celui-ci prend fin après que l'utilisateur a confirmé par coche « exécution immédiate, perte du droit de rétractation dès la recharge » ; (2) En l'absence de confirmation d'exécution immédiate et si aucun jeton n'a été consommé, le délai légal de rétractation s'applique ; la portion consommée reste non remboursable.\n" +
				"Remboursement pour faute de la plateforme : Si des raisons imputables à la plateforme empêchent l'utilisateur d'utiliser les services de base correspondant aux jetons achetés et non expirés pendant 30 jours calendaires consécutifs, la plateforme doit, dans les 15 jours ouvrables suivant la notification de l'utilisateur, proposer des mesures de remédiation telles que la compensation par des jetons équivalents, la prolongation de la validité ou le remboursement de la portion inutilisée ; à défaut, l'utilisateur a le droit de demander un remboursement au prorata des jetons inutilisés.",

			"Article 7 Obligations de conformité de l'utilisateur\n" +
				"L'utilisateur s'engage à respecter les lois de la République populaire de Chine ainsi que les lois, réglementations, exigences de contrôle, normes industrielles, contrôles des exportations, sanctions économiques, règles de transfert transfrontalier de données et l'ordre public et les bonnes mœurs de son lieu de résidence.\n" +
				"L'utilisateur ne doit pas générer, propager ou aider à générer du contenu illégal, contrefait, frauduleux, mensonger, haineux, harcelant, violent, pornographique, terroriste, extrémiste, incitant à l'automutilation, au suicide, lié à la drogue, aux armes, aux codes malveillants, aux cyberattaques, aux activités financières illégales ou tout autre contenu non conforme.\n" +
				"L'utilisateur ne doit pas enfreindre les droits de propriété intellectuelle d'autrui, les secrets commerciaux, le droit à la vie privée, les droits d'information personnelle, le droit à l'image, le droit à la réputation ou d'autres droits et intérêts légitimes.\n" +
				"L'utilisateur ne doit pas, sans autorisation, collecter (scraping), copier, entraîner, faire de l'ingénierie inverse, pirater, contourner les limitations techniques, les mécanismes de sécurité, les contrôles d'accès ou les systèmes de facturation de la plateforme.\n" +
				"L'utilisateur ne doit pas utiliser le service pour l'envoi automatisé de spams, l'augmentation artificielle de trafic (刷量), les faux avis, l'hameçonnage (phishing), l'escroquerie, l'usurpation d'identité, le contournement de la modération de contenu ou la génération en masse de contenus illégaux.\n" +
				"L'utilisateur ne doit pas soumettre à la plateforme des données, des informations personnelles, des informations personnelles sensibles, des informations confidentielles, des secrets d'État, des secrets commerciaux ou des données soumises à des restrictions de contrôle des exportations qu'il n'a pas le droit de traiter.\n" +
				"La plateforme a le droit d'effectuer les examens de sécurité, l'identification des risques et le traitement de conformité nécessaires sur le contenu d'entrée, le contenu de sortie, le comportement d'appel et le comportement de compte de l'utilisateur conformément aux lois et réglementations, aux exigences réglementaires, aux règles de la plateforme, aux politiques de gestion des risques ou aux exigences des fournisseurs tiers.",

			"Article 8 Risques et exclusions de responsabilité\n" +
				"L'utilisateur comprend et accepte que le contenu de sortie de l'IA est généré par des modèles algorithmiques sur la base du contenu d'entrée, des paramètres du modèle, des données d'entraînement, du contexte et des configurations du système, et peut s'avérer inexact, incomplet, obsolète, fictif, biaisé, inapplicable ou inexplicable.\n" +
				"Le contenu généré ne constitue pas un avis juridique, financier, médical, d'investissement, d'ingénierie, de sécurité, de conformité ou tout autre avis professionnel ; si l'utilisateur doit s'appuyer sur le contenu généré pour prendre des décisions importantes, il doit procéder lui-même à une vérification humaine et consulter des professionnels qualifiés.\n" +
				"L'utilisateur assume l'entière responsabilité du contenu qu'il saisit, de sa manière d'utiliser le service, ainsi que de l'examen, de la publication, de la diffusion, de l'application et des conséquences du contenu généré.\n" +
				"La plateforme ne garantit pas que le contenu généré par l'IA soit tout à fait exact, que le service soit disponible en continu, qu'il réponde à des fins spécifiques de l'utilisateur, ni que le contenu généré ne soit identique ou similaire au contenu généré par d'autres utilisateurs.\n" +
				"En cas d'utilisation du contenu de sortie à des fins commerciales, de publication publique, de publicité, de documents d'appel d'offres, de documents juridiques, de conseils médicaux, de conseils financiers, de conceptions techniques, de prise de décision automatisée ou dans d'autres cas à haut risque, l'utilisateur en assumera le contrôle de conformité et les risques d'utilisation.",

			"Article 9 Propriété intellectuelle\n" +
				"La plateforme, ses sociétés affiliées et ses concédants de licence détiennent l'intégralité des droits sur le système de la plateforme, les logiciels, les modèles, les algorithmes, les interfaces, la documentation, la conception des pages, les marques, les logos, les bases de données, les solutions techniques, les modèles commerciaux et les droits de propriété intellectuelle y afférents.\n" +
				"L'achat de jetons d'IA (AI Token) par l'utilisateur ne lui confère aucun droit de propriété sur les modèles, algorithmes, logiciels, codes sources, architecture système, marques ou autres droits de propriété intellectuelle de la plateforme.\n" +
				"L'utilisateur conserve les droits qu'il détient légalement sur le contenu d'entrée fourni par ses soins de manière licite.\n" +
				"Sous réserve que l'utilisateur respecte le présent contrat et ait payé les frais correspondants, la plateforme l'autorise à utiliser le contenu généré par le service dans la limite de ses droits. Cependant, la question de savoir si le contenu généré peut être exclusif, s'il enfreint les droits de tiers, s'il peut être enregistré en tant que propriété intellectuelle ou s'il peut être commercialisé doit être évaluée par l'utilisateur lui-même en fonction du contenu spécifique et des lois applicables ; l'utilisateur assume seul le risque de contrefaçon.\n" +
				"L'utilisateur autorise la plateforme à traiter le contenu d'entrée, le contenu de sortie et les données d'utilisation dans la mesure nécessaire à la fourniture des services, à la facturation, à la gestion des risques, aux examens de sécurité, au dépannage, à l'optimisation des services, aux audits de conformité et à l'exécution des obligations légales.\n" +
				"Si la plateforme utilise des modèles de base tiers, des services cloud, des services de paiement ou d'autres services tiers, les limitations de droits, politiques de contenu et restrictions d'utilisation correspondantes peuvent également s'appliquer à l'utilisateur.",

			"Article 10 Protection des données et vie privée\n" +
				"Utilisateurs nationaux : Conformément à la « Loi sur la protection des informations personnelles de la République populaire de Chine », selon le principe de minimisation et de nécessité, la plateforme ne collecte que les informations strictement nécessaires à la fourniture du service : ① informations d'inscription (numéro de téléphone, e-mail) ; ② informations de paiement (numéro de commande Alipay, montant du paiement) ; ③ informations sur l'appareil (adresse IP, type de navigateur) ; ④ informations d'utilisation (historique de recharge / consommation de jetons). Aucune information sensible telle que le numéro d'identité ou le numéro de carte bancaire n'est collectée.\n" +
				"Utilisateurs internationaux : Soumis aux lois sur la protection des données de leur lieu de résidence (telles que le RGPD). La plateforme traite les informations personnelles conformément aux exigences légales locales, ne transfère pas de données à l'étranger vers des régions ne répondant pas aux normes de protection des données et peut, si nécessaire, signer un accord de traitement des données avec l'utilisateur.\n" +
				"Durée de conservation des données : Les informations personnelles de l'utilisateur sont conservées pendant 1 an à compter de la suppression du compte ou de la résiliation du service, après quoi elles sont automatiquement anonymisées ou supprimées ; les enregistrements de commandes, les registres fiscaux et les historiques de recharge / consommation de jetons sont conservés pendant 3 ans (pour répondre aux exigences d'audit fiscal et de conformité).\n" +
				"Si l'utilisateur saisit sur la plateforme des informations personnelles, des informations personnelles sensibles, des secrets commerciaux, des informations confidentielles ou des données réglementées concernant des tiers, il doit s'assurer de disposer d'une base légale pour ce traitement et d'avoir accompli les obligations nécessaires d'information, de consentement, d'autorisation, d'évaluation ou d'enregistrement.\n" +
				"La plateforme peut, dans la mesure nécessaire, confier le traitement des données à des prestataires de services tiers conformes, y compris des fournisseurs de services cloud, des fournisseurs de modèles, des institutions de paiement, des fournisseurs de SMS, des prestataires de facturation, des systèmes de support client, des prestataires de sécurité des données, etc., et exige de ces tiers qu'ils prennent des mesures raisonnables de protection des données.",

			"Article 11 Taxes, factures et paiements transfrontaliers\n" +
				"L'utilisateur doit payer les frais conformément au montant indiqué sur la page de commande. Le fait que les frais correspondants incluent ou non des taxes est soumis aux indications de la page de la plateforme, de la commande ou de la facture.\n" +
				"Les frais supplémentaires générés en raison de la localisation de l'utilisateur, du mode de paiement, du canal de paiement, de l'acquéreur, de la banque, du règlement de change, de la TVA, de la taxe à la consommation, de la taxe de vente, de la retenue à la source, de la taxe sur les services numériques ou d'autres taxes et frais sont supportés par les parties conformément aux lois applicables et aux conditions de la commande. Pour les paiements en devises étrangères, le montant réel débité peut différer en raison des taux de change, des frais bancaires, des frais de canal de paiement ou des délais de règlement ; ces différences sont à la charge exclusive de l'utilisateur, sauf promesse contraire de la plateforme.\n" +
				"Si l'utilisateur a besoin de factures, de reçus ou de justificatifs fiscaux, il doit fournir des informations de facturation réelles, exactes et complètes conformément aux exigences de la plateforme ; si des informations erronées fournies par l'utilisateur ou une incohérence entre l'entité payeuse et l'entité facturée empêchent l'établissement de la facture ou nécessitent sa correction, l'utilisateur doit coopérer et assumer les pertes correspondantes.",

			"Article 12 Modification, interruption et maintenance du service\n" +
				"La plateforme peut ajuster le contenu du service, les types de modèles, les règles de tarification, les ratios de consommation de jetons, les limites d'utilisation et l'étendue des fonctions en fonction du développement de l'activité, des mises à niveau techniques, de la fourniture de modèles, des variations de coûts, des lois et réglementations, des exigences réglementaires ou des modifications des services tiers. Toutefois, si ces modifications entraînent une réduction significative de la portée d'utilisation des jetons achetés mais inutilisés, ou une augmentation significative du coût unitaire du service (le ratio de consommation de jetons augmentant de plus de 20 %), l'utilisateur a le droit de demander un remboursement au prorata des frais correspondants aux jetons inutilisés avant que les modifications ne prennent effet ; lesdites modifications n'affectent pas les droits d'utilisation des jetons achetés mais non consommés dans la portée initiale du service.\n" +
				"Lors de modifications majeures, la plateforme doit en informer l'utilisateur par le biais d'annonces sur le site, de notifications internes, d'e-mails, de SMS ou d'autres moyens raisonnables.\n" +
				"En cas d'interruption, de retard ou d'indisponibilité du service en raison de la maintenance du système, de mises à niveau, de pannes de réseau, de pannes de services cloud, de pannes de services de modèles tiers, de pannes de canaux de paiement, d'incidents de sécurité, de force majeure ou d'exigences réglementaires, la plateforme fera des efforts raisonnables pour restaurer le service, mais n'assumera aucune responsabilité pour les pertes indirectes.\n" +
				"La plateforme n'assume aucune responsabilité pour les échecs d'utilisation causés par le réseau de l'utilisateur, ses équipements, les configurations d'interface, la gestion des clés, les erreurs de code, les erreurs de paramètres d'appel, la violation des règles d'utilisation ou des raisons imputables à des tiers.",

			"Article 13 Restrictions de compte et traitement des violations\n" +
				"Si l'utilisateur viole le présent contrat ou les lois et réglementations applicables, la plateforme a le droit, selon les circonstances, de prendre des mesures telles que demander une rectification, exiger la suppression du contenu non conforme ou suspendre les fonctions d'appel associées.\n" +
				"La plateforme peut restreindre certaines fonctions, limiter la concurrence, réduire les quotas, suspendre les clés API, geler les comptes, suspendre le service ou résilier le service.\n" +
				"La plateforme peut déduire, geler ou annuler les jetons obtenus ou utilisés de manière illicite ; pour les commandes générées en violation des règles, la plateforme a le droit de refuser le remboursement, d'annuler les réductions et de récupérer les pertes.\n" +
				"Dans la mesure permise par la loi, la plateforme peut fournir les informations nécessaires aux autorités de régulation, aux organes judiciaires, aux titulaires de droits ou aux parties lésées, et prendre d'autres mesures autorisées par la loi.\n" +
				"Si le comportement de violation de l'utilisateur cause des pertes, des réclamations, des sanctions, des enquêtes, des litiges, des frais d'avocat, des frais de notaire, des frais d'expertise, des frais de déplacement ou d'autres frais à la plateforme, à ses affiliés, partenaires, autres utilisateurs ou tiers, l'utilisateur assumera l'entière responsabilité de l'indemnisation.",

			"Article 14 Limitation de responsabilité\n" +
				"Dans la mesure maximale permise par la loi applicable, la responsabilité totale de la plateforme pour le service est plafonnée au montant réellement payé par l'utilisateur pour la commande contestée, non consommé et non remboursé.\n" +
				"La plateforme n'assume aucune responsabilité pour les dommages indirects, la perte de bénéfices, la perte de clientèle, l'interruption d'activité, la perte de données, les coûts d'acquisition de services de substitution, les réclamations de tiers, la perte de revenus attendus ou les dommages punitifs, sauf disposition contraire de la loi impérative.\n" +
				"En ce qui concerne les jetons gratuits, d'essai, promotionnels et d'activité, la plateforme les fournit « en l'état » et ne promet aucune disponibilité, stabilité, continuité ou responsabilité d'indemnisation, sauf disposition contraire de la loi impérative.\n" +
				"Cet article ne limite pas la responsabilité découlant d'une faute intentionnelle de la plateforme, d'une négligence grave, d'une fraude, de dommages corporels ou d'autres situations où la responsabilité ne peut être limitée par la loi.",

			"Article 15 Force majeure\n" +
				"Si une partie est incapable d'exécuter ou retarde l'exécution du présent contrat en raison de catastrophes naturelles, de guerres, d'émeutes, d'attaques terroristes, d'actes gouvernementaux, de modifications de lois et réglementations, d'exigences réglementaires, de cyberattaques, d'incidents de piratage, de virus, de pannes des télécommunications de base, de pannes de services cloud, de pannes d'énergie, de pannes de canaux de paiement, d'interruptions de services de modèles tiers, de sanctions internationales, de contrôles des exportations, d'épidémies ou d'autres événements de force majeure imprévisibles, inévitables et insurmontables, cette partie sera partiellement ou totalement exonérée de sa responsabilité dans la limite de l'impact de ces événements.\n" +
				"La partie touchée par la force majeure doit en informer l'autre partie dans un délai raisonnable et faire des efforts raisonnables pour atténuer les pertes.",

			"Article 16 Modification du contrat\n" +
				"La plateforme peut mettre à jour le présent contrat conformément aux lois et réglementations, aux exigences réglementaires, aux ajustements de l'activité ou aux modifications du service.\n" +
				"Le contrat mis à jour sera présenté à l'utilisateur par le biais d'annonces sur le site, de pages de paiement, de notifications internes, d'e-mails ou d'autres méthodes raisonnables.\n" +
				"Si l'utilisateur continue d'acheter, de recharger ou d'utiliser des jetons d'IA après la mise à jour du contrat, il est réputé avoir accepté le contrat mis à jour.\n" +
				"Si la modification du contrat a un impact négatif majeur sur les droits et intérêts fondamentaux existants de l'utilisateur concernant les jetons achetés mais non consommés (réduction substantielle du nombre de jetons, réduction significative de la période de validité ou interdiction totale d'utiliser des modèles initialement promis), la plateforme fournira un préavis raisonnable et pourra choisir de proposer une solution alternative de valeur équivalente ou de rembourser au prorata les frais correspondant aux jetons inutilisés, sauf disposition contraire de la loi impérative.",

			"Article 17 Droit applicable et règlement des différends\n" +
				"Utilisateurs nationaux : La formation, la validité, l'exécution, l'interprétation et le règlement des différends du présent contrat sont régis par les lois de la République populaire de Chine (à l'exclusion des règles de conflit de lois) ; tout différend découlant de ou lié au présent contrat doit d'abord être résolu par des négociations amicales ; en cas d'échec des négociations, chaque partie doit soumettre le différend au tribunal populaire compétent du lieu du prestataire de services.\n" +
				"Utilisateurs internationaux : Le présent contrat est régi par le droit du lieu du prestataire de services (tout en respectant les lois impératives du lieu de résidence de l'utilisateur) ; tout différend découlant de ou lié au présent contrat doit d'abord être résolu par des négociations amicales ; en cas d'échec des négociations, il sera soumis à la Commission d'arbitrage économique et commercial international de Chine (CIETAC) pour arbitrage conformément à ses règles d'arbitrage en vigueur au moment de la demande d'arbitrage ; le lieu de l'arbitrage est Shanghai, Chine, la langue de l'arbitrage est le chinois / anglais, et la sentence arbitrale est définitive et obligatoire pour les deux parties.\n" +
				"Si le lieu de résidence de l'utilisateur consommateur contient des dispositions impératives sur la juridiction, ces dispositions prévaudront.",

			"Article 18 Notifications\n" +
				"La plateforme peut envoyer des notifications à l'utilisateur par le biais d'annonces sur le site, de messages internes, d'e-mails, de SMS, d'appels téléphoniques, de notifications sur la console API, de relevés de facturation ou via les coordonnées fournies par l'utilisateur.\n" +
				"L'utilisateur doit s'assurer que ses coordonnées sont réelles, exactes et valides ; en cas de modification de ses coordonnées, il doit les mettre à jour en temps utile.\n" +
				"Règles de notification : Réputée reçue au moment de la publication de l'annonce sur le site ; réputée reçue lors de l'envoi réussi du message interne, de l'e-mail ou du SMS ; réputée reçue lors de la signature pour réception en cas d'envoi par courrier express ou recommandé, ou considérée comme reçue en cas de refus de réception ou d'impossibilité de livraison pour des raisons imputables à l'utilisateur.",

			"Article 19 Divers\n" +
				"Le présent contrat constitue l'intégralité de l'accord entre les parties concernant l'achat et l'utilisation des jetons d'IA (AI Token), et remplace tout accord antérieur, oral ou écrit, conclu entre elles.\n" +
				"Si une disposition du présent contrat est jugée nulle, illégale ou inapplicable, cela n'affectera pas la validité des autres dispositions.\n" +
				"Sans le consentement écrit de la plateforme, l'utilisateur ne peut transférer ses droits et obligations au titre du présent contrat.\n" +
				"La plateforme peut transférer ses droits et obligations au titre du présent contrat en cas de réorganisation de l'activité, de fusion, de scission, de transfert d'actifs, de reprise par une entreprise affiliée ou de migration de services, à condition de s'assurer que les droits existants des utilisateurs ne soient pas substantiellement réduits.\n" +
				"Les titres du présent contrat ne sont fournis que pour en faciliter la lecture et n'affectent en rien l'interprétation de ses dispositions.\n" +
				"Ce contrat peut être rédigé en chinois, en anglais ou dans d'autres langues ; en cas de conflit entre les différentes versions linguistiques, la version chinoise prévaudra, sauf disposition contraire de la loi impérative du lieu de résidence de l'utilisateur.\n" +
				"L'utilisateur doit télécharger, enregistrer ou imprimer le présent contrat de sa propre initiative avant de payer ; la plateforme peut également proposer des options de consultation des versions du contrat, de l'historique des commandes ou des transactions dans l'espace client.",
		},
	},
	"ja": {
		Title: "AI トークン購入および利用規約",
		MetaKeys: map[string]string{
			"version":     "バージョン：V1.1",
			"effective":   "発効日：{{date}}",
			"provider":    "サービス提供者（プラットフォーム）：OSS Energietechnik GmbH",
			"address":     "登録住所：Adam-Opel-Straβe 16-18, 60386 Frankfurt am Main",
			"email":       "連絡先メールアドレス：info@oss-energietechnik.de",
			"user":        "ユーザー：{{user}}",
			"order":       "支払注文：{{order}} ({{amount}})",
			"signMethod":  "署名方法：ユーザーが本規約に同意し、「規約に同意して支払う」「購入を確認する」などの同等のボタンをクリックすることにより、電子署名が有効に成立したものとみなされ、本規約は即時に発効します。",
			"applicable":  "適用地域：国内（中国本土）および国外のユーザーに適用されます。国内ユーザーには中華人民共和国の法律が優先して適用され、国外ユーザーは同時に、その所在地の消費者保護、データ保護、税務、電子商取引、輸出管理、およびデジタルサービスに関する強制法を遵守するものとします。",
			"declaration": "重要な声明：本規約は商業契約書であり、法律上の意見ではありません。正式に導入する前に、ユーザーの所在地法、支払いチャネル、税務、およびデータコンプライアンスの要件に基づいて最終確認を行うことをお勧めします。",
		},
		Sections: []string{
			"第一条 定義\n" +
				"AI トークン：ユーザーがプラットフォームから購入する、プラットフォームの人工知能モデル、アルゴリズム機能、API サービス、エージェントサービス、テキスト生成、画像生成、音声処理、データ分析、またはその他の AI 機能を呼び出すためのサービス利用枠／消費証明を指します。通貨、暗号資産（仮想通貨）、電子マネー、証券、金融商品、プリペイドカード、または譲渡可能な資産には該当しません。\n" +
				"サービス：プラットフォームがユーザーに提供する AI モデルの呼び出し、API 接続、オンライン生成、データ処理、アカウント管理、使用量統計、技術サポートなどの関連サービスを指します。\n" +
				"アカウント：ユーザーがプラットフォームに登録、ログイン、またはプラットフォームによって割り当てられたアカウント、API キー、秘密鍵、組織 ID、プロジェクト ID、またはその他の識別情報を指します。\n" +
				"入力内容：ユーザーがサービスを利用する際に送信、アップロード、入力、送信、または呼び出すテキスト、画像、音声、動画、コード、ファイル、データ、プロンプト、インターフェースパラメータなどのコンテンツを指します。\n" +
				"输出内容：サービスがユーザーの入力内容に基づいて生成または返却するテキスト、画像、音声、動画、コード、分析結果、モデルの回答、またはその他の結果を指します。\n" +
				"消耗：ユーザーが AI サービスを呼び出し、リクエストを送信し、コンテンツを生成し、データを処理し、または関連機能を使用することによって AI トークンを差し引く行為を指します。\n" +
				"合规定性：本 AI トークンは、プラットフォーム専用の AI サービス利用枠であり、流通性や法的支払効力はなく、法定通貨、暗号資産（仮想通貨）との交換や取引目的の投機行為は禁止されています。国内では中央銀行などの監督管理規定が適用され、国外ではユーザーの所在地における暗号資産およびデジタル資産に関する規制規則を遵守するものとします。",

			"第二条 規約の締結と電子的確認\n" +
				"本規約は、電子データの形式で表示、確認、署名、保存され、書面による合意と同等の法的効力を有します。\n" +
				"ユーザーが本規約に同意することを選択し、支払いの実行、購入の確認、注文の送信、支払いの完了、または実際に AI トークンを使用する行為は、本規約の有効な確認および締結を構成します。\n" +
				"プラットフォームは、ユーザーが本規約を確認した際のアカウント情報、注文番号、支払記録、IP アドレス、デバイス情報、ブラウザ情報、操作日時、規約バージョン番号、同意確認のチェック記録、クリック記録、ログ記録などを記録および保存する権利を有し、これらは規約の成立、履行、紛争解決、およびコンプライアンス監査の証明として使用されます。\n" +
				"ユーザーは、書面による契約が締結されていないこと、物理的な押印がないこと、または書面での署名がないことを理由に、本規約の効力を否定することはできません。\n" +
				"ユーザーが企業、機関、またはその他の組織である場合、その主体を代表してサービスを購入または利用する担当者は、十分な権限を取得していることを確認します。正当な授権なしに他者を代表してサービスを購入または利用した場合、実際の操作者がそれに応じた責任を負うものとします。\n" +
				"本規約における返金不可や責任の制限など、プラットフォームの責任を免除または制限する条項は、本文中での太字表示や支払ページ上での独立したポップアップ等により明示的に提示されており、ユーザーはこれらを十分に理解し、自発的に受け入れたことを確認します。",

			"第三条 トークンの購入、価格設定および引渡し\n" +
				"ユーザーは、プラットフォームのページに表示されるパッケージ、価格、通貨、数量、有効期限、適用モデル、消費ルール、およびその他の説明に従って、AI トークンを購入できます。\n" +
				"AI トークンの価格、交換比率、消費ルール、対応モデル、コンテキスト長、同時接続数制限、および機能範囲は、ユーザーが注文した時点でプラットフォームに表示されている情報が適用されます。プラットフォームが変更を公表した場合は、公表された内容に従って実行されます。\n" +
				"プラットフォームが明確に別段の指定をしない限り、AI トークンはプラットフォーム内の専用サービスでの利用に限定され、金融資産には該当せず、他の資産への譲渡、取引、または交換を行うことはできません。\n" +
				"プラットフォームは、ユーザーからの支払金を受領し、注文を確認した後、24時間以内に対応する AI トークンをユーザーのアカウントに付与（チャージ）します。期限内にアカウントに反映されない場合、ユーザーは無条件で全額返金を申請する権利を有します。\n" +
				"支払チャネル、銀行、サードパーティの決済機関、外国交換決済、不正防止審査、税務審査、またはコンプライアンス審査に起因するチャージの遅延、手数料の差し引きなどについて、プラットフォームはサードパーティの理由による責任を負いませんが、合理的な範囲でユーザーの問い合わせを支援するものとします。\n" +
				"ユーザーは、注文情報、アカウント情報、請求書（領収書）情報、および支払主体の情報が真実、正確、かつ完全であることを保証するものとします。ユーザー情報の誤りによりチャージの失敗、請求書の誤り、アカウント所有権の紛争、またはその他の損失が発生した場合、ユーザーが自らその責任を負うものとします。",

			"第四条 トークンの利用ルール\n" +
				"ユーザーが購入した AI トークンは、プラットフォームが指定する範囲内でのみ使用できます。具体的にどのモデル、サービス、API、機能、または製品に使用できるかは、プラットフォームのページ、管理コンソール、または注文の記載に従うものとします。\n" +
				"AI トークンは、プラットフォームが公表する測定方法に従って消費されます。これには、入力文字数、出力文字数、トークン数、画像枚数、音声時間、動画時間、リクエスト回数、モデルタイプ、計算リソース、ストレージリソース、プラグイン呼び出し、ツール呼び出しなどが含まれますが、これらに限定されません。\n" +
				"モデル、機能、リージョンノード、およびサービスレベルによってトークンの消費基準が異なる場合があり、ユーザーは使用前に対応する料金説明を確認する必要があります。\n" +
				"AI トークンが消費された時点で、サービスは引渡しが完了した（または部分的に完了した）ものとみなされます。本規約に別段の定めがある場合、または法律による強制規定がある場合を除き、消費された AI トークンの復元、返金、譲渡、または現金化はサポートされません。\n" +
				"ユーザーは、アカウント、パスワード、API キー、アクセスキー、確認コード、およびその他の個人認証情報を適切に管理する必要があります。ユーザーのアカウントまたは秘密鍵を通じて開始されたすべての呼び出しは、ユーザー本人またはその授権者による行為とみなされます。\n" +
				"ユーザーは、アカウントまたは秘密鍵の漏洩、盗用、または異常な呼び出しを発見した場合、直ちにプラットフォームに通知し、パスワードのリセット、キーの無効化、インターフェースの閉鎖などの措置を講じる必要があります。プラットフォームが通知を受理する前に発生したトークンの消費は、プラットフォームに故意または重大な過失がある場合を除き、原則としてユーザーの負担となります。\n" +
				"ユーザーは、プラットフォームの外で AI トークンを譲渡、販売、賃貸、貸与、出金、現金化、または取引してはならず、プラットフォーム外での取引や資金調達活動に利用してはなりません。",

			"第五条 有効期限、失効および更新\n" +
				"AI トークンの有効期限は、プラットフォームが購入ページ、注文ページ、パッケージ説明、またはユーザーアカウントのダッシュボードに明示的に表示する情報に従うものとします。明示的に表示されていない場合、デフォルトの有効期限はチャージ完了日から12ヶ月間です。\n" +
				"有効期限が切れた場合、未使用の AI トークンは自動的に失効し、プラットフォームは利用、返金、交換、または期限延長のサービスを提供しません。ただし、法律の強制規定により別段の要求がある場合はこの限りではありません。\n" +
				"プラットフォームは、運営の都合により、期限延長、チャージ更新、アップグレード、およびパッケージの移行サービスを提供する場合があります。具体的なルールは、その時点でプラットフォームに表示されている内容に従うものとします。\n" +
				"ユーザーが新しいパッケージを購入した場合、新旧トークンの消費順序はプラットフォームのシステム規則に従うものとします。特段の説明がない限り、プラットフォームは有効期限が近いトークンを優先的に消費することができます。",

			"第六条 返金と解約（国内外共通＋法域適応）\n" +
				"コアルール：AI トークンは即時に提供され、即時に消費されるデジタルサービスであり、8日間の無条件返品ルール（クーリングオフなど）は適用されません。支払完了後は、原則として無条件の返金は行われません。ただし、本条第2項、第3項、第4項、および法律の強制規定による場合は除きます。\n" +
				"国内ユーザー向けの返金ルール：(1) チャージ完了後は、返金ルールは適用されず、返金は行われません。(2) すでに消費されたトークン、無償トークン、キャンペーンで付与されたトークン、体験トークン、クーポン適用分、本規約違反により利用制限されたトークン、およびユーザー自身の誤操作（誤チャージ、誤呼び出しなど）により消費されたトークンは、返金、出金、譲渡、または現金化の対象外となります。\n" +
				"国外ユーザー向けの返金ルール（EU、英国などの法域を含む）：(1) ユーザーの所在地において法定の契約撤回権（キャンセル権）が認められている場合、ユーザーが「即時履行を希望し、チャージ完了と同時に撤回権を喪失すること」を確認・同意してチェックを入れた時点で、その喪失が有効になります。(2) 即時履行の同意がなく、かつトークンが未消費である場合は、ユーザーの所在地の法定撤回期間が適用されます。消費された部分については返金できません。\n" +
				"プラットフォームの過失による返金：プラットフォームの責めに帰すべき事由により、ユーザーが購入した未期限切れトークンに対応するコアサービスを連続して30暦日間利用できなかった場合、プラットフォームはユーザーから有効な通知を受理した後15営業日以内に、同等価値のトークンの補填、有効期限の延長、または未使用部分に相当する返金措置を提供するものとします。期限内に是正されない場合、ユーザーは未使用トークンの割合に応じた返金を請求する権利を有します。",

			"第七条 ユーザーの法的義務と遵守事項\n" +
				"ユーザーは、中華人民共和国の法律、およびユーザーの所在地の法律、規制要件、業界の規範、輸出管理、経済制裁、データの国境を越えた転送規則、ならびに公序良俗を遵守することを誓約するものとします。\n" +
				"ユーザーは、違法、他者の権利を侵害する、詐欺的、虚危、憎悪表現、ハラスメント、暴力、アダルト、テロリズム、極端主義、自傷行為、自殺、薬物、武器、悪意のあるコード、サイバー攻撃、違法な金融活動など、規制に違反するコンテンツの生成や伝播を行ってはならず、その生成を支援してはなりません。\n" +
				"ユーザーは、他者の知的財産権、営業秘密、プライバシー権、個人情報に関する権利、肖像権、名誉権、またはその他の合法的な権利や利益を侵害してはなりません。\n" +
				"ユーザーは、プラットフォームの技術的制限、安全管理メカニズム、アクセス制御、または課金システムを、許可なくクロール、複製、学習、リバースエンジニアリング、ハッキング、または回避してはなりません。\n" +
				"ユーザーは、自動化されたスパム、アクセス数の偽装（刷量）、虚偽のレビュー、フィッシング、詐欺、他者へのなりすまし、コンテンツ審査の回避、または違法コンテンツの大量生成のためにサービスを利用してはなりません。\n" +
				"ユーザーは、自身が処理する正当な権限を持たないデータ、個人情報、機微な個人情報、秘密情報、国家秘密、営業秘密、または輸出管理制限対象のデータをプラットフォームに送信してはなりません。\n" +
				"プラットフォームは、法律規制、監督管理要求、プラットフォームの規則、リスク管理ポリシー、またはサードパーティのモデル提供元からの要請に従い、ユーザーの入力内容、出力内容、呼び出し行為、およびアカウント行為に対して、必要な安全確認、リスク認識、およびコンプライアンス処理を行う権利を有します。",

			"第八条 AI 出力内容とリスクに関する免責提示\n" +
				"ユーザーは、AI の出力結果が入力内容、モデルパラメータ、学習データ、コンテキスト、およびシステム設定に基づきアルゴリズムモデルによって生成されるものであり、不正確、不完全、古い内容、創作、偏見、適用不可能、または説明がつかない内容が含まれる可能性があることを理解し、同意するものとします。\n" +
				"出力結果は、法律、財務、医療、投資、設計、安全、コンプライアンス、またはその他の専門的なアドバイスを構成するものではありません。ユーザーが出力結果に基づいて重要な決定を下す必要がある場合は、自ら人的な検証を行い、資格を持つ専門家に相談するものとします。\n" +
				"ユーザーは、自身の入力内容、サービスの使用方法、出力結果の検証、公開、伝播、応用、およびそれに起因する結果について、一切の責任を負うものとします。\n" +
				"プラットフォームは、AI 出力結果が完全に正確であること、サービスが継続的に利用可能であること、ユーザーの特定の目的に合致すること、また出力結果が他のユーザーの生成結果と同一または類似しないことを保証しません。\n" +
				"ユーザーが出力結果を商業目的、一般公開、広告宣伝、入札書類、法的書類、医療アドバイス、投資アドバイス、設計工学、自動意思決定、またはその他の高リスクのシナリオで使用する場合、ユーザー自身がコンプライアンスの確認および使用リスクを負うものとします。",

			"第九条 知的財産権\n" +
				"プラットフォーム、その関連会社、およびライセンサーは、プラットフォームシステム、ソフトウェア、モデル、アルゴリズム、インターフェース、ドキュメント、ページデザイン、商標、ロゴ、データベース、技術的スキーム、ビジネスモデル、およびそれらに関連する知的財産権に対して、完全な権利を有します。\n" +
				"ユーザーによる AI トークンの購入は、プラットフォームのモデル、アルゴリズム、ソフトウェア、ソースコード、システムアーキテクチャ、商標、またはその他の知的財産権の所有権をユーザーに付与するものではありません。\n" +
				"ユーザーは、合法的に提供した入力内容に対して、自らが法的に享受する権利を保持します。\n" +
				"ユーザーが本規約を遵守し、対応する料金を支払っていることを前提として、プラットフォームは自らの権利の範囲内で、ユーザーに対して本サービスによって生成された出力内容の利用を許諾します。ただし、出力結果が排他的であるか、第三者の権利を侵害していないか、知的財産権として登録可能か、商業利用可能かについては、ユーザーが特定のコンテンツや適用される法律に基づいて自己判断するものとし、侵害のリスクはユーザー自身が負うものとします。\n" +
				"ユーザーは、サービスの提供、料金計算、リスク管理、安全確認、障害対応、サービスの最適化、コンプライアンス監査、および法的義務の履行に必要な範囲内で、プラットフォームが入力内容、出力内容、および利用データを処理することを許諾します。\n" +
				"プラットフォームがサードパーティのベースモデル、クラウドサービス、決済サービス、またはその他のサードパーティのサービスを利用する場合、対応する権利制限、コンテンツポリシー、および利用制限がユーザーにも同時に適用される場合があります。",

			"第十条 個人情報保護とプライバシー（国内外のコンプライアンス対応）\n" +
				"国内ユーザー：『中華人民共和国個人情報保護法』が適用され、必要最小限の原則に従います。プラットフォームは、サービスの提供に必要な情報のみを収集します：①登録情報（携帯電話番号、メールアドレス）、②支払情報（Alipay注文番号、支払金額）、③デバイス情報（IPアドレス、ブラウザタイプ）、④使用状況情報（トークンのチャージ・消費履歴）。身分証明書番号やクレジットカード番号などの機微な個人情報は収集しません。\n" +
				"国外ユーザー：ユーザーの所在地における個人情報保護法（GDPRなど）が適用されます。プラットフォームは、ユーザーの所在地の法律に従って個人情報を処理し、十分な個人情報保護水準に達していない地域への国境を越えたデータ転送を行いません。必要に応じて、ユーザーとデータ処理合意（DPA）を締結することができます。\n" +
				"個人情報の保存期間：アカウントの削除またはサービスの終了日から1年間保存され、期限終了後は自動的に匿名化または削除されます。注文履歴、税務記録、トークンのチャージ・消費履歴は3年間保存されます（税務およびコンプライアンス監査の要件に対応するため）。\n" +
				"ユーザーが第三者に関する個人情報、機微な個人情報、営業秘密、秘密情報、または規制データなどをプラットフォームに入力する場合、ユーザーは合法的な処理基盤を有していること、および必要な告知、同意、授権、評価、または届出義務を履行していることを保証するものとします。\n" +
				"プラットフォームは、必要かつ合理的な範囲内で、データのクラウド保管サービス提供元、モデル提供元、決済代行会社、SMS送信会社、インボイス発行サービス、カスタマーサポートシステム、データセキュリティ関連企業などのコンプライアンスを遵守した第三者に対しデータ処理を委託することができ、委託先に対して合理的なデータ保護管理措置を講じることを求めます。",

			"第十一条 税金、請求書およびクロスボーダー決済\n" +
				"ユーザーは、注文ページに表示された金額に従って料金を支払うものとします。関連する料金に税金が含まれているかどうかは、プラットフォームのページ、注文、または請求書の説明に従うものとします。\n" +
				"ユーザーの所在地、支払方法、決済チャネル、加盟店契約会社、銀行、為替決済、消費税、付加価値税、売上税、源泉徴収税、デジタルサービス税、またはその他の税金や手数料により発生する追加費用は、適用される法律および注文の合意に従って双方が負担するものとします。外貨による支払の場合、為替レート、銀行手数料、決済チャネル費用、または決済時間の違いにより、実際の引き落とし額が異なる場合があります。プラットフォームが別途約束した場合を除き、このような差額はユーザー自身の負担とします。\n" +
				"ユーザーが請求書、領収書、または税務証明書を必要とする場合、プラットフォームの要請に従って、真実、正確、かつ完全な宛先情報を提出する必要があります。ユーザーが提出した情報の誤り、または支払主体と宛先主体の不一致により、請求書の発行ができなかったり、修正が必要になったりした場合、ユーザーは処理に協力し、それに伴う損失を負担するものとします。",

			"第十二条 サービスの変更、中断および保守\n" +
				"プラットフォームは、ビジネスの発展、技術のアップグレード、モデルの提供状況、コストの変動、法律規制、監督管理要求、またはサードパーティサービスの変更に基づき、サービス内容、モデルタイプ、料金規則、トークン消費比率、使用制限、および機能範囲を調整することができます。ただし、前述の変更により、ユーザーが購入した未使用トークンの利用範囲が著しく縮小されたり、単位サービスコストが著しく上昇（トークン消費比率が20%以上増加）したりした場合、ユーザーは変更の発効前に、未使用トークンに相当する料金の比例返金を申請する権利を有します。前述の変更は、すでに購入され未消費のトークンが元のサービス範囲内で有する利用権利には影響しません。\n" +
				"プラットフォームが重大な変更を行う場合、ウェブサイトのアナウンス、サイト内メッセージ、メール、SMS、またはその他の合理的な方法によりユーザーに通知するものとします。\n" +
				"システムの保守、アップグレード、ネットワーク障害、クラウドサービスの障害、サードパーティのモデルサービスの障害、決済チャネルの障害、セキュリティインシデント、不可抗力、または監督管理要求により、サービスの中断、遅延、または利用不能が発生した場合、プラットフォームはサービスの復旧に合理的な努力を払いますが、間接的な損害については責任を負いません。\n" +
				"プラットフォームは、ユーザー側のネットワーク、デバイス、インターフェース設定、キーの管理、コードエラー、呼び出しパラメータの誤り、利用ルールの違反、またはサードパーティの事由に起因する利用失敗について、一切の責任を負いません。",

			"第十三条 アカウントの制限と違約処理\n" +
				"ユーザーが本規約または適用される法律規制に違反した場合、プラットフォームは状況に応じて、是正勧告、違法コンテンツの削除要求、関連する呼び出し機能の一時停止などの措置を講じる権利を有します。\n" +
				"プラットフォームは、一部の機能の制限、同時呼び出し数の制限、クォータの削減、API キーの一時停止、アカウントの凍結、サービスの休止、または提供の終了を行うことができます。\n" +
				"プラットフォームは、規制違反により取得または使用されたトークンを差し引き、凍結、または無効化することができます。違反により生じた注文について、プラットフォームは返金を拒否し、優待をキャンセルし、損失を回収する権利を有します。\n" +
				"法律で認められる範囲内で、プラットフォームは監督官庁、司法機関、権利者、または被害を受けた当事者に対して必要な情報を提供し、法律で認められるその他の措置を講じることができます。\n" +
				"ユーザーの違約行為により、プラットフォーム、その関連会社、パートナー、他のユーザー、または第三者が損失、損害賠償請求、罰則、調査、訴訟、弁護士費用、公証費用、鑑定費用、旅費交通費などを被った場合、ユーザーは一切の賠償責任を負うものとします。",

			"第十四条 責任の制限\n" +
				"法律で認められる最大の範囲内で、サービスに関するプラットフォームの賠償責任は、紛争の生じた注文に対してユーザーが実際に支払った、未消費かつ未返金の金額を上限とします。\n" +
				"プラットフォームは、間接的損害、逸失利益、営業信用の喪失、業務の中断、データの紛失、代替サービスの調達コスト、第三者からの損害賠償請求、期待収益の損失、または懲罰的損害賠償について責任を負いません。ただし、法律の強制規定により別段の要求がある場合はこの限りではありません。\n" +
				"無償で提供されたトークン、体験トークン、ボーナストークン、およびキャンペーンで付与されたトークンについて、プラットフォームはそれらを「現状有姿」で提供するものとし、その可用性、安定性、継続性、または賠償責任を約束しません。ただし、法律の強制規定により別段の要求がある場合はこの限りではありません。\n" +
				"本条は、プラットフォームの故意、重大な過失、詐欺、人身傷害、または法律により責任を制限できない状況に起因する責任を制限するものではありません。",

			"第十五条 不可抗力\n" +
				"自然災害、戦争、暴動、テロ攻撃、政府行為、法律規制の変更、監督管理要求、サイバー攻撃、ハッカー事件、ウイルス、基礎電気通信障害、クラウドサービス障害、電力供給の中断、決済チャネルの障害、サードパーティのモデルサービスの中断、国際制裁、輸出管理、疫病の流行、またはその他の予見不可能、回避不可能、かつ克服不可能な不可抗力イベントにより、一方が本規約の履行を怠り、または遅延した場合、その当事者は影響を受けた範囲内で責任を一部または全部免除されるものとします。\n" +
				"不可抗力の影響を受けた当事者は、合理的な期間内に他方に通知し、損失を最小限に抑えるために合理的な努力を払うものとします。",

			"第十六条 規約の変更\n" +
				"プラットフォームは、法律規制、監督管理要求、ビジネスの調整、またはサービスの変更に伴い、本規約を更新することができます。\n" +
				"更新された規約は、ウェブサイトのアナウンス、決済ページ、サイト内通知、メール、またはその他の合理的な方法によってユーザーに提示されます。\n" +
				"規約が更新された後、ユーザーが引き続き AI トークンを購入、チャージ、または利用した場合、ユーザーは更新後の規約に同意したものとみなされます。\n" +
				"規約の変更が、ユーザーがすでに購入した未消費トークンの主要な権利や利益に重大な悪影響を与える場合（トークン数量の実質的な削減、有効期限の著しい短縮、または当初約束されたモデルの利用の完全な禁止）、プラットフォームは合理的な事前通知を行い、同等価値の代替手段を提供するか、未使用トークンの割合に応じた比例返金を提供するものとします。ただし、法律の強制規定により別段の要求がある場合はこの限りではありません。",

			"第十七条 準拠法および紛争解決（国内外別設定）\n" +
				"国内ユーザー：本規約の成立、効力、履行、解釈および紛争解決には、中華人民共和国の法律（抵触法規則を除く）が適用されます。本規約に起因または関連する紛争について、双方はまず友好関係に基づき解決を図るものとします。解決できない場合、いずれの当事者も、サービス提供者の所在地を管轄する権利を有する人民法院に訴訟を提起するものとします。\n" +
				"国外ユーザー：本規約は、サービス提供者の所在地の法律に準拠します（同時にユーザー所在地の強制規則を遵守します）。本規約に起因または関連する紛争について、双方はまず友好関係に基づき解決を図るものとします。解決できない場合、中国国際経済貿易仲裁委員会（CIETAC）に仲裁を申請し、申請時点で有効な同委員会の仲裁規則に従って仲裁解決を図るものとします。仲裁地は中国・上海とし、仲裁言語は中国語／英語とします。仲裁判断は最終的なものであり、双方を拘束するものとします。\n" +
				"消費者ユーザーの所在地に管轄に関する強制的な規定がある場合は、その規定に従うものとします。",

			"第十八条 通知\n" +
				"プラットフォームは、ウェブサイトのアナウンス、サイト内メッセージ、電子メール、SMS、電話、API コンソールの通知、請求書の通知、またはユーザーが登録した連絡先を通じて通知を送信することができます。\n" +
				"ユーザーは、自身の連絡先情報が真実、正確、かつ有効であることを保証するものとし、連絡先が変更された場合は、速やかに更新しなければなりません。\n" +
				"通知の到達ルール：ウェブサイトのアナウンスは公表された時点で到達したものとみなします。サイト内メッセージ、電子メール、SMS は送信に成功した時点で到達したものとみなします。宅配便または書留郵便は受領署名がされた時点、またはユーザーの事由により受領拒否や配達不能となった時点で到達したものとみなします。",

			"第十九条 その他\n" +
				"本規約は、AI トークンの購入および使用事項に関する双方の完全な合意を構成し、これより前に口頭または書面でなされた合意に優先します。\n" +
				"本規約のいずれかの条項が無効、違法、または執行不能と判断された場合でも、その他の条項の効力には影響しません。\n" +
				"ユーザーは、プラットフォームの書面による事前の同意なしに、本規約に基づく権利または義務を譲渡することはできません。\n" +
				"プラットフォームは、事業再編、合併、分割、資産譲渡、関連会社への承継、またはサービスの移行に伴い、本規約に基づく権利および義務を譲渡することができます。ただし、ユーザーの既存の権利や利益が実質的に減少されないことを保証するものとします。\n" +
				"本規約の各条項の見出しは、単に読みやすさのために付されたものであり、条項の解釈に影響を与えるものではありません。\n" +
				"本規約は、中国語、英語、またはその他の言語で提供される場合があります。多言語版の間で齟齬や矛盾が生じた場合は、ユーザーの所在地の強制法による場合を除き、中国語版が優先して適用されます。\n" +
				"ユーザーは支払を行う前に、自ら本規約をダウンロード、保存、または印刷する必要があります。プラットフォームは、ユーザーのアカウント内で規約のバージョン、注文履歴、または取引履歴の照会機能を提供する場合があります。",
		},
	},
	"ru": {
		Title: "Соглашение о покупке и использовании AI токенов",
		MetaKeys: map[string]string{
			"version":     "Версия: V1.1",
			"effective":   "Дата вступления в силу: {{date}}",
			"provider":    "Поставщик услуг (Платформа): OSS Energietechnik GmbH",
			"address":     "Юридический адрес: Adam-Opel-Straβe 16-18, 60386 Frankfurt am Main",
			"email":       "Контактный email: info@oss-energietechnik.de",
			"user":        "Пользователь: {{user}}",
			"order":       "Платежный заказ: {{order}} ({{amount}})",
			"signMethod":  "Способ подписания: Отметка пользователем согласия с настоящим соглашением и нажатие кнопок типа «Согласиться с соглашением и оплатить», «Подтвердить покупку» рассматривается как подписание соглашения в электронной форме, и соглашение вступает в силу немедленно.",
			"applicable":  "Применимый регион: Применимо как к внутренним (материковый Китай), так и к зарубежным пользователям. К внутренним пользователям применяется право Китайской Народной Республики; зарубежные пользователи также должны соблюдать обязательное законодательство своей страны о защите прав потребителей, защите данных, налогообложении, электронной коммерции, экспортном контроле и цифровых услугах.",
			"declaration": "Важное заявление: Настоящее соглашение является коммерческим контрактом, а не юридическим заключением; перед официальным запуском рекомендуется выполнить окончательную проверку в соответствии с юрисдикцией пользователя, каналами оплаты, налоговыми требованиями и требованиями к комплаенсу данных.",
		},
		Sections: []string{
			"Статья 1 Определения\n" +
				"AI Token: Означает приобретенную пользователем у платформы квоту на вызовы услуг / расходный ваучер для вызова моделей искусственного интеллекта платформы, возможностей алгоритмов, услуг API, услуг агентов, генерации текста, генерации изображений, обработки голоса, анализа данных или других функций ИИ. Не является валютой, виртуальной валютой, электронными деньгами, ценными бумагами, финансовыми продуктами, картами предоплаты или передаваемыми активами.\n" +
				"Услуги: Означает вызов моделей ИИ, интеграцию API, онлайн-генерацию, обработку данных, управление учетными записями, статистику использования, техническую поддержку и другие сопутствующие услуги, предоставляемые платформой пользователям.\n" +
				"Учетная запись: Означает учетную запись, API-ключ, секретный ключ, идентификатор организации, идентификатор проекта или другие идентификаторы, зарегистрированные, используемые для входа или выделенные платформой пользователю.\n" +
				"Входные данные: Означает текст, изображения, аудио, видео, код, файлы, данные, промпты, параметры интерфейса и другие материалы, отправленные, загруженные, введенные, переданные или вызванные пользователем при использовании услуг.\n" +
				"Выходные данные: Означает текст, изображения, аудио, видео, код, результаты анализа, ответы моделей или другие результаты, сгенерированные или возвращенные услугой на основе входных данных пользователя.\n" +
				"Расход: Означает списание AI токенов в результате вызова пользователем услуг ИИ, отправки запросов, генерации контента, обработки данных или использования связанных функций.\n" +
				"Соответствие требованиям: Настоящие AI токены являются исключительно собственной квотой платформы для вызова услуг ИИ, не имеют обращения или статуса законного платежного средства, не могут быть обменяны на фиатную валюту, виртуальную валюту или использованы для спекуляций; внутри страны применяются правила центрального банка и других регулирующих органов, за рубежом — правила регулирования виртуальных валют и цифровых активов по месту нахождения пользователя.",

			"Статья 2 Подписание соглашения и электронное подтверждение\n" +
				"Настоящее соглашение отображается, подтверждается, подписывается и сохраняется в виде электронных данных и имеет ту же юридическую силу, что и письменное соглашение.\n" +
				"Действия пользователя по проставлению отметки о согласии с настоящим соглашением и нажатию кнопки оплаты, подтверждению покупки, отправке заказа, завершению платежа или фактическому использованию AI токенов составляют действительное подтверждение и подписание настоящего соглашения.\n" +
				"Платформа имеет право записывать и сохранять информацию об учетной записи пользователя, номер заказа, платежную информацию, IP-адрес, информацию об устройстве, браузере, времени операции, номере версии соглашения, отметках, кликах, логах и т.д. при подтверждении пользователем настоящего соглашения для доказательства заключения, исполнения, разрешения споров и аудита соответствия.\n" +
				"Пользователь не может отрицать действительность настоящего соглашения на том основании, что письменный договор не был подписан, физическая печать не была поставлена или подпись на бумажном носителе отсутствует.\n" +
				"Если пользователем является предприятие, учреждение или иная организация, лицо, совершающее покупку или использующее услуги от имени этого субъекта, подтверждает наличие у него полных полномочий; в случае покупки или использования услуг от имени других лиц без надлежащих полномочий ответственность несет лицо, фактически совершившее эти действия.\n" +
				"Положения настоящего соглашения, освобождающие платформу от ответственности или ограничивающие ее (такие как правила возврата и ограничения ответственности), были выделены жирным шрифтом и продублированы в виде отдельного всплывающего окна на странице оплаты. Пользователь подтверждает, что полностью понял и добровольно принимает их.",

			"Статья 3 Покупка, тарификация и предоставление токенов\n" +
				"Пользователи могут приобретать AI токены в соответствии с пакетами услуг, ценами, валютой, количеством, сроком действия, применимыми моделями, правилами расхода и другими сведениями, представленными на странице платформы.\n" +
				"Цена, курс обмена, правила расхода, поддерживаемые модели, длина контекста, ограничения на параллельные запросы и набор функций AI токенов определяются информацией, представленной на платформе в момент размещения заказа пользователем; в случае публикации платформой изменений они исполняются в соответствии с опубликованным содержанием.\n" +
				"Если платформа явно не указала иное, AI токены могут быть использованы только в рамках внутренних услуг платформы, не являются финансовыми активами, не могут быть переданы, проданы или обменены на другие активы.\n" +
				"После получения оплаты от пользователя и подтверждения заказа платформа начислит соответствующие AI токены на учетную запись пользователя в течение 24 часов; в случае задержки начисления пользователь имеет право безоговорочно потребовать полный возврат средств.\n" +
				"Платформа не несет ответственности за задержки в начислении средств, списание комиссий и т.д., вызванные действиями платежных каналов, банков, сторонних платежных организаций, расчетами по обмену валюты, антифрод-проверками, налоговыми проверками или проверками на соответствие требованиям. Тем не менее платформа в разумных пределах окажет пользователю содействие в выяснении причин.\n" +
				"Пользователь должен гарантировать, что информация о заказе, учетной записи, счете-фактуре и лице, осуществляющем платеж, является достоверной, точной и полной; в случае если ошибки в данных пользователя приведут к сбою при пополнении, ошибкам в счете-фактуре, спорам о принадлежности аккаунта или иным убыткам, ответственность несет сам пользователь.",

			"Статья 4 Правила использования токенов\n" +
				"Приобретенные пользователем AI токены могут быть использованы только в пределах, установленных платформой. Конкретные модели, услуги, API, функции или продукты, для которых они могут применяться, определяются информацией на странице платформы, в консоли управления или в описании заказа.\n" +
				"Списание AI токенов осуществляется согласно опубликованным платформой методам измерения, включая, помимо прочего, количество входных символов, количество выходных символов, количество токенов, количество изображений, длительность аудио, длительность видео, количество запросов, тип модели, вычислительные ресурсы, ресурсы хранения, вызовы плагинов, вызовы инструментов и т. д.\n" +
				"Стандарты списания токенов могут отличаться для разных моделей, функций, региональных узлов и уровней обслуживания. Пользователь должен ознакомиться с соответствующими правилами тарификации перед использованием.\n" +
				"После расходования AI токенов услуга считается оказанной или частично оказанной; если иное не оговорено в настоящем соглашении и не предусмотрено применимым законодательством, израсходованные AI токены не подлежат восстановлению, возврату, передаче или обналичиванию.\n" +
				"Пользователь должен надлежащим образом хранить имя пользователя, пароли, API-ключи, ключи доступа, коды подтверждения и другие учетные данные; все вызовы, инициированные с использованием учетной записи или ключей пользователя, считаются действиями самого пользователя или уполномоченных им лиц.\n" +
				"Если пользователь обнаружит, что учетная запись или ключ были скомпрометированы, похищены или совершаются аномальные вызовы, он должен немедленно уведомить платформу и принять меры, такие как сброс пароля, деактивация ключа, закрытие интерфейса и т. д.; расход токенов, произошедший до получения уведомления платформой, в принципе относится на счет пользователя, за исключением случаев умысла или грубой неосторожности со стороны платформы.\n" +
				"Пользователь не должен передавать, продавать, сдавать в аренду, давать взаймы, выводить, обналичивать или торговать AI токенами вне платформы, а также использовать их для транзакций или финансовой деятельности за пределами платформы.",

			"Статья 5 Срок действия, истечение срока и продление\n" +
				"Срок действия AI токенов определяется информацией, явно указанной платформой на странице покупки, странице заказа, в описании тарифа или в личном кабинете пользователя; если срок действия явно не указан, по умолчанию он составляет 12 месяцев с даты зачисления токенов на счет.\n" +
				"По истечении срока действия неиспользованные AI токены автоматически аннулируются, платформа не предоставляет услуги по их продлению, возврату или обмену, если иное не предусмотрено императивными нормами применимого права.\n" +
				"Платформа может предлагать услуги по продлению срока действия, дозакупке, обновлению тарифов или переходу на другие тарифные планы в соответствии с операционной политикой; конкретные правила определяются информацией, отображаемой на платформе в соответствующее время.\n" +
				"При покупке нового пакета порядок расходования старых и новых токенов определяется системными правилами платформы; если не указано иное, платформа может в приоритетном порядке списывать токены с наиболее близким сроком истечения.",

			"Статья 6 Возврат средств и отказ от услуг\n" +
				"Основные правила: AI токены являются цифровыми услугами с мгновенным предоставлением и потреблением, к ним не применяется право на возврат в течение 7 дней без объяснения причин; после совершения платежа и зачисления средств возврат без объяснения причин по общему правилу не производится, за исключением случаев, предусмотренных пунктами 2, 3 и 4 настоящей Статьи, а также императивными нормами законодательства.\n" +
				"Правила возврата средств для внутренних пользователей: (1) После зачисления токенов правила возврата не применяются, средства не возвращаются; (2) Израсходованные AI токены, бонусные токены, токены по акциям, пробные токены, суммы скидок по купонам, токены, использование которых ограничено в связи с нарушением настоящего соглашения, а также токены, израсходованные из-за ошибок пользователя (ошибочное пополнение, случайные вызовы), не подлежат возврату, выводу, передаче или обналичиванию.\n" +
				"Правила возврата средств для зарубежных пользователей (включая ЕС, Великобританию и др.): (1) Если по месту нахождения пользователя законом предусмотрено право на отзыв соглашения, оно прекращает действие после подтверждения пользователем согласия с тем, что «исполнение начинается немедленно, право на отзыв теряется при зачислении токенов»; (2) При отсутствии подтверждения немедленного исполнения и при условии, что токены не были израсходованы, применяется установленный законом срок для отзыва; израсходованная часть токенов возврату не подлежит.\n" +
				"Возврат средств по вине платформы: Если по вине платформы пользователь не может использовать базовые услуги, соответствующие приобретенным и неиспользованным токенам, в течение 30 календарных дней подряд, платформа обязана в течение 15 рабочих дней после получения уведомления от пользователя предоставить меры компенсации (начисление эквивалентных токенов, продление срока действия или возврат стоимости неиспользованной части); при невыполнении обязательств в срок пользователь имеет право потребовать пропорционального возврата стоимости неиспользованных токенов.",

			"Статья 7 Обязанности пользователя по соблюдению требований\n" +
				"Пользователь обязуется соблюдать законы Китайской Народной Республики, а также законы, нормативные требования, отраслевые стандарты, экспортный контроль, экономические санкции, правила трансграничной передачи данных и общественный порядок по месту своего нахождения.\n" +
				"Пользователю запрещается генерировать, распространять или способствовать генерации контента, нарушающего закон, авторские права, содержащего элементы мошенничества, лжи, ненависти, преследования, насилия, порнографии, терроризма, экстремизма, призывов к самоповреждению, самоубийству, связанного с наркотиками, оружием, вредоносным кодом, кибератаками, незаконной финансовой деятельностью или иного нарушающего правила контента.\n" +
				"Пользователь не должен нарушать права интеллектуальной собственности третьих лиц, коммерческую тайну, права на конфиденциальность, личные данные, право на изображение, репутацию или иные законные права.\n" +
				"Пользователю запрещается без разрешения осуществлять сбор данных (парсинг), копирование, обучение моделей, обратный инжиниринг, взлом или обход технических ограничений, механизмов безопасности, контроля доступа или биллинговых систем платформы.\n" +
				"Пользователь не должен использовать услуги для рассылки спама, искусственного завышения трафика, создания ложных отзывов, фишинга, мошенничества, выдачи себя за другое лицо, обхода модерации или массового создания запрещенного контента.\n" +
				"Пользователь не должен передавать платформе данные, персональную информацию, конфиденциальную информацию, государственные секреты или данные, подпадающие под экспортный контроль, на обработку которых у него нет законных прав.\n" +
				"Платформа имеет право проводить необходимые проверки безопасности, выявлять риски и предпринимать меры по обеспечению соответствия в отношении входных данных пользователя, выходных данных, действий по вызову услуг и поведения учетной записи в соответствии с законами и правилами, требованиями регуляторов, правилами платформы, политикой управления рисками или требованиями сторонних поставщиков моделей.",

			"Статья 8 Выходные данные ИИ и предупреждение о рисках\n" +
				"Пользователь понимает и соглашается с тем, что выходные данные ИИ генерируются алгоритмическими моделями на основе входных данных, параметров модели, обучающих данных, контекста и конфигурации системы, и могут быть неточными, неполными, устаревшими, вымышленными, предвзятыми, неприменимыми или необъяснимыми.\n" +
				"Выходные данные не представляют собой юридическую, финансовую, медицинскую, инвестиционную, инженерную, техническую консультацию или заключение о соответствии требованиям; если пользователю необходимо принять важные решения на основе выходных данных, он должен самостоятельно провести проверку с участием человека и проконсультироваться с квалифицированными специалистами.\n" +
				"Пользователь несет полную ответственность за свои входные данные, способы использования услуг, а также за проверку, публикацию, распространение, применение и последствия использования выходных данных.\n" +
				"Платформа не гарантирует, что выходные данные ИИ являются абсолютно точными, что услуги будут предоставляться непрерывно или отвечать конкретным целям пользователя, а также что выходные данные не будут идентичными или схожими с контентом, сгенерированным другими пользователями.\n" +
				"В случае использования пользователем выходных данных в коммерческих целях, для открытой публикации, рекламы, тендерной документации, юридических документов, медицинских рекомендаций, инвестиционных советов, инженерного проектирования, автоматизированного принятия решений или в иных сценариях с высоким уровнем риска, пользователь самостоятельно несет риски соответствия законодательству и использования.",

			"Статья 9 Интеллектуальная собственность\n" +
				"Платформа, ее аффилированные лица и лицензиары обладают всеми правами в отношении системы платформы, программного обеспечения, моделей, алгоритмов, интерфейсов, документации, дизайна страниц, товарных знаков, логотипов, баз данных, технических решений, бизнес-моделей и связанных с ними прав интеллектуальной собственности.\n" +
				"Покупка пользователем AI токенов не дает права собственности на модели платформы, алгоритмы, программное обеспечение, исходный код, архитектуру системы, товарные знаки или иные объекты интеллектуальной собственности.\n" +
				"Пользователь сохраняет за собой законные права на законно предоставленные им входные данные.\n" +
				"При условии соблюдения пользователем настоящего соглашения и оплаты им соответствующих сборов платформа в пределах своих прав разрешает пользователю использовать выходные данные, сгенерированные услугой. Однако вопросы о том, являются ли выходные данные исключительными, нарушают ли они права третьих лиц, могут ли быть зарегистрированы как объекты интеллектуальной собственности или использованы в коммерческих целях, должны решаться пользователем самостоятельно на основе конкретного содержания и применимого законодательства; риск нарушения прав относится исключительно на счет пользователя.\n" +
				"Пользователь разрешает платформе обрабатывать входные данные, выходные данные и данные об использовании в пределах, необходимых для оказания услуг, тарификации, управления рисками, проверки безопасности, устранения неполадок, оптимизации услуг, аудита соответствия и выполнения юридических обязательств.\n" +
				"Если платформа использует сторонние базовые модели, облачные службы, платежные сервисы или иные сторонние услуги, соответствующие ограничения прав, политики контента и ограничения использования могут применяться также и к пользователю.",

			"Статья 10 Защита данных и конфиденциальность\n" +
				"Внутренние пользователи: Применяется Закон КНР «О защите персональной информации», соблюдается принцип минимизации и необходимости; платформа собирает только информацию, необходимую для оказания услуг: ① регистрационные данные (номер телефона, адрес электронной почты); ② платежная информация (номер заказа Alipay, сумма платежа); ③ данные устройства (IP-адрес, тип браузера); ④ данные об использовании (записи о пополнении и расходе токенов). Платформа не собирает конфиденциальную информацию, такую как номера удостоверений личности или банковских карт.\n" +
				"Зарубежные пользователи: Применяется законодательство о защите данных по месту их нахождения (например, GDPR); платформа обрабатывает персональную информацию в соответствии с требованиями законодательства страны пользователя, не осуществляет трансграничную передачу данных в регионы, не обеспечивающие надлежащие стандарты защиты данных, и при необходимости может подписать с пользователем соглашение об обработке данных.\n" +
				"Срок хранения персональной информации пользователя: 1 год с даты удаления аккаунта или прекращения обслуживания, по истечении которого данные автоматически обезличиваются или удаляются; записи заказов, налоговые записи, записи о пополнении и расходе токенов хранятся в течение 3 лет (для удовлетворения требований налогового аудита и комплаенса).\n" +
				"Если пользователь передает платформе персональные данные, чувствительные персональные данные, коммерческую тайну, конфиденциальную информацию или регулируемые данные других лиц, он должен гарантировать, что имеет на это законные основания и выполнил необходимые обязательства по уведомлению, получению согласию, получению разрешений, оценке или регистрации.\n" +
				"Платформа может в необходимых пределах поручить обработку данных сторонним поставщикам услуг, соответствующим требованиям, включая облачных провайдеров, поставщиков моделей, платежные организации, провайдеров SMS-услуг, сервисы выставления счетов, системы поддержки клиентов, провайдеров безопасности данных и т.д., и требовать от них принятия разумных мер защиты данных.",

			"Статья 11 Налоги, счета-фактуры и трансграничные платежи\n" +
				"Пользователь должен оплатить сборы в размере, указанном на странице заказа. Вопрос о том, включают ли соответствующие сборы налоги, определяется информацией на странице платформы, в описании заказа или счета-фактуры.\n" +
				"Дополнительные расходы, возникающие в связи с местом нахождения пользователя, способом оплаты, платежным каналом, банком-эквайером, банком, конвертацией валюты, НДС, потребительским налогом, налогом с продаж, налогом у источника, налогом на цифровые услуги или иными налогами и сборами, распределяются между сторонами в соответствии с применимым законодательством и условиями заказа. При оплате в иностранной валюте фактически списанная сумма может отличаться из-за курса обмена, комиссий банков, комиссий платежных систем или времени расчета; указанная разница относится на счет пользователя, если платформа не пообещала иное.\n" +
				"Если пользователю требуются счета-фактуры, квитанции или налоговые документы, он должен предоставить достоверные, точные и полные реквизиты по запросу платформы; в случае невозможности выставления счетов-фактур или необходимости их корректировки из-за ошибок пользователя или несовпадения плательщика и лица, на имя которого выставляется счет, пользователь должен оказать содействие в исправлении и несет соответствующие убытки.",

			"Статья 12 Изменение, прерывание и обслуживание услуг\n" +
				"Платформа может корректировать содержание услуг, типы моделей, правила тарификации, коэффициенты расхода токенов, ограничения на использование и объем функций в соответствии с развитием бизнеса, обновлением технологий, доступностью моделей, изменением затрат, законами и правилами, требованиями регулирующих органов или изменениями в услугах третьих лиц. Однако если указанные изменения приведут к существенному сокращению объема использования токенов, приобретенных, но не израсходованных пользователем, или к существенному увеличению удельной стоимости услуг (рост коэффициента расхода токенов более чем на 20%), пользователь имеет право до вступления изменений в силу потребовать пропорционального возврата стоимости неиспользованных токенов; указанные изменения не влияют на права использования приобретенных, но не израсходованных токенов в рамках первоначального объема услуг.\n" +
				"При внесении существенных изменений платформа должна уведомить пользователя посредством объявлений на сайте, системных уведомлений, электронной почты, SMS или других разумных способов.\n" +
				"В случае прерывания, задержки или недоступности услуг вследствие технического обслуживания, обновления системы, сбоев в сети, сбоев облачных служб, сбоев сторонних служб моделей, сбоев платежных каналов, инцидентов безопасности, форс-мажора или требований регуляторов, платформа приложит разумные усилия для восстановления обслуживания, но не несет ответственности за косвенные убытки.\n" +
				"Платформа не несет ответственности за сбои при использовании, вызванные сетью пользователя, оборудованием, конфигурацией интерфейсов, управлением ключами, ошибками в коде, ошибками параметров вызова, нарушением правил использования или действиями третьих лиц.",

			"Статья 13 Ограничение учетной записи и меры при нарушении\n" +
				"В случае нарушения пользователем настоящего соглашения или применимого законодательства платформа имеет право в зависимости от ситуации потребовать устранения нарушений, удаления запрещенного контента, приостановить соответствующие функции вызова и т. д.\n" +
				"Платформа может ограничить некоторые функции, ограничить количество параллельных запросов, снизить лимиты, приостановить действие API-ключей, заморозить учетные записи, приостановить или прекратить обслуживание.\n" +
				"Платформа может списать, заморозить или аннулировать токены, полученные или использованные с нарушением правил; в отношении заказов, созданных с нарушением правил, платформа имеет право отказать в возврате средств, отменить скидки и взыскать убытки.\n" +
				"В пределах, разрешенных законом, платформа может предоставлять необходимую информацию регулирующим органам, судебным органам, правообладателям или пострадавшим сторонам, а также принимать иные меры, разрешенные законом.\n" +
				"Если нарушение пользователем соглашения приведет к возникновению убытков, претензий, штрафов, расследований, судебных исков, расходов на адвокатов, нотариусов, экспертизы, командировочных расходов и т.д. у платформы, ее аффилированных лиц, партнеров, других пользователей или третьих лиц, пользователь несет полную ответственность по возмещению ущерба.",

			"Статья 14 Ограничение ответственности\n" +
				"В максимальной степени, разрешенной законом, ответственность платформы за предоставление услуг ограничивается суммой, фактически уплаченной пользователем за спорный заказ, которая не была израсходована и возвращена.\n" +
				"Платформа не несет ответственности за косвенные убытки, упущенную выгоду, ущерб деловой репутации, перерывы в работе, потерю данных, расходы на приобретение заменяющих услуг, претензии третьих лиц, потерю ожидаемых доходов или штрафные убытки, за исключением случаев, когда иное предусмотрено императивными нормами закона.\n" +
				"В отношении бесплатных токенов, пробных токенов, бонусных токенов и токенов по акциям платформа предоставляет их на условиях «как есть», не гарантирует их доступность, стабильность, непрерывность или ответственность за компенсацию, за исключением случаев, когда иное предусмотрено императивными нормами закона.\n" +
				"Настоящая Статья не ограничивает ответственность, возникшую вследствие умысла, грубой неосторожности платформы, мошенничества, причинения вреда жизни или здоровью, либо в иных случаях, когда ответственность не может быть ограничена по закону.",

			"Статья 15 Форс-мажор\n" +
				"Если сторона не может исполнить или задерживает исполнение настоящего соглашения вследствие стихийных бедствий, военных действий, массовых беспорядков, террористических актов, действий государственных органов, изменений в законодательстве, требований регуляторов, кибератак, хакерских инцидентов, вирусов, сбоев в работе сетей связи, сбоев облачных служб, отключения электроэнергии, сбоев платежных каналов, приостановки услуг сторонних моделей, международных санкций, экспортного контроля, эпидемий или иных непредвиденных, неизбежных и непреодолимых форс-мажорных обстоятельств, указанная сторона освобождается от ответственности частично или полностью в пределах влияния таких обстоятельств.\n" +
				"Сторона, пострадавшая от форс-мажорных обстоятельств, должна в разумный срок уведомить другую сторону и приложить разумные усилия для минимизации убытков.",

			"Статья 16 Изменение соглашения\n" +
				"Платформа может обновлять настоящее соглашение в связи с изменениями законов и правил, требований регуляторов, корректировкой бизнеса или изменениями услуг.\n" +
				"Обновленное соглашение будет представлено пользователю посредством объявлений на сайте, страниц оплаты, системных уведомлений, электронной почты или других разумных способов.\n" +
				"Если после обновления соглашения пользователь продолжает приобретать, пополнять или использовать AI токены, считается, что он принял обновленное соглашение.\n" +
				"Если изменения соглашения оказывают существенное негативное влияние на существующие основные права пользователя в отношении приобретенных, но не израсходованных токенов (существенное сокращение количества токенов, значительное сокращение срока действия или полный запрет на использование моделей, которые были обещаны изначально), платформа должна предоставить разумное уведомление и может предложить альтернативное решение равной стоимости либо вернуть стоимость неиспользованных токенов пропорционально их количеству, если иное не предусмотрено императивными нормами закона.",

			"Статья 17 Применимое право и разрешение споров\n" +
				"Внутренние пользователи: К заключению, действительности, исполнению, толкованию и разрешению споров по настоящему соглашению применяется законодательство Китайской Народной Республики (за исключением коллизионных норм); споры, возникшие из настоящего соглашения или в связи с ним, должны сначала решаться сторонами путем дружеских переговоров; при недостижении согласия любая из сторон может передать спор на рассмотрение в народный суд по месту нахождения поставщика услуг.\n" +
				"Зарубежные пользователи: Настоящее соглашение регулируется законодательством по месту нахождения поставщика услуг (при одновременном соблюдении императивных норм страны нахождения пользователя); споры, возникшие из настоящего соглашения или в связи с ним, должны сначала решаться сторонами путем дружеских переговоров; при недостижении согласия спор передается на рассмотрение в Китайскую международную экономическую и торговую арбитражную комиссию (CIETAC) для проведения арбитража в соответствии с ее правилами арбитража, действующими на момент подачи заявления; местом проведения арбитража является Шанхай (Китай), язык арбитража — китайский / английский, арбитражное решение является окончательным и обязательным для обеих сторон.\n" +
				"Если в месте нахождения пользователя-потребителя установлены императивные правила о подсудности, применяются такие правила.",

			"Статья 18 Уведомления\n" +
				"Платформа может направлять пользователю уведомления посредством объявлений на сайте, системных уведомлений, электронной почты, SMS, телефонных звонков, уведомлений в консоли управления API, уведомлений о выставлении счетов или по контактным данным, предоставленным пользователем.\n" +
				"Пользователь должен гарантировать, что его контактные данные являются достоверными, точными и действующими; в случае изменения контактных данных пользователь должен своевременно обновить их.\n" +
				"Правила доставки уведомлений: Считаются доставленными в момент публикации объявления на сайте; в момент успешной отправки для системных сообщений, электронной почты, SMS; в момент подписания при доставке курьерской службой или заказным письмом, либо считаются доставленными в случае отказа от получения или невозможности доставки по причинам, зависящим от пользователя.",

			"Статья 19 Прочие положения\n" +
				"Настоящее соглашение составляет полное соглашение между сторонами по вопросам покупки и использования AI токенов и заменяет любые предшествующие устные или письменные договоренности между сторонами.\n" +
				"Если какое-либо положение настоящего соглашения будет признано недействительным, незаконным или не подлежащим исполнению, это не влияет на действительность остальных положений.\n" +
				"Пользователь не имеет права передавать свои права и обязанности по настоящему соглашению без письменного согласия платформы.\n" +
				"Платформа может передать свои права и обязанности по настоящему соглашению в случае реорганизации бизнеса, слияния, разделения, передачи активов, перехода права к аффилированным лицам или миграции услуг, но должна обеспечить, чтобы существующие права и интересы пользователей не были существенно ущемлены.\n" +
				"Заголовки настоящего соглашения приведены исключительно для удобства чтения и не влияют на толкование положений.\n" +
				"Настоящее соглашение может быть составлено на китайском, английском или других языках; в случае расхождений между различными языковыми версиями преимущественную силу имеет китайская версия, если иное не требуется императивными нормами законодательства страны нахождения пользователя.\n" +
				"Пользователь должен самостоятельно скачать, сохранить или распечатать настоящее соглашение перед оплатой; платформа также может предоставить возможность просмотра версии соглашения, истории заказов или транзакций в личном кабинете пользователя.",
		},
	},
	"vi": {
		Title: "Thỏa thuận mua và sử dụng AI Token",
		MetaKeys: map[string]string{
			"version":     "Phiên bản: V1.1",
			"effective":   "Ngày hiệu lực: {{date}}",
			"provider":    "Bên cung cấp dịch vụ (Nền tảng): OSS Energietechnik GmbH",
			"address":     "Địa chỉ đăng ký: Adam-Opel-Straβe 16-18, 60386 Frankfurt am Main",
			"email":       "Email liên hệ: info@oss-energietechnik.de",
			"user":        "Người dùng: {{user}}",
			"order":       "Đơn hàng thanh toán: {{order}} ({{amount}})",
			"signMethod":  "Phương thức ký kết: Người dùng tích chọn đồng ý với thỏa thuận này và nhấp vào các nút tương tự như \"Đồng ý thỏa thuận và thanh toán\" hoặc \"Xác nhận mua\" được coi là ký kết điện tử hợp lệ, thỏa thuận có hiệu lực ngay lập tức.",
			"applicable":  "Khu vực áp dụng: Áp dụng cho cả người dùng trong nước (Trung Quốc đại lục) và ngoài nước. Người dùng trong nước ưu tiên áp dụng luật pháp Cộng hòa Nhân dân Trung Hoa; người dùng ngoài nước đồng thời tuân thủ luật pháp bắt buộc tại nơi cư trú liên quan đến bảo vệ người tiêu dùng, bảo vệ dữ liệu, thuế, thương mại điện tử, kiểm soát xuất khẩu và dịch vụ kỹ thuật số.",
			"declaration": "Tuyên bố quan trọng: Thỏa thuận này là hợp đồng thương mại, không phải ý kiến pháp lý; trước khi chính thức đưa lên hệ thống, khuyến nghị nên hoàn tất việc kiểm tra dưới cùng kết hợp với quyền tài phán của người dùng, kênh thanh toán, yêu cầu tuân thủ về thuế và dữ liệu.",
		},
		Sections: []string{
			"Điều 1 Định nghĩa\n" +
				"AI Token: Là hạn mức gọi dịch vụ / chứng từ tiêu thụ do người dùng mua từ nền tảng để sử dụng cho việc gọi mô hình trí tuệ nhân tạo, năng lượng thuật toán, dịch vụ API, dịch vụ agent, tạo văn bản, tạo hình ảnh, xử lý âm thanh, phân tích dữ liệu hoặc các chức năng AI khác của nền tảng, không thuộc về tiền tệ, tiền ảo, tiền điện tử, chứng khoán, sản phẩm tài chính, thẻ tích lũy giá trị hoặc tài sản có thể chuyển nhượng.\n" +
				"Dịch vụ: Là các dịch vụ liên quan do nền tảng cung cấp cho người dùng như gọi mô hình AI, kết nối API, tạo trực tuyến, xử lý dữ liệu, quản lý tài khoản, thống kê lượng sử dụng, hỗ trợ kỹ thuật.\n" +
				"Tài khoản: Là tài khoản, API Key, khóa mật khẩu, ID tổ chức, ID dự án hoặc nhận dạng danh tính khác do người dùng đăng ký, đăng nhập trên nền tảng hoặc do nền tảng phân bổ.\n" +
				"Nội dung đầu vào: Là các nội dung như văn bản, hình ảnh, âm thanh, video, mã, tệp, dữ liệu, từ gợi ý (prompt), tham số giao diện do người dùng gửi, tải lên, nhập, truyền hoặc gọi khi sử dụng dịch vụ.\n" +
				"Nội dung đầu ra: Là văn bản, hình ảnh, âm thanh, video, mã, kết quả phân tích, câu trả lời của mô hình hoặc các kết quả khác do dịch vụ tạo ra hoặc trả về dựa trên nội dung đầu vào của người dùng.\n" +
				"Tiêu thụ: Là hành vi khấu trừ AI Token do người dùng gọi dịch vụ AI, gửi yêu cầu, tạo nội dung, xử lý dữ liệu hoặc sử dụng các chức năng liên quan.\n" +
				"Tính tuân thủ: AI Token này chỉ là hạn mức gọi dịch vụ AI độc quyền của nền tảng, không có tính lưu thông, tính pháp lý, không được đổi lấy tiền pháp định, tiền ảo hoặc dùng cho giao dịch đầu cơ; trong nước áp dụng quy định quản lý của Ngân hàng Trung ương, ngoài nước đồng thời phù hợp với quy tắc quản lý chống tiền ảo và tài sản số tại nơi cư trú của người dùng.",

			"Điều 2 Ký kết thỏa thuận và xác nhận điện tử\n" +
				"Thỏa thuận này được hiển thị, xác nhận, ký kết, lưu trữ dưới dạng dữ liệu điện tử, có hiệu lực pháp lý tương đương với thỏa thuận bằng văn bản.\n" +
				"Hành vi người dùng tích chọn đồng ý thỏa thuận này và nhấp vào thanh toán, xác nhận mua, gửi đơn hàng, hoàn tất thanh toán hoặc thực tế sử dụng AI Token cấu thành sự xác nhận và ký kết có hiệu lực đối với thỏa thuận này.\n" +
				"Nền tảng có quyền ghi lại và lưu giữ thông tin tài khoản, mã đơn hàng, lịch sử thanh toán, địa chỉ IP, thông tin thiết bị, thông tin trình duyệt, thời gian thao tác, số phiên bản thỏa thuận, lịch sử tích chọn, lịch sử nhấp chuột, nhật ký hệ thống khi người dùng xác nhận thỏa thuận này, dùng để chứng minh việc thiết lập, thực hiện thỏa thuận, giải quyết tranh chấp và kiểm toán tuân thủ.\n" +
				"Người dùng không được phủ nhận hiệu lực của thỏa thuận này với lý do chưa ký hợp đồng giấy, chưa đóng dấu thực tế hoặc chưa ký tên trực tiếp.\n" +
				"Nếu người dùng là doanh nghiệp, tổ chức hoặc đoàn thể khác, người thay mặt đơn vị đó thực hiện mua hoặc sử dụng dịch vụ xác nhận đã được ủy quyền đầy đủ; trường hợp chưa được ủy quyền mà thay mặt người khác mua hoặc sử dụng dịch vụ thì người thao tác thực tế phải tự chịu trách nhiệm tương ứng.\n" +
				"Các điều khoản miễn trừ / giới hạn trách nhiệm của nền tảng như hoàn tiền, giới hạn trách nhiệm trong thỏa thuận này đã được nhắc nhở nổi bật bằng cách in đậm trong văn bản, hiển thị cửa sổ độc lập tại trang thanh toán, người dùng xác nhận đã hiểu rõ và tự nguyện chấp nhận.",

			"Điều 3 Mua, định giá và bàn giao Token\n" +
				"Người dùng có thể mua AI Token căn cứ theo gói dịch vụ, giá cả, loại tiền tệ, số lượng, thời hạn hiệu lực, mô hình áp dụng, quy tắc tiêu thụ và các mô tả khác hiển thị trên trang của nền tảng.\n" +
				"Giá cả, tỷ lệ đổi, quy tắc tiêu thụ, mô hình hỗ trợ, độ dài ngữ cảnh, giới hạn cuộc gọi đồng thời, phạm vi chức năng của AI Token sẽ căn cứ theo thông tin hiển thị trên nền tảng tại thời điểm người dùng đặt hàng; trường hợp nền tảng công bố điều chỉnh thì thực hiện theo nội dung công bố.\n" +
				"Trừ khi nền tảng có mô tả rõ ràng khác, AI Token chỉ được dùng cho dịch vụ độc quyền trong nền tảng, không thuộc về tài sản tài chính, không thể chuyển nhượng, giao dịch, quy đổi tài sản khác.\n" +
				"Nền tảng sau khi nhận được tiền thanh toán của người dùng và xác nhận đơn hàng, sẽ nạp AI Token tương ứng vào tài khoản người dùng trong vòng 24 giờ; quá thời hạn mà chưa vào tài khoản, người dùng có quyền xin hoàn tiền toàn bộ vô điều kiện.\n" +
				"Đối với việc chậm trễ vào tài khoản, khấu trừ phí dịch vụ do kênh thanh toán, ngân hàng, tổ chức thanh toán bên thứ ba, quyết toán ngoại hối, kiểm tra chống gian lận, kiểm tra thuế hoặc kiểm tra tuân thủ, nền tảng không chịu trách nhiệm do bên thứ ba gây ra, nhưng sẽ hỗ trợ người dùng truy vấn trong phạm vi hợp lý.\n" +
				"Người dùng phải đảm bảo thông tin đơn hàng, thông tin tài khoản, thông tin hóa đơn, thông tin chủ thể thanh toán là xác thực, chính xác và đầy đủ; trường hợp do thông tin người dùng bị sai dẫn đến nạp tiền thất bại, hóa đơn bị sai, tranh chấp quyền sở hữu tài khoản hoặc tổn thất khác thì người dùng tự chịu trách nhiệm.",

			"Điều 4 Quy tắc sử dụng Token\n" +
				"AI Token do người dùng mua chỉ được dùng trong phạm vi chỉ định của nền tảng, cụ thể dùng được cho mô hình, dịch vụ, API, chức năng hoặc sản phẩm nào sẽ căn cứ theo trang nền tảng, trang điều khiển hoặc mô tả đơn hàng.\n" +
				"AI Token được tiêu thụ theo phương thức tính toán do nền tảng công bố, bao gồm nhưng không giới hạn số lượng ký tự đầu vào, số lượng ký tự đầu ra, số token, số lượng ảnh, thời lượng âm thanh, thời lượng video, số lần yêu cầu, loại mô hình, tài nguyên tính toán, tài nguyên lưu trữ, gọi plugin, gọi công cụ.\n" +
				"Mức tiêu thụ Token của các mô hình khác nhau, chức năng khác nhau, nút khu vực khác nhau, cấp độ dịch vụ khác nhau có thể khác nhau, người dùng cần xem mô tả tính phí liên quan trước khi dùng.\n" +
				"AI Token một khi đã tiêu thụ, coi như dịch vụ đã bàn giao hoặc bàn giao một phần; trừ khi thỏa thuận này có quy định khác hoặc luật pháp bắt buộc, AI Token đã tiêu thụ không hỗ trợ khôi phục, hoàn lại, chuyển nhượng hoặc quy đổi tiền mặt.\n" +
				"Người dùng phải bảo quản an toàn tài khoản, mật khẩu, API Key, khóa truy cập, mã xác nhận và các thông tin xác nhận danh tính khác; các cuộc gọi bắt nguồn từ tài khoản hoặc khóa mật khẩu của người dùng đều coi là hành vi của chính người dùng hoặc người được ủy quyền.\n" +
				"Trường hợp người dùng phát hiện tài khoản hoặc khóa mật khẩu bị rò rỉ, bị đánh cắp hoặc có cuộc gọi bất thường, phải thông báo ngay cho nền tảng và thực hiện các biện pháp như reset mật khẩu, vô hiệu hóa khóa, đóng cổng kết nối; lượng Token tiêu thụ phát sinh trước khi nền tảng nhận được thông báo về nguyên tắc do người dùng chịu, trừ trường hợp nền tảng có hành vi cố ý hoặc lỗi vô ý nghiêm trọng.\n" +
				"Người dùng không được chuyển nhượng, bán, cho thuê, cho mượn, rút tiền mặt, quy đổi hoặc giao dịch AI Token ngoài nền tảng, không được dùng chúng cho giao dịch hoặc hoạt động tài trợ vốn bên ngoài nền tảng.",

			"Điều 5 Thời hạn hiệu lực, hết hạn và nạp tiếp\n" +
				"Thời hạn hiệu lực của AI Token căn cứ theo thông tin hiển thị rõ ràng trên trang mua hàng, trang đơn hàng, mô tả gói dịch vụ hoặc trang quản lý tài khoản người dùng của nền tảng; nếu không hiển thị rõ ràng thì thời hạn hiệu lực mặc định là 12 tháng kể từ ngày nạp vào tài khoản.\n" +
				"Sau khi hết thời hạn hiệu lực, AI Token chưa dùng sẽ tự động mất hiệu lực, nền tảng không cung cấp dịch vụ sử dụng, hoàn trả, quy đổi hoặc gia hạn nữa, trừ trường hợp pháp luật bắt buộc có yêu cầu khác.\n" +
				"Nền tảng có thể cung cấp dịch vụ gia hạn, tiếp tục nạp, nâng cấp, chuyển đổi gói dịch vụ theo kế hoạch vận hành, quy tắc cụ thể căn cứ theo hiển thị trên nền tảng tại thời điểm đó.\n" +
				"Sau khi người dùng mua gói dịch vụ mới, thứ tự sử dụng Token cũ và mới thực hiện theo quy tắc hệ thống của nền tảng; nếu không có mô tả đặc biệt, nền tảng có thể ưu tiên tiêu thụ Token sắp hết hạn.",

			"Điều 6 Chính sách hoàn tiền và rút lui\n" +
				"Quy tắc cốt lõi: AI Token là dịch vụ kỹ thuật số bàn giao tức thì, tiêu thụ tức thì, không áp dụng trả hàng không lý do trong vòng 7 ngày; sau khi thanh toán vào tài khoản về nguyên tắc không hoàn tiền không lý do, trừ quy định tại khoản 2, 3, 4 Điều này và quy định bắt buộc của pháp luật.\n" +
				"Quy tắc hoàn tiền cho người dùng trong nước: (1) Sau khi nạp tiền vào tài khoản, không áp dụng quy tắc hoàn tiền, không hoàn tiền; (2) Các Token AI đã tiêu thụ, Token tặng, Token hoạt động, Token dùng thử, phần khấu trừ bằng coupon, Token bị hạn chế sử dụng do vi phạm thỏa thuận này, Token bị tiêu thụ do thao tác sai của người dùng (như nạp nhầm, gọi nhầm) không được hoàn tiền, rút tiền mặt, chuyển nhượng hoặc quy đổi tiền mặt.\n" +
				"Quy tắc hoàn tiền cho người dùng ngoài nước (bao gồm cả EU, Vương quốc Anh và các khu vực tài phán khác): (1) Nơi cư trú của người dùng có quyền rút lui theo luật định, việc tích chọn xác nhận của người dùng về \"thực hiện ngay, nạp tiền đồng nghĩa mất quyền rút lui\" sẽ có hiệu lực; (2) Nếu chưa xác nhận thực hiện ngay và chưa tiêu thụ Token thì thực hiện theo thời hạn rút lui luật định tại nơi cư trú của người dùng; phần đã tiêu thụ không được hoàn tiền.\n" +
				"Hoàn tiền do lỗi nền tảng: Do nguyên nhân từ nền tảng dẫn đến người dùng không thể sử dụng dịch vụ cốt lõi tương ứng với Token đã mua và chưa hết hạn trong 30 ngày liên tục, nền tảng sau khi nhận được thông báo có hiệu lực của người dùng trong vòng 15 ngày làm việc sẽ đưa ra biện pháp khắc phục như đền bù Token có giá trị tương đương, gia hạn thời gian hiệu lực hoặc hoàn tiền phần chưa sử dụng; quá thời hạn mà chưa khắc phục, người dùng có quyền yêu cầu hoàn tiền theo tỷ lệ Token chưa sử dụng.",

			"Điều 7 Nghĩa vụ tuân thủ của người dùng\n" +
				"Người dùng cam kết tuân thủ luật pháp Cộng hòa Nhân dân Trung Hoa và các luật pháp quy định, yêu cầu quản lý, quy phạm ngành nghề, kiểm soát xuất khẩu, lệnh trừng phạt kinh tế, quy tắc truyền dữ liệu qua biên giới, trật tự công cộng và thuần phong mỹ tục tại nơi cư trú của người dùng.\n" +
				"Người dùng không được tạo, truyền bá hoặc hỗ trợ tạo các nội dung vi phạm như nội dung bất hợp pháp, xâm hại quyền lợi, lừa đảo, giả mạo, thù hận, quấy rối, bạo lực, khiêu dâm, khủng bố, cực đoan, tự hại, tự sát, ma túy, vũ khí, mã độc, tấn công mạng, hoạt động tài chính bất hợp pháp.\n" +
				"Người dùng không được xâm phạm quyền sở hữu trí tuệ, bí mật thương mại, quyền riêng tư, quyền thông tin cá nhân, quyền hình ảnh, quyền danh dự hoặc các quyền và lợi ích hợp pháp khác của người khác.\n" +
				"Người dùng không được thu thập (scrape), sao chép, huấn luyện, đảo ngược kỹ thuật, bẻ khóa, vượt qua giới hạn kỹ thuật, cơ chế bảo mật, kiểm soát truy cập hoặc hệ thống tính phí của nền tảng khi chưa được phép.\n" +
				"Người dùng không được sử dụng dịch vụ cho mục đích phát tán thư rác tự động, tăng tương tác giả (刷量), đánh giá giả, lừa đảo trực tuyến (phishing), mạo danh người khác, vượt qua kiểm duyệt nội dung hoặc tạo hàng loạt nội dung vi phạm pháp luật.\n" +
				"Người dùng không được gửi lên nền tảng các dữ liệu, thông tin cá nhân, thông tin cá nhân nhạy cảm, thông tin bảo mật, bí mật quốc gia, bí mật thương mại hoặc dữ liệu bị hạn chế kiểm soát xuất khẩu mà mình không có quyền xử lý.\n" +
				"Nền tảng có quyền thực hiện kiểm tra an toàn, nhận diện kiểm soát rủi ro và xử lý tuân thủ cần thiết đối với nội dung đầu vào, nội dung đầu ra, hành vi gọi cuộc gọi và hành vi tài khoản của người dùng theo quy định của pháp luật, yêu cầu quản lý, quy tắc nền tảng, chiến lược kiểm soát rủi ro hoặc yêu cầu của nhà cung cấp mô hình bên thứ ba.",

			"Điều 8 Nội dung đầu ra AI và Cảnh báo rủi ro\n" +
				"Người dùng hiểu và đồng ý rằng nội dung đầu ra AI được tạo ra bởi mô hình thuật toán dựa trên nội dung đầu vào, tham số mô hình, dữ liệu huấn luyện, ngữ cảnh và cấu hình hệ thống, có thể tồn tại những trường hợp không chính xác, không đầy đủ, lỗi thời, giả tưởng, sai lệch, không áp dụng được hoặc không thể giải thích.\n" +
				"Nội dung đầu ra không cấu thành ý kiến chuyên môn về pháp lý, tài chính, y tế, đầu tư, kỹ thuật, an toàn, tuân thủ hoặc các ý kiến chuyên môn khác; nếu người dùng cần dựa vào nội dung đầu ra để đưa ra các quyết định quan trọng, cần tự mình tiến hành kiểm tra lại bằng con người và tham khảo ý kiến của các chuyên gia có trình độ.\n" +
				"Người dùng phải tự chịu toàn bộ trách nhiệm đối với nội dung đầu vào, cách thức sử dụng, việc kiểm duyệt, công bố, truyền bá, ứng dụng và hậu quả của nội dung đầu ra của mình.\n" +
				"Nền tảng không đảm bảo nội dung đầu ra AI chính xác hoàn toàn, dịch vụ luôn sẵn sàng liên tục, đáp ứng mục đích cụ thể của người dùng, cũng không đảm bảo nội dung đầu ra sẽ không trùng lặp hoặc tương tự với nội dung do người dùng khác tạo ra.\n" +
				"Đối với việc người dùng sử dụng nội dung đầu ra cho mục đích thương mại, công bố công khai, quảng cáo, hồ sơ thầu, tài liệu pháp lý, tư vấn y tế, tư vấn đầu tư, thiết kế kỹ thuật, quyết định tự động hoặc các tình huống rủi ro cao khác, người dùng phải tự chịu rủi ro về kiểm tra tuân thủ và sử dụng.",

			"Điều 9 Quyền sở hữu trí tuệ\n" +
				"Nền tảng cùng các bên liên kết, bên cấp phép của nền tảng sở hữu toàn bộ quyền đối với hệ thống nền tảng, phần mềm, mô hình, thuật toán, giao diện, tài liệu, thiết kế trang web, nhãn hiệu, logo, cơ sở dữ liệu, phương án kỹ thuật, mô hình kinh doanh và các quyền sở hữu trí tuệ liên quan.\n" +
				"Việc người dùng mua AI Token không đồng nghĩa với việc có được quyền sở hữu đối với mô hình, thuật toán, phần mềm, mã nguồn, cấu trúc hệ thống, nhãn hiệu hoặc các quyền sở hữu trí tuệ khác của nền tảng.\n" +
				"Người dùng bảo lưu các quyền hợp pháp đối với nội dung đầu vào do mình cung cấp một cách hợp pháp.\n" +
				"Với điều kiện người dùng tuân thủ thỏa thuận này và đã thanh toán các khoản phí tương ứng, nền tảng trong phạm vi quyền hạn của mình cho phép người dùng sử dụng nội dung đầu ra do dịch vụ tạo ra; tuy nhiên việc nội dung đầu ra có thể độc quyền hay không, có xâm phạm quyền lợi của bên thứ ba hay không, có thể đăng ký sở hữu trí tuệ hay không, có thể thương mại hóa hay không cần do người dùng tự đánh giá dựa trên nội dung cụ thể và luật pháp áp dụng, rủi ro xâm phạm quyền lợi do người dùng tự gánh chịu.\n" +
				"Người dùng ủy quyền cho nền tảng xử lý nội dung đầu vào, nội dung đầu ra và dữ liệu sử dụng trong phạm vi cần thiết cho việc cung cấp dịch vụ, tính phí, kiểm soát rủi ro, kiểm tra an toàn, khắc phục sự cố, tối ưu hóa dịch vụ, kiểm toán tuân thủ và thực hiện các nghĩa vụ pháp lý.\n" +
				"Nếu nền tảng sử dụng mô hình cơ sở bên thứ ba, dịch vụ đám mây, dịch vụ thanh toán hoặc các dịch vụ bên thứ ba khác, các giới hạn quyền lợi, chính sách nội dung, giới hạn sử dụng tương ứng có thể đồng thời áp dụng cho người dùng.",

			"Điều 10 Bảo vệ dữ liệu và Quyền riêng tư\n" +
				"Người dùng trong nước: Áp dụng \"Luật Bảo vệ thông tin cá nhân Cộng hòa Nhân dân Trung Hoa\", tuân thủ nguyên tắc tối thiểu và cần thiết; nền tảng chỉ thu thập thông tin cần thiết để cung cấp dịch vụ: ① thông tin đăng ký (số điện thoại, email); ② thông tin thanh toán (mã đơn hàng Alipay, số tiền thanh toán); ③ thông tin thiết bị (địa chỉ IP, loại trình duyệt); ④ thông tin lượng sử dụng (lịch sử nạp / tiêu thụ Token), không thu thập thông tin nhạy cảm như số CMND, số thẻ ngân hàng.\n" +
				"Người dùng ngoài nước: Áp dụng luật bảo vệ dữ liệu tại nơi cư trú của người dùng (như GDPR); nền tảng xử lý thông tin cá nhân theo yêu cầu luật pháp tại nơi cư trú của người dùng, không truyền dữ liệu xuyên biên giới đến khu vực không đạt tiêu chuẩn bảo vệ dữ liệu, khi cần thiết có thể ký thỏa thuận xử lý dữ liệu với người dùng.\n" +
				"Thời hạn lưu trữ thông tin cá nhân của người dùng: 1 năm kể từ ngày hủy tài khoản hoặc chấm dứt dịch vụ, sau khi hết hạn sẽ tự động ẩn danh hoặc xóa; hồ sơ đơn hàng, hồ sơ thuế, hồ sơ nạp / tiêu thụ Token được lưu trữ trong 3 năm (phù hợp với yêu cầu kiểm toán thuế và tuân thủ).\n" +
				"Trường hợp người dùng nhập vào nền tảng thông tin cá nhân, thông tin cá nhân nhạy cảm, bí mật thương mại, thông tin bảo mật hoặc dữ liệu chịu sự quản lý liên quan đến người khác, người dùng phải đảm bảo mình có cơ sở xử lý hợp pháp, và đã thực hiện nghĩa vụ thông báo, đồng ý, ủy quyền, đánh giá hoặc báo cáo cần thiết.\n" +
				"Nền tảng có thể ủy quyền cho nhà cung cấp dịch vụ bên thứ ba tuân thủ xử lý dữ liệu trong phạm vi cần thiết, bao gồm nhà cung cấp dịch vụ đám mây, dịch vụ mô hình, tổ chức thanh toán, dịch vụ tin nhắn, dịch vụ hóa đơn, hệ thống hỗ trợ khách hàng, dịch vụ an toàn dữ liệu, và yêu cầu các bên thứ ba liên quan thực hiện các biện pháp bảo vệ dữ liệu hợp lý.",

			"Điều 11 Thuế, Hóa đơn và Thanh toán xuyên biên giới\n" +
				"Người dùng phải thanh toán phí theo số tiền hiển thị trên trang đơn hàng. Phí liên quan có bao gồm thuế hay không sẽ căn cứ theo trang nền tảng, đơn hàng hoặc mô tả hóa đơn.\n" +
				"Các chi phí phát sinh thêm do nơi cư trú của người dùng, phương thức thanh toán, kênh thanh toán, tổ chức thanh toán, ngân hàng, quyết toán ngoại hối, thuế GTGT, thuế tiêu thụ, thuế bán hàng, thuế khấu trừ tại nguồn, thuế dịch vụ kỹ thuật số hoặc các loại thuế phí khác sẽ do hai bên tự gánh chịu theo quy định của pháp luật và thỏa thuận đơn hàng. Đối với các khoản thanh toán bằng ngoại tệ, số tiền thực tế bị khấu trừ có thể khác biệt do tỷ giá, phí giao dịch ngân hàng, phí kênh thanh toán hoặc thời gian quyết toán; sự khác biệt này do người dùng tự gánh chịu, trừ khi nền tảng có cam kết khác.\n" +
				"Trường hợp người dùng cần hóa đơn, biên lai hoặc chứng từ thuế, phải cung cấp thông tin xuất hóa đơn xác thực, chính xác và đầy đủ theo yêu cầu của nền tảng; do người dùng cung cấp sai thông tin hoặc chủ thể thanh toán không thống nhất với chủ thể nhận hóa đơn dẫn đến không xuất được hóa đơn hoặc cần sửa đổi hóa đơn, người dùng phải phối hợp xử lý và tự chịu tổn thất tương ứng.",

			"Điều 12 Thay đổi, Gián đoạn và Bảo trì dịch vụ\n" +
				"Nền tảng có thể điều chỉnh nội dung dịch vụ, loại mô hình, quy tắc tính phí, tỷ lệ tiêu thụ Token, giới hạn sử dụng, phạm vi chức năng theo sự phát triển kinh doanh, nâng cấp kỹ thuật, cung cấp mô hình, thay đổi chi phí, luật pháp quy định, yêu cầu quản lý hoặc thay đổi của dịch vụ bên thứ ba; tuy nhiên nếu thay đổi nêu trên dẫn đến phạm vi sử dụng của Token đã mua nhưng chưa sử dụng của người dùng bị thu hẹp đáng kể, chi phí dịch vụ đơn vị tăng lên đáng kể (tỷ lệ tiêu thụ Token tăng vượt quá 20%), người dùng có quyền xin hoàn tiền theo tỷ lệ chi phí tương ứng với Token chưa sử dụng trước khi thay đổi có hiệu lực; thay đổi nêu trên không ảnh hưởng đến quyền lợi sử dụng của Token đã mua nhưng chưa tiêu thụ trong phạm vi dịch vụ ban đầu.\n" +
				"Khi thực hiện thay đổi lớn, nền tảng sẽ thông báo cho người dùng qua thông báo trang web, tin nhắn trong trang, email, tin nhắn điện thoại hoặc phương thức hợp lý khác.\n" +
				"Do bảo trì hệ thống, nâng cấp, sự cố mạng, sự cố dịch vụ đám mây, sự cố dịch vụ mô hình bên thứ ba, sự cố kênh thanh toán, sự kiện an toàn, bất khả kháng hoặc yêu cầu quản lý dẫn đến dịch vụ bị gián đoạn, chậm trễ, không khả dụng, nền tảng sẽ nỗ lực hợp lý để khôi phục dịch vụ, và không chịu trách nhiệm đối với tổn thất gián tiếp.\n" +
				"Nền tảng không chịu trách nhiệm đối với việc sử dụng thất bại do mạng của người dùng, thiết bị, cấu hình giao diện, quản lý khóa bảo mật, lỗi mã nguồn, lỗi tham số gọi, vi phạm quy tắc sử dụng hoặc nguyên nhân từ bên thứ ba.",

			"Điều 13 Hạn chế tài khoản và Xử lý vi phạm hợp đồng\n" +
				"Nếu người dùng vi phạm thỏa thuận này hoặc quy định pháp luật áp dụng, nền tảng có quyền tùy theo tình hình để thực hiện các biện pháp như nhắc nhở sửa đổi, yêu cầu xóa nội dung vi phạm, tạm ngưng chức năng gọi liên quan.\n" +
				"Nền tảng có thể hạn chế một phần chức năng, hạn chế số cuộc gọi đồng thời, hạ thấp hạn mức, tạm ngưng API Key, đóng băng tài khoản, tạm ngưng dịch vụ hoặc chấm dứt dịch vụ.\n" +
				"Nền tảng có thể khấu trừ, đóng băng hoặc hủy bỏ Token có được hoặc sử dụng vi phạm quy định; đối với các đơn hàng phát sinh do vi phạm, nền tảng có quyền từ chối hoàn tiền, hủy bỏ ưu đãi và thu hồi tổn thất.\n" +
				"Trong phạm vi pháp luật cho phép, nền tảng có thể cung cấp thông tin cần thiết cho cơ quan quản lý, cơ quan tư pháp, chủ thể quyền lợi hoặc bên bị hại, và thực hiện các biện pháp khác được pháp luật cho phép.\n" +
				"Trường hợp hành vi vi phạm hợp đồng của người dùng dẫn đến nền tảng, bên liên kết, bên hợp tác, người dùng khác hoặc bên thứ ba chịu tổn thất, khiếu nại, xử phạt, điều tra, kiện tụng, phí luật sư, phí công chứng, phí giám định, phí đi lại, v.v., người dùng phải chịu toàn bộ trách nhiệm bồi thường.",

			"Điều 14 Giới hạn trách nhiệm\n" +
				"Trong phạm vi tối đa pháp luật cho phép, trách nhiệm của nền tảng đối với dịch vụ được giới hạn tối đa bằng số tiền thực tế mà người dùng đã thanh toán cho đơn hàng phát sinh tranh chấp nhưng chưa tiêu thụ và chưa hoàn tiền.\n" +
				"Nền tảng không chịu trách nhiệm đối với tổn thất gián tiếp, tổn thất lợi nhuận, tổn thất uy tín thương mại, gián đoạn kinh doanh, mất dữ liệu, chi phí mua dịch vụ thay thế, khiếu nại của bên thứ ba, tổn thất doanh thu kỳ vọng hoặc bồi thường mang tính trừng phạt, trừ khi pháp luật bắt buộc có yêu cầu khác.\n" +
				"Đối với Token miễn phí, Token dùng thử, Token tặng, Token hoạt động, nền tảng cung cấp theo hiện trạng \"as is\", không cam kết tính khả dụng, tính ổn định, tính liên tục hoặc trách nhiệm bồi thường, trừ khi pháp luật bắt buộc có yêu cầu khác.\n" +
				"Điều này không giới hạn trách nhiệm phát sinh do hành vi cố ý, lỗi vô ý nghiêm trọng, lừa đảo, thiệt hại về người của nền tảng hoặc các trường hợp pháp luật không cho phép giới hạn trách nhiệm.",

			"Điều 15 Bất khả kháng\n" +
				"Do thiên tai, chiến tranh, bạo loạn, tấn công khủng bố, hành vi của chính phủ, thay đổi luật pháp quy định, yêu cầu quản lý, tấn công mạng, sự kiện hacker, virus, sự cố viễn thông cơ bản, sự cố dịch vụ đám mây, gián đoạn năng lượng, sự sự cố kênh thanh toán, gián đoạn dịch vụ mô hình bên thứ ba, trừng phạt quốc tế, kiểm soát xuất khẩu, dịch bệnh hoặc các sự kiện bất khả kháng khác không thể dự báo, không thể tránh khỏi và không thể khắc phục dẫn đến một bên không thể thực hiện hoặc chậm thực hiện thỏa thuận này, bên đó có thể được miễn trừ một phần hoặc toàn bộ trách nhiệm trong phạm vi bị ảnh hưởng.\n" +
				"Bên chịu ảnh hưởng của sự kiện bất khả kháng phải thông báo cho bên kia trong thời gian hợp lý, và nỗ lực hợp lý để giảm thiểu tổn thất.",

			"Điều 16 Thay đổi thỏa thuận\n" +
				"Nền tảng có thể cập nhật thỏa thuận này theo quy định của pháp luật, yêu cầu quản lý, điều chỉnh kinh doanh hoặc thay đổi dịch vụ.\n" +
				"Thỏa thuận sau khi cập nhật sẽ được hiển thị cho người dùng thông qua thông báo trên trang web, trang thanh toán, tin nhắn trong trang, email hoặc phương thức hợp lý khác.\n" +
				"Nếu người dùng tiếp tục mua, nạp tiền hoặc sử dụng AI Token sau khi thỏa thuận được cập nhật, coi như người dùng chấp nhận thỏa thuận đã cập nhật.\n" +
				"Nếu việc sửa đổi thỏa thuận gây ra ảnh hưởng bất lợi lớn đến quyền lợi cốt lõi của Token đã mua và chưa tiêu thụ của người dùng (giảm đáng kể số lượng Token, rút ngắn đáng kể thời hạn hiệu lực hoặc cấm hoàn toàn việc sử dụng mô hình đã cam kết ban đầu), nền tảng sẽ đưa ra thông báo hợp lý, và có thể lựa chọn cung cấp phương án thay thế có giá trị tương đương hoặc hoàn tiền theo tỷ lệ Token chưa sử dụng, trừ khi pháp luật bắt buộc có yêu cầu khác.",

			"Điều 17 Luật áp dụng và Giải quyết tranh chấp (phân biệt trong và ngoài nước)\n" +
				"Người dùng trong nước: Việc thiết lập, hiệu lực, thực hiện, giải thích và giải quyết tranh chấp của thỏa thuận này áp dụng luật pháp Cộng hòa Nhân dân Trung Hoa (không bao gồm quy tắc xung đột luật); tranh chấp phát sinh từ hoặc liên quan đến thỏa thuận này, hai bên cần giải quyết trước qua thương lượng thân thiện; thương lượng không thành, bất kỳ bên nào cũng phải nộp đơn lên Tòa án nhân dân có thẩm quyền tại nơi cư trú của bên cung cấp dịch vụ để giải quyết bằng kiện tụng.\n" +
				"Người dùng ngoài nước: Thỏa thuận này áp dụng luật pháp nơi cư trú của bên cung cấp dịch vụ (đồng thời tuân thủ luật pháp bắt buộc tại nơi cư trú của người dùng); tranh chấp phát sinh từ hoặc liên quan đến thỏa thuận này, hai bên cần giải quyết trước qua thương lượng thân thiện; thương lượng không thành, nộp lên Ủy ban Trọng tài Thương mại và Kinh tế Quốc tế Trung Quốc (CIETAC) để tiến hành trọng tài theo quy tắc trọng tài có hiệu lực của ủy ban tại thời điểm nộp đơn; địa điểm trọng tài là Thượng Hải, Trung Quốc, ngôn ngữ trọng tài là tiếng Trung / tiếng Anh, phán quyết trọng tài là chung thẩm, có hiệu lực ràng buộc đối với cả hai bên.\n" +
				"Trường hợp nơi cư trú của người dùng là người tiêu dùng có quy định bắt buộc về thẩm quyền xét xử thì thực hiện theo quy định đó.",

			"Điều 18 Thông báo\n" +
				"Nền tảng có thể gửi thông báo cho người dùng qua thông báo trang web, tin nhắn trong trang, thư điện tử, tin nhắn điện thoại, gọi điện thoại, thông báo trang điều khiển API, thông báo hóa đơn hoặc thông tin liên hệ do người dùng đăng ký.\n" +
				"Người dùng phải đảm bảo thông tin liên hệ của mình là xác thực, chính xác và có hiệu lực; nếu thông tin liên hệ thay đổi, người dùng phải cập nhật kịp thời.\n" +
				"Quy tắc gửi thông báo: Gửi thông báo trang web khi công bố; gửi tin nhắn trong trang, thư điện tử, tin nhắn điện thoại khi gửi thành công; gửi thư chuyển phát nhanh hoặc thư bảo đảm khi ký nhận, hoặc coi như đã gửi khi bị từ chối nhận hoặc không thể gửi do nguyên nhân từ người dùng.",

			"Điều 19 Điều khoản khác\n" +
				"Thỏa thuận này cấu thành thỏa thuận hoàn chỉnh giữa hai bên về việc mua và sử dụng AI Token, thay thế cho bất kỳ thỏa thuận nào đã đạt được trước đó bằng lời nói hoặc văn bản giữa hai bên.\n" +
				"Bất kỳ điều khoản nào của thỏa thuận này bị xác định là vô hiệu, bất hợp pháp hoặc không thể thực thi sẽ không ảnh hưởng đến hiệu lực của các điều khoản khác.\n" +
				"Nếu không có sự đồng ý bằng văn bản của nền tảng, người dùng không được chuyển nhượng quyền và nghĩa vụ theo thỏa thuận này.\n" +
				"Nền tảng có thể chuyển nhượng quyền và nghĩa vụ theo thỏa thuận này khi tái cơ cấu kinh doanh, sáp nhập, chia tách, chuyển nhượng tài sản, bên liên kết tiếp quản hoặc chuyển dịch dịch vụ, nhưng phải đảm bảo không làm giảm đáng kể quyền lợi hiện có của người dùng.\n" +
				"Tiêu đề của thỏa thuận này chỉ nhằm thuận tiện cho việc đọc, không ảnh hưởng đến việc giải thích các điều khoản.\n" +
				"Thỏa thuận này có thể cung cấp các phiên bản tiếng Trung, tiếng Anh hoặc ngôn ngữ khác; nếu có xung đột giữa các phiên bản ngôn ngữ khác nhau, phiên bản tiếng Trung sẽ được ưu tiên, trừ khi luật pháp bắt buộc tại nơi cư trú của người dùng có yêu cầu khác.\n" +
				"Người dùng nên tự tải xuống, lưu trữ hoặc in thỏa thuận này trước khi thanh toán; nền tảng cũng có thể cung cấp chức năng tra cứu phiên bản thỏa thuận, hồ sơ đơn hàng hoặc lịch sử giao dịch trong tài khoản người dùng.",
		},
	},
}
