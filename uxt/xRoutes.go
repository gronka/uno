package uxt

// uf_aim
var AimMessageReceiveV1 = "/v1/message.receive"
var AimChatGetBySurferIdV1 = "/v1/chat/get.bySurferId"
var AimChatGetByIdV1 = "/v1/chat/get.byId"
var AimChatGetNewMsgsV1 = "/v1/chat/get.newMsgs"
var AimChatsGetListV1 = "/v1/chats/get.list"
var AimChatUpsertV1 = "/api/v1/chat/upsert"

// uf_border
var BorderJsonHookReceiveV1 = "/v1/json/hook.receive"
var BorderLoopHookReceiveV1 = "/v1/loop/hook.receive"
var BorderSipHookReceiveV1 = "/v1/sip/hook.receive"
var BorderSendMsgOutV1 = "/v1/msgOut.send"

// uf_order
var OrderPaymentAttemptCreateV1 = "/v1/paymentAttempt.create"
var OrderOrderCancelV1 = "/v1/order.cancel"
var OrderOrderCreateV1 = "/v1/order.create"
var OrderOrderGetByIdDetailsV1 = "/v1/order/get.byId.details"
var OrderOrderGetBySurferIdNewestV1 = "/v1/order/get.bySurferId.newest"

var OrderOrdersGetBySurferIdV1 = "/v1/orders/get.bySurferId"

// uf_public
var PublicHealth = "/api"
var PublicToddSignInWithEmailV1 = "/api/v1/todd/signIn.withEmail"
var PublicToddSignInWithPhoneV1 = "/api/v1/todd/signIn.withPhone"
var PublicChatGetByPhoneV1 = "/api/v1/chat/get.byPhone"
var PublicChatGetBySurferIdV1 = "/api/v1/chat/get.bySurferId"
var PublicChatMsgsGetByIdV1 = "/api/v1/chatMsgs/get.byId"
var PublicChatSendMsgV1 = "/api/v1/chat/sendMsg"
var PublicChatsGetListV1 = "/api/v1/chats/get.list"

var PublicGlogsGetByUfIdV1 = "/api/v1/glogs/get.byUfId"

var PublicHookLoopReceive = "/api/hook/loop.receive"
var PublicHookSipReceive = "/api/hook/sip.receive"
var PublicHookWebV1 = "/api/v1/hook/web.receive"
var PublicHookWebCreateUserV1 = "/api/v1/hook/web.createUser"

var PublicSurferGetByIdV1 = "/api/v1/surfer/get.byId"

var PublicStripeSuccessHookV1 = "/api/v1/stripe/hook.success"

//var PublicStripeFailureHookV1 = "/api/v1/stripe/hook.failure"
var PublicStripeCancelHookV1 = "/api/v1/stripe/hook.cancel"

// uf_user
var UserSurferCreateStripeCustomerIdV1 = "/v1/surfer/create.stripeCustomerId"
var UserSurferGetOrCreateFromPhoneOrEmailV1 = "/v1/surfer/get.orCreateFromPhoneOrEmail"
var UserSurferGetByEmailV1 = "/v1/surfer/get.byEmail"
var UserSurferGetByPhoneV1 = "/v1/surfer/get.byPhone"
var UserSurferGetByIdV1 = "/v1/surfer/get.byId"
var UserSurferUpdateNameV1 = "/v1/surfer/update.name"

var UserAddressCreateBuilderV1 = "/v1/address/create.builder"
var UserAddressDeleteV1 = "/v1/address/delete"
var UserAddressDeleteBuilderV1 = "/v1/address/delete.builder"
var UserAddressGetByIdV1 = "/v1/address/get.byId"
var UserAddressGetBuilderBySurferIdV1 = "/v1/address/get.builderBySurferId"
var UserAddressGetNonBuilderBySurferIdV1 = "/v1/address/get.nonBuilderBySurferId"
var UserAddressUpdateNameV1 = "/v1/address/update.name"
var UserAddressUpdatePostalPlusV1 = "/v1/address/update.postalPlus"
var UserAddressUpdateLine1V1 = "/v1/address/update.line1"
var UserAddressUpdateLine2V1 = "/v1/address/update.line2"
var UserAddressValidateUspsV1 = "/v1/address/validate.usps"
