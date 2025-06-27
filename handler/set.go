package handler

import (
	"github.com/PlayerNeo42/gvalkey/resp"
)

func (h *Handler) handleSet(args resp.Array) (resp.Payload, error) {
	parsedArgs, err := resp.ParseSetArgs(args)
	if err != nil {
		return nil, err
	}

	oldValue, success := h.store.Set(*parsedArgs)

	// handle the GET option: return the old value or NULL.
	if parsedArgs.Get {
		if !success || oldValue == nil {
			return resp.NULL, nil
		}

		// marshal the old value for the response safely.
		switch val := oldValue.(type) {
		case resp.Payload:
			return val, nil
		default:
			// this case prevents a panic if the store returns an unexpected type.
			return resp.NewSimpleError("internal error: stored value has an unmarshalable type"), nil
		}
	}

	// if a conditional SET (NX/XX) failed, return NULL.
	if !success {
		return resp.NULL, nil
	}

	// otherwise, the SET was successful.
	return resp.OK, nil
}
