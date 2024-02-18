package auth

type UC interface {
	VerifySignature(service, signature, body, timestamp, requestId string) (bool, error)
}
