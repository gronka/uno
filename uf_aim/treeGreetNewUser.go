package uf_aim

func (aimInfo *AimInfo) MakeBranchGreetNewUserIntroductionResponse() {
	hi := ""
	if aimInfo.Surfer.Name == "" {
		hi = "Hi!"
	} else {
		hi = "Hi " + aimInfo.Surfer.Name + "!"
	}

	out := &aimInfo.MsgOut
	out.Code = "welcome_response"
	out.MediaUrl = "https://storage.googleapis.com/friday_odc/vcard.vcf"
	out.Content = hi +
		"I'm Fridayy - the simplest way to make quick purchases. " +
		"Visit https://fridayy.me to view our privacy policy and " +
		"terms of service. (tldr: We respect your data and we do not share it.)"
}

func (aimInfo *AimInfo) MakeBranchGreetNewUserWhatDoYouNeedResponse() {
	out := &aimInfo.MsgOut
	out.Code = "greet_what_do_you_need"
	out.MediaUrl = "https://storage.googleapis.com/friday_odc/vcard.vcf"
	out.Content = "Do you need to buy anything online today? I'll bet I can " +
		"find it right away!"
}
