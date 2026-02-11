package user

import "context"

type StaticPasswordVerifier struct {
	secrets map[string]string
}

func NewStaticPasswordVerifier(secrets map[string]string) StaticPasswordVerifier {
	copied := make(map[string]string, len(secrets))
	for id, secret := range secrets {
		copied[id] = secret
	}
	return StaticPasswordVerifier{secrets: copied}
}

func (v StaticPasswordVerifier) Verify(_ context.Context, id string, secret string) (bool, error) {
	expected, ok := v.secrets[id]
	if !ok {
		return false, nil
	}
	return expected == secret, nil
}
