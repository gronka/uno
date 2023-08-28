package uf

func AddOurMargin(in int) int {
	og := float32(in)
	// try1 is our default margin case
	try1 := og * 1.1
	// try2 ensures that we do not lose money on cheap purchases
	//try2 := og*1.05 + 1

	//if try1 > try2 {
	//return int(try1)
	//} else {
	//return int(try2)
	//}
	return int(try1)
}
