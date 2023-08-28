# uf_aim: AI Management

uf_aim is used to "aim" incoming requests. It will be Fridayy's 1st layer of 
filtering of requests.

Long-term, we envision Fridayy holding multiple conversations/contexts at once.
To make this happen, we need a top-layer filter which will determine where we
should send each incoming user request. Once the aim is determined, uf_aim will
send the request to the appropriate microservice.
