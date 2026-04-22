package handler

func invitationResponseMessage(emailSent bool) string {
	if emailSent {
		return "Invitation email has been sent"
	}

	return "Invitation created, but email could not be sent"
}
