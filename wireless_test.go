package gotwilio

import "testing"

var simId = ""
var ratePlanId = ""
var accountSid = ""
var token = ""

func TestGetSimCardByIccId(t *testing.T) {
	twilio := NewTwilioClient(accountSid, token)
	card, ex, err := twilio.GetSimCardByIccId(simId)
	if err != nil {
		t.Fatal(err)
	}

	if ex != nil {
		t.Fatal(ex.Message)
	}

	if card.IccId != simId {
		t.Fatal("Sim card id does not match")
	}
}

func TestUpdateSimCardStatus(t *testing.T) {
	twilio := NewTwilioClient(accountSid, token)
	card, ex, err := twilio.GetSimCardByIccId(simId)
	if err != nil {
		t.Fatal(err)
	}

	if ex != nil {
		t.Fatal(ex.Message)
	}

	if card.IccId != simId {
		t.Fatal("Sim card id does not match")
	}

	if card.Status != SimStatus_Active && card.Status != SimStatus_Ready {
		ex, err = twilio.UpdateSimStatus(SimStatus_Ready, card.SId)
		if err != nil {
			t.Fatal(err)
		}

		if ex != nil {
			t.Fatal(ex.Message)
		}
	}
}

func TestUpdateCardRatePlan(t *testing.T) {
	twilio := NewTwilioClient(accountSid, token)
	card, ex, err := twilio.GetSimCardByIccId(simId)
	if err != nil {
		t.Fatal(err)
	}

	if ex != nil {
		t.Fatal(ex.Message)
	}

	if card.IccId != simId {
		t.Fatal("Sim card id does not match")
	}

	if card.RatePlanSId != ratePlanId {
		ex, err = twilio.UpdateSimRatePlan(ratePlanId, card.SId)
		if err != nil {
			t.Fatal(err)
		}

		if ex != nil {
			t.Fatal(ex.Message)
		}
	}
}

func TestGetRatePlans(t *testing.T) {
	twilio := NewTwilioClient("AC3d5310c77122fa60cb646fdb7b99202a", "4039f760773fa5bdd7cb75f4a6dfa3d8")
	plans, ex, err := twilio.GetAllRatePlans()
	if err != nil {
		t.Fatal(err)
	}

	if ex != nil {
		t.Fatal(ex.Message)
	}

	if len(plans) == 0 {
		t.Fatal("No rate plans found")
	}
}
