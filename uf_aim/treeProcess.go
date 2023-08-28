package uf_aim

import (
	"strings"

	"gitlab.com/textfridayy/uno/uf"
	"gitlab.com/textfridayy/uno/uf_user"
	"gitlab.com/textfridayy/uno/ut"
	"gitlab.com/textfridayy/uno/uxt"
)

func (aimInfo *AimInfo) processTree(
	gibs *uf.Gibs,
) (ures uf.UfResponse, err error) {
	uf.Trace(aimInfo.RunningTreeName)
	uf.Trace(aimInfo.RunningBranch)
	switch aimInfo.RunningTreeName {
	case AimGreetNewUser:

		switch aimInfo.RunningBranch {
		case BranchStart:
			aimInfo.MakeBranchGreetNewUserIntroductionResponse()
			aimInfo.RunningBranch = BranchGreetNewUserWhatDoYouNeed

		case BranchGreetNewUserWhatDoYouNeed:
			aimInfo.MakeBranchGreetNewUserWhatDoYouNeedResponse()
			aimInfo.RunningBranch = BranchDone
			aimInfo.MustWaitUserResponse = true
		}

	case AimSearchItemVariant1:
		fallthrough
	case AimSearchItemVariant2:
		fallthrough
	case AimSearchItemVariant3:
		fallthrough
	case AimSearchItem:
		aimInfo.Cart.LCartGetCurrent(gibs, aimInfo.Surfer.SurferId)
		uf.Trace("88888888888888888888888888888888888888888")
		uf.Trace(aimInfo.Cart)

		switch aimInfo.RunningBranch {
		case BranchStart:
			// on BranchStart, we perform a new search
			product := aimInfo.PerformSearchInitial(gibs)

			aimInfo.MakeBranchSearchInitiateResponse(gibs, product)
			aimInfo.MustWaitUserResponse = true
			aimInfo.RunningBranch = BranchSearchRefine
			aimInfo.Query.LQueryUpsert(gibs)

		case BranchSearchRefine:
			// on BranchSearchRefine, we load the prior query
			uf.Trace("on refine search branch")
			aimInfo.LoadQueryFromUno(gibs)
			if aimInfo.IsShortcut() {
				cmd := aimInfo.GetMsgIn()
				switch cmd {
				case "1":
					fallthrough
				case "y":
					aimInfo.Cart = &CartPg{}
					uf.Trace("9999999999999999999999999999 CREATING CART")
					aimInfo.Cart.AddToCart(
						aimInfo.Query.CurrentProductId,
						int64(aimInfo.Query.Count),
						int64(uf.AddOurMargin(aimInfo.Query.CurrentProductPrice)),
					)
					uf.Trace(aimInfo.Cart)
					aimInfo.Cart.CartId = uf.RandomUUID()
					aimInfo.Cart.SurferId = aimInfo.Surfer.SurferId
					aimInfo.Cart.LCartUpsert(gibs)
					uf.Trace(aimInfo.Cart)

					aimInfo.MsgOut.Content = "refined search with " + cmd
					aimInfo.RunningBranch = BranchSearchSelectShippingAddress

				case "2":
					fallthrough
				case "c":
					aimInfo.Query.MinPrice = 0
					aimInfo.Query.MaxPrice = aimInfo.Query.CurrentProductPrice
					aimInfo.PerformSearchRefine(gibs)
					aimInfo.MustWaitUserResponse = true
					aimInfo.Query.LQueryUpsert(gibs)

				case "3":
					fallthrough
				case "q":
					aimInfo.Query.MinPrice = aimInfo.Query.CurrentProductPrice
					aimInfo.Query.MaxPrice = 0
					aimInfo.PerformSearchRefine(gibs)
					aimInfo.MustWaitUserResponse = true
					aimInfo.Query.LQueryUpsert(gibs)

				case "9":
					aimInfo.MsgOut.Content = "Search canceled. Do you want to search" +
						"for something else?"
					aimInfo.RunningBranch = BranchStart
					aimInfo.MustWaitUserResponse = true

				}
			} else {
				// reprocess the tree, but from from another Branch
				if aimInfo.MakoResponse.Function == AimSearchItem {
					aimInfo.RunningBranch = BranchStart
				}
			}

		case BranchSearchSelectShippingAddress:
			address, bures := uxt.AddressGetNonBuilderBySurferId(
				gibs,
				aimInfo.Surfer.SurferId)
			ures.AddErrors(bures.Errors)

			uf.Trace("looking for address")
			uf.Trace(address)
			mustCreateAddress := false
			if uf.UuidIsZero(address.AddressId) {
				mustCreateAddress = true
			} else {
				if aimInfo.IsShortcut() {
					cmd := aimInfo.GetMsgIn()
					switch cmd {
					case "1":
						aimInfo.Cart.ShippingAddressId = address.AddressId
						aimInfo.Cart.LCartUpsert(gibs)
						aimInfo.RunningBranch = BranchSearchSelectPayment
						aimInfo.MsgOut.Content = "Shipping address confirmed!"

					case "2":
						mustCreateAddress = true
					}

				} else {
					aimInfo.MsgOut.Content = "I have your delivery address as: " +
						address.AsLetterhead() + "\n" +
						"Is this correct?\n1) yes\n2) no"
					aimInfo.MustWaitUserResponse = true
				}
			}

			if mustCreateAddress {
				uxt.AddressCreateBuilder(gibs, aimInfo.Surfer.SurferId)
				aimInfo.RunningBranch = BranchSearchInputName
			}

		case BranchSearchInputName:
			quit := aimInfo.SetReplyIfUserWantsToQuitShipping(gibs)
			if !quit {
				name := aimInfo.GetMsgIn()
				if name == "" {
					aimInfo.MsgOut.Content = "Okay! I need to create a shipping label. " +
						"Can I get the name? You can say 'quit' at anytime."
					aimInfo.MustWaitUserResponse = true
				} else {
					names := strings.SplitN(name, " ", 2)
					if len(names) < 2 {
						aimInfo.MsgOut.Content = "Please enter your first name and last " +
							"name."
						aimInfo.MustWaitUserResponse = true
					} else {
						builder, _ := uxt.AddressGetBuilderBySurferId(gibs, aimInfo.Surfer.SurferId)
						uxt.AddressUpdateName(gibs, builder.AddressId, name)
						aimInfo.RunningBranch = BranchSearchInputPostal
					}
				}
			}

		case BranchSearchInputPostal:
			quit := aimInfo.SetReplyIfUserWantsToQuitShipping(gibs)
			if !quit {
				postal := aimInfo.GetMsgIn()
				if postal == "" {
					aimInfo.MsgOut.Content = "Next, I need your zip code."
					aimInfo.MustWaitUserResponse = true
				} else {
					lookup := uf_user.UspsZipcodeToCityState(gibs, postal)
					uf.Trace(postal)
					uf.Trace(lookup)

					if lookup.ZipC.City == "" {
						aimInfo.MsgOut.Content = "I wasn't able to find this zip code " +
							"can you send it again?"
						aimInfo.MustWaitUserResponse = true
					} else {
						builder, _ := uxt.AddressGetBuilderBySurferId(gibs, aimInfo.Surfer.SurferId)
						uxt.AddressUpdatePostalPlus(gibs, builder.AddressId, postal)
						aimInfo.RunningBranch = BranchSearchInputAddressLine1
					}
				}
			}

		case BranchSearchInputAddressLine1:
			quit := aimInfo.SetReplyIfUserWantsToQuitShipping(gibs)
			if !quit {
				line1 := aimInfo.GetMsgIn()
				if line1 == "" {
					aimInfo.MsgOut.Content = "What is the house number and street? " +
						"(108 Blue Hill Rd)"
					aimInfo.MustWaitUserResponse = true
				} else {
					builder, _ := uxt.AddressGetBuilderBySurferId(gibs, aimInfo.Surfer.SurferId)
					uxt.AddressUpdateLine1(gibs, builder.AddressId, line1)
					aimInfo.RunningBranch = BranchSearchInputAddressLine2
				}
			}

		case BranchSearchInputAddressLine2:
			quit := aimInfo.SetReplyIfUserWantsToQuitShipping(gibs)

			if !quit {
				line2 := aimInfo.GetMsgIn()
				if line2 == "" {
					aimInfo.MsgOut.Content = "Can you tell me the apartment or suite " +
						"number for the 2nd lind of your address? (such as Suite E or " +
						"Apartment 404). If you do not have one, say 'None'. If you " +
						"need to start over, say 'quit'."
					aimInfo.MustWaitUserResponse = true
				} else {
					builder, _ := uxt.AddressGetBuilderBySurferId(gibs, aimInfo.Surfer.SurferId)
					uxt.AddressUpdateLine2(gibs, builder.AddressId, line2)

					valid, _ := uxt.AddressValidateUsps(gibs, builder.AddressId)
					if valid.IsValid {
						aimInfo.MsgOut.Content = "Got it! Address verified!"
						aimInfo.RunningBranch = BranchSearchSelectShippingAddress
					} else {
						aimInfo.MsgOut.Content = "I failed to validate your address with " +
							"USPS. Try sending us the 2nd line again. If you need to start " +
							"over, say 'quit'. This is what you told me:\n" +
							builder.AsLetterhead()
						aimInfo.MustWaitUserResponse = true
					}
				}
			}

		//case BranchSearchSelectShippingMethod:
		//quit := aimInfo.SetReplyIfUserWantsToQuitShipping(gibs)

		//if !quit {
		//choice := aimInfo.GetMsgIn()

		//if choice == "" {
		//aimInfo.Cart.GetShippingChoices(gibs)
		//aimInfo.Cart.LCartUpsert(gibs)
		//aimInfo.MsgOut.Content = aimInfo.Cart.ShowShippingChoicesText()
		//aimInfo.MustWaitUserResponse = true
		//} else {
		//switch choice {
		//case "1":
		//aimInfo.Cart.ShippingChoice = "cheapest"
		//case "2":
		//aimInfo.Cart.ShippingChoice = "fastest"
		//case "3":
		//aimInfo.Cart.ShippingChoice = "balanced"
		//}
		//aimInfo.Cart.LCartUpsert(gibs)
		//aimInfo.RunningBranch = BranchSearchSelectPayment
		//}
		//}

		case BranchSearchSelectPayment:
			quit := aimInfo.SetReplyIfUserWantsToQuitShipping(gibs)
			if !quit {

				appendLink := false
				uf.Trace("55555555555555555555555555555555555555555")
				uf.Trace(aimInfo.Cart.ProductIds)
				uf.Trace("3")
				uf.Trace("selecting payment")

				basketItems := []ut.BasketItemPg{}
				product := ProductPg{}
				for index, productId := range aimInfo.Cart.ProductIds {
					product, _ = LProductGetById(gibs, productId)
					tempItem := ut.BasketItemPg{
						ProductId:             productId,
						Title:                 product.Title,
						Quantity:              int(aimInfo.Cart.Counts[index]),
						UnitPrice:             product.Price,
						Currency:              product.Currency,
						ScrapeEngine:          product.ScrapeEngine,
						ScrapeEngineProductId: product.ScrapeEngineProductId,
						StoreUrl:              product.StoreUrl,
						ProductUrl:            product.ProductUrl,
						ImageUrl:              product.ImageUrl,
					}
					basketItems = append(basketItems, tempItem)
				}

				productId := aimInfo.Cart.ProductIds[0]
				product, _ = LProductGetById(gibs, productId)

				orderPkg := ut.OrderCreateIn{
					SurferId:          aimInfo.Surfer.SurferId,
					Title:             product.Title,
					Price:             product.Price,
					Margin:            product.GetMargin(),
					Currency:          "USD",
					ShippingDays:      product.ShippingDays,
					ShippingPrice:     product.ShippingPrice,
					ShippingAddressId: aimInfo.Cart.ShippingAddressId,
					BasketItems:       basketItems,
				}

				orderResponse, _ := uxt.OrderCreate(gibs, orderPkg)

				// prepare for scoping
				paResponse := ut.PaymentAttemptCreateOut{}
				if orderResponse.PaymentMethodId != "" {
					uf.Trace("is payment expired?")
					if orderResponse.IsPaymentMethodExpired {
						uf.Trace("payment IS expired?")
						aimInfo.MsgOut.Content = "Your default card has expired. Please " +
							"follow the link to use another card."
						appendLink = true
						aimInfo.MustWaitUserResponse = true
					} else {
						uf.Trace("payment IS NOT expired?")
						if aimInfo.IsShortcut() {
							uf.Trace("IS shortcut")
							cmd := aimInfo.GetMsgIn()

							switch cmd {
							case "1":
								challengeWord := ut.GetChallengeWord()
								aimInfo.MsgOut.Content = "To verify you are a human, I need " +
									"you to send me a text that says '" + challengeWord + "'. " +
									"You have 3 attempts. To cancel, write 'quit'."
								aimInfo.ChallengeWord = challengeWord
								aimInfo.ChallengeWordCounter = 0
								aimInfo.RunningBranch = BranchSearchChallengeWord

							case "2":
								aimInfo.MsgOut.Content = "Please follow the link to pay " +
									"with your card."
								aimInfo.RunningBranch = BranchDone
								appendLink = true
								aimInfo.MustWaitUserResponse = true

							}
						} else {
							uf.Trace("IS NOT shortcut")
							aimInfo.MsgOut.Content = "Do you want to charge to your card " +
								"ending in " + orderResponse.CardLast4 + "?\n1) yes\n2) no"
							aimInfo.MustWaitUserResponse = true
						}
					}
				} else {
					aimInfo.MsgOut.Content = "Almost done! Finally, I need your payment" +
						"information. Stripe (tm) handles payments for us (and millions of " +
						" others). Please click the link below to finalize your purchase."
					aimInfo.RunningBranch = BranchDone
					appendLink = true
					aimInfo.MustWaitUserResponse = true

				}

				if appendLink {
					paResponse, _ = uxt.PaymentAttemptCreate(
						gibs,
						ut.PaymentAttemptCreateIn{
							OrderId:      orderResponse.OrderId,
							SurferId:     aimInfo.Surfer.SurferId,
							UseSavedCard: false,
						})

					uf.Trace("appending link")
					uf.Trace(uf.IntToPriceString(aimInfo.Cart.GetFinalPrice(gibs)))
					uf.Trace(aimInfo.Query.AsQueryString())
					uf.Trace(paResponse.PaymentUrl)

					if err != nil {
						uf.Glog(gibs, uf.GlogStruct{
							Level:     uf.LevelError,
							Code:      "aim.101",
							Interface: err,
						})
					}
					aimInfo.MsgOut.Content += "\n\n" + paResponse.PaymentUrl
				}

			}
		case BranchSearchChallengeWord:
			order, _ := uxt.OrderGetBySurferIdNewest(gibs, aimInfo.Surfer.SurferId)
			quit := aimInfo.SetReplyIfUserWantsToQuitShipping(gibs)
			if !quit {
				word := aimInfo.GetMsgIn()

				if aimInfo.ChallengeWord != word {
					uf.Trace("Challenge word incorrect!")
					aimInfo.MsgOut.Content = "That's incorrect - please try again."
					aimInfo.ChallengeWordCounter++

					if aimInfo.ChallengeWordCounter == 3 {
						aimInfo.MsgOut.Content = "I'm sorry - you have failed the word " +
							"challenge too many times. For your account's security, I " +
							"will cancel your order."
						aimInfo.RunningBranch = BranchDone
					}

				} else {
					uf.Trace("Challenge word correct!")
					paResponse, _ := uxt.PaymentAttemptCreate(
						gibs,
						ut.PaymentAttemptCreateIn{
							OrderId:      order.OrderId,
							SurferId:     aimInfo.Surfer.SurferId,
							UseSavedCard: true,
						})

					if paResponse.PaymentSucceeded {
						aimInfo.MsgOut.Content = "Okay! Your order has been " +
							"placed! I will send you tracking information as soon " +
							"as I have it."
						aimInfo.RunningBranch = BranchDone
					} else {
						aimInfo.MsgOut.Content = "We were unable to charge your card. " +
							"Please follow the link below to complete payment in a web " +
							"browser."
						paResponse, _ = uxt.PaymentAttemptCreate(
							gibs,
							ut.PaymentAttemptCreateIn{
								OrderId:      order.OrderId,
								SurferId:     aimInfo.Surfer.SurferId,
								UseSavedCard: false,
							})
						aimInfo.MsgOut.Content += "\n\n" + paResponse.PaymentUrl
						aimInfo.RunningBranch = BranchDone
					}
				}
			}

		} // end of switch RunningBranch

	case AimCancelOrder:
		uf.Trace(AimCancelOrder)
		order, _ := uxt.OrderGetBySurferIdNewest(gibs, aimInfo.Surfer.SurferId)
		if uf.UuidIsZero(order.OrderId) {
			aimInfo.MsgOut.Content = "I was going to help you cancel an order, " +
				"but you haven't bought anything yet!"
		} else {
			if !aimInfo.IsShortcut() {
				uf.Trace("IS NOT shortcut")
				aimInfo.MsgOut.Content = "I looked up your most recent order. You " +
					"purchased " + order.Title + ". What would you like me to do? " +
					"\n1) cancel this order\n2) nothing"
				aimInfo.MustWaitUserResponse = true

			} else {
				uf.Trace("IS shortcut")
				cmd := aimInfo.GetMsgIn()
				switch cmd {
				case "1":
					aimInfo.MsgOut.Content = "Okay, I am going to try and cancel this " +
						"order. Not all orders can be canceled. I will emssage you with " +
						"any updates."
					aimInfo.RunningBranch = BranchDone

				case "2":
					aimInfo.MsgOut.Content = "Okay, I will not modify your order."
					aimInfo.RunningBranch = BranchDone
				default:
					aimInfo.MsgOut.Content = "I didn't catch that. Please reply with " +
						"\n1) to cancel this order\n2) to do nothing"
				}
			}
		}

	case AimTrackOrder:
		uf.Trace(AimTrackOrder)
		order, _ := uxt.OrderGetBySurferIdNewest(gibs, aimInfo.Surfer.SurferId)
		if uf.UuidIsZero(order.OrderId) {
			aimInfo.MsgOut.Content = "I tried to look up tracking information " +
				"for you, but you haven't bought anything yet!"
		} else {
			details, _ := uxt.OrderGetDetails(gibs, order.OrderId)

			if len(details.Tracking) > 0 {
				carrier := details.Tracking[0].Carrier
				num := details.Tracking[0].TrackingNumber
				aimInfo.MsgOut.Content = "I looked up your most recent order. You " +
					"purchased " + order.Title + ". It is " +
					"being shipped to you by " + carrier + ". The tracking number is\n" +
					num
			} else {
				url := "http://makeAUrl"
				aimInfo.MsgOut.Content = "I could not find any tracking information " +
					"for your most recent purchase. You can view your recent purchases " +
					"online at\n" + url
				//TODO: online URL
			}
		}
		aimInfo.RunningBranch = BranchDone

	case AimReorder:
		uf.Trace(AimReorder)
		order, _ := uxt.OrderGetBySurferIdNewest(gibs, aimInfo.Surfer.SurferId)
		if uf.UuidIsZero(order.OrderId) {
			aimInfo.MsgOut.Content = "I was going to help you place an order " +
				"again, but you haven't bought anything yet!"
		} else {
			if !aimInfo.IsShortcut() {
				uf.Trace("IS NOT shortcut")
				aimInfo.MsgOut.Content = "I looked up your most recent order. You " +
					"purchased " + order.Title + ". What would you like me to do? " +
					"\n1) order this again\n2) nothing"
				aimInfo.MustWaitUserResponse = true

			} else {
				uf.Trace("IS shortcut")
				cmd := aimInfo.GetMsgIn()
				switch cmd {
				case "1":
					aimInfo.MsgOut.Content = "Okay, I am going to try and cancel this " +
						"order. Not all orders can be canceled. I will emssage you with " +
						"any updates."
					aimInfo.RunningBranch = BranchDone

				case "2":
					aimInfo.MsgOut.Content = "Okay, I will not modify your order."
					aimInfo.RunningBranch = BranchDone
				default:
					aimInfo.MsgOut.Content = "I didn't catch that. Please reply with " +
						"\n1) to cancel this order\n2) to do nothing"
				}
			}
		}

	case AimNone:
		fallthrough
	case AimBlank:
		//NOTE: increment to force an exit
		//aimInfo.RunningBranch++

	case AimDoSomething:
		fallthrough
	case AimAskSomething:
		fallthrough
	case AimDescribeSomething:
		//NOTE: shopping case
		fallthrough
	default:
		aimInfo.RunningBranch = BranchDone
		uf.Error("Undeveloped or unknown AimSeries case: " +
			string(aimInfo.RunningTreeName))
		aimInfo.RunningTreeName = AimBlank
		//aimInfo.MsgOut.Content = "unknown aim for aim.Content: " +
		//aimInfo.MsgIn.Content
	}
	return ures, err
}
