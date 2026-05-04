package errmsg

const (
	MessageYourRequestHasFailedToProcess                                         = "your_request_has_failed_to_process"
	MessageYourRequestHasBeenFailedToProcess                                     = "your_request_has_been_failed_to_process"
	MessageYourRequestHasBeenSuccessfullyProcessed                               = "your_request_has_been_successfully_processed"
	MessageCompanyNotFound                                                       = "company_not_found"
	MessageCompanyCodeAlreadyExists                                              = "company_code_already_exists"
	MessageAttendanceRadiusMustBeGte0                                            = "attendance_radius_must_be_gte_0"
	MessageTimezoneIsRequired                                                    = "timezone_is_required"
	MessageDateFormatIsRequired                                                  = "date_format_is_required"
	MessageTimeFormatIsRequired                                                  = "time_format_is_required"
	MessagePreferredLanguageIsRequired                                           = "preferred_language_is_required"
	MessageSupportedLanguagesIsRequired                                          = "supported_languages_is_required"
	MessageSupportedLanguagesContainsUnsupportedLanguage                         = "supported_languages_contains_unsupported_language"
	MessagePreferredLanguageMustExistInSupportedLanguages                        = "preferred_language_must_exist_in_supported_languages"
	MessageInvalidCredentials                                                    = "invalid_credentials"
	MessageFailedToGenerateToken                                                 = "failed_to_generate_token"
	MessageEmailIsAlreadyRegistered                                              = "email_is_already_registered"
	MessageTenantCodeIsAlreadyRegistered                                         = "tenant_code_is_already_registered"
	MessageTenantCodeIsInvalid                                                   = "tenant_code_is_invalid"
	MessageTenantCodeIsReserved                                                  = "tenant_code_is_reserved"
	MessageTenantSetupConfirmationIsRequired                                     = "tenant_setup_confirmation_is_required"
	MessageInvalidGoogleToken                                                    = "invalid_google_token"
	MessageGoogleRegistrationSuccessful                                          = "google_registration_successful"
	MessageGoogleLoginSuccessful                                                 = "google_login_successful"
	MessageGoogleAccountIsNotConnectedToAPresenseUser                            = "google_account_is_not_connected_to_a_presense_user"
	MessageThisGoogleEmailExistsInMultipleTenantsSignInWithEmailAndPasswordFirst = "this_google_email_exists_in_multiple_tenants_sign_in_with_email_and_password_first"
	MessageGoogleEmailIsAlreadyRegistered                                        = "google_email_is_already_registered"
	MessageGoogleIdentityIsAlreadyLinked                                         = "google_identity_is_already_linked"
	MessageTenantNotFound                                                        = "tenant_not_found"
	MessageRoleNotFound                                                          = "role_not_found"
	MessageUserNotFound                                                          = "user_not_found"
	MessageUserHasNoRoles                                                        = "user_has_no_roles"
	MessageInvalidOrExpiredRefreshToken                                          = "invalid_or_expired_refresh_token"
	MessagePleaseVerifyYourEmailToContinue                                       = "please_verify_your_email_to_continue"
	MessageEmailIsNotVerified                                                    = "email_is_not_verified"
	MessageEmailIsAlreadyVerified                                                = "email_is_already_verified"
	MessageEmailVerifiedSuccessfully                                             = "email_verified_successfully"
	MessageVerificationEmailHasBeenResent                                        = "verification_email_has_been_resent"
	MessageIfTheEmailIsRegisteredAPasswordResetLinkHasBeenSent                   = "if_the_email_is_registered_a_password_reset_link_has_been_sent"
	MessagePasswordHasBeenResetSuccessfully                                      = "password_has_been_reset_successfully"
	MessageInvalidOrExpiredEmailVerificationToken                                = "invalid_or_expired_email_verification_token"
	MessageInvalidOrExpiredPasswordResetToken                                    = "invalid_or_expired_password_reset_token"
	MessageTooManyVerificationEmailRequestsPleaseWaitBeforeTryingAgain           = "too_many_verification_email_requests_please_wait_before_trying_again"
	MessageTooManyPasswordResetRequestsPleaseWaitBeforeTryingAgain               = "too_many_password_reset_requests_please_wait_before_trying_again"
	MessageEmployeeNotFound                                                      = "employee_not_found"
	MessageWorkLocationNotFound                                                  = "work_location_not_found"
	MessageWorkShiftNotFound                                                     = "work_shift_not_found"
	MessageOrgUnitNotFound                                                       = "org_unit_not_found"
	MessageOrganizationUnitNotFound                                              = "organization_unit_not_found"
	MessageYouCannotDeleteYourOwnAccount                                         = "you_cannot_delete_your_own_account"
	MessageFileCanNoLongerBeDeleted                                              = "file_can_no_longer_be_deleted"
	MessageInvalidJoinDateOrJoinDateTimezone                                     = "invalid_join_date_or_join_date_timezone"
	MessageEmployeeNumberAlreadyExistsInThisTenant                               = "employee_number_already_exists_in_this_tenant"
	MessageWorkShiftIsRequired                                                   = "work_shift_is_required"
	MessageFailedToGetEmployeeRole                                               = "failed_to_get_employee_role"
	MessageFailedToCreateUserAccount                                             = "failed_to_create_user_account"
	MessageFailedToUpdateEmployeePassword                                        = "failed_to_update_employee_password"
	MessageFailedToCreateUser                                                    = "failed_to_create_user"
	MessageFailedToUpdateUser                                                    = "failed_to_update_user"
	MessageWorkLocationNameAlreadyExists                                         = "work_location_name_already_exists"
	MessageAttendanceLogNotFound                                                 = "attendance_log_not_found"
	MessageAttendanceRecordedSuccessfully                                        = "attendance_recorded_successfully"
	MessageCheckInAlreadyRecordedForToday                                        = "check_in_already_recorded_for_today"
	MessageCheckOutAlreadyRecordedForToday                                       = "check_out_already_recorded_for_today"
	MessageCheckInMustBeRecordedBeforeCheckOut                                   = "check_in_must_be_recorded_before_check_out"
	MessageCompanyPolicyNotFound                                                 = "company_policy_not_found"
	MessageEmployeeProfileNotFound                                               = "employee_profile_not_found"
	MessageExportJobMessageBusIsNotConfigured                                    = "export_job_message_bus_is_not_configured"
	MessageEmailMessageBusIsNotConfigured                                        = "email_message_bus_is_not_configured"
	MessageExportJobNotFound                                                     = "export_job_not_found"
	MessageExportJobResourceMetadataIsRequired                                   = "export_job_resource_metadata_is_required"
	MessageFileFilterIsRequired                                                  = "file_filter_is_required"
	MessageFileNotFound                                                          = "file_not_found"
	MessageInvalidCategoryForTheSelectedParent                                   = "invalid_category_for_the_selected_parent"
	MessageInvalidLoggedAtFormat                                                 = "invalid_logged_at_format"
	MessageInvalidSelfieFile                                                     = "invalid_selfie_file"
	MessageJobPositionNotFound                                                   = "job_position_not_found"
	MessageOrganizationUnitCodeAlreadyExists                                     = "organization_unit_code_already_exists"
	MessageRouteNotFound                                                         = "route_not_found"
	MessageSelfieFileIsNoLongerAvailable                                         = "selfie_file_is_no_longer_available"
	MessageSelfieFileNotFound                                                    = "selfie_file_not_found"
	MessageUserRoleNotFound                                                      = "user_role_not_found"
	MessageYouAreNotAllowedToAccessExportJobs                                    = "you_are_not_allowed_to_access_export_jobs"
	MessageYouAreNotAllowedToAccessThisResource                                  = "you_are_not_allowed_to_access_this_resource"
	MessageYouAreNotAllowedToExportAttendanceReports                             = "you_are_not_allowed_to_export_attendance_reports"
	MessageYouDontHaveAccessToThisCompany                                        = "you_dont_have_access_to_this_company"
	MessageProductNotFound                                                       = "product_not_found"
	MessageFileIsRequired                                                        = "file_is_required"
	MessageBodyIsRequired                                                        = "body_is_required"
	MessageFilenameIsRequired                                                    = "filename_is_required"
	MessageFailedToGetPageOfResults                                              = "failed_to_get_page_of_results"
	MessageAcceptedInvitationCannotBeResent                                      = "accepted_invitation_cannot_be_resent"
	MessageAcceptedInvitationCannotBeRevoked                                     = "accepted_invitation_cannot_be_revoked"
	MessageActiveProductPricingNotFound                                          = "active_product_pricing_not_found"
	MessageAmountMustBePositive                                                  = "amount_must_be_positive"
	MessageAmountsMustNotBeNegative                                              = "amounts_must_not_be_negative"
	MessageAnActiveInvitationAlreadyExists                                       = "an_active_invitation_already_exists"
	MessageClientOwnerRoleNotFound                                               = "client_owner_role_not_found"
	MessageCompanyCanOnlyBeUpdatedFromCompanyDetails                             = "company_can_only_be_updated_from_company_details"
	MessageCompanyCannotBeDeleted                                                = "company_cannot_be_deleted"
	MessageCurrencyCodeFormatIsInvalid                                           = "currency_code_format_is_invalid"
	MessageCurrencyCodeAlreadyExists                                             = "currency_code_already_exists"
	MessageCurrencyIsNotActive                                                   = "currency_is_not_active"
	MessageCurrencyNotFound                                                      = "currency_not_found"
	MessageProviderCurrencyIsNotActive                                           = "provider_currency_is_not_active"
	MessagePriceIsNotActive                                                      = "price_is_not_active"
	MessageDefaultCurrencyCannotBeInactive                                       = "default_currency_cannot_be_inactive"
	MessageDefaultCurrencyCannotBeDeleted                                        = "default_currency_cannot_be_deleted"
	MessageCurrencyIsStillUsed                                                   = "currency_is_still_used"
	MessageDokuClientIsNotConfigured                                             = "doku_client_is_not_configured"
	MessageEmployeeIdIsNotAllowedForAdminInvitations                             = "employee_id_is_not_allowed_for_admin_invitations"
	MessageEmployeeIdIsRequiredForEmployeeInvitations                            = "employee_id_is_required_for_employee_invitations"
	MessageEmployeeAlreadyHasAnActiveUser                                        = "employee_already_has_an_active_user"
	MessageFailedToAcceptInvitation                                              = "failed_to_accept_invitation"
	MessageFailedToCreateInvitation                                              = "failed_to_create_invitation"
	MessageFailedToHashPassword                                                  = "failed_to_hash_password"
	MessageFailedToResendInvitation                                              = "failed_to_resend_invitation"
	MessageFullNameIsRequired                                                    = "full_name_is_required"
	MessageGoogleRegistrationSessionCreated                                      = "google_registration_session_created"
	MessageGoogleRegistrationSessionIsInvalidOrExpired                           = "google_registration_session_is_invalid_or_expired"
	MessageInternalProductCodeAlreadyExists                                      = "internal_product_code_already_exists"
	MessageInternalProductCodeFormatIsInvalid                                    = "internal_product_code_format_is_invalid"
	MessageInternalProductNotFound                                               = "internal_product_not_found"
	MessageInternalProductPriceNotFound                                          = "internal_product_price_not_found"
	MessageInternalProductPricingCodeAlreadyExists                               = "internal_product_pricing_code_already_exists"
	MessageInternalProductPricingCodeFormatIsInvalid                             = "internal_product_pricing_code_format_is_invalid"
	MessageInternalProductPricingNotFound                                        = "internal_product_pricing_not_found"
	MessageInternalProductPricingStatusIsInvalid                                 = "internal_product_pricing_status_is_invalid"
	MessageInternalProductStatusIsInvalid                                        = "internal_product_status_is_invalid"
	MessageInternalUserEmailAlreadyExists                                        = "internal_user_email_already_exists"
	MessageInternalUserIsNotActive                                               = "internal_user_is_not_active"
	MessageInternalUserNotFound                                                  = "internal_user_not_found"
	MessageInternalUserRoleIsInvalid                                             = "internal_user_role_is_invalid"
	MessageInternalUserStatusIsInvalid                                           = "internal_user_status_is_invalid"
	MessageInvalidInvitationRole                                                 = "invalid_invitation_role"
	MessageInvalidRefreshToken                                                   = "invalid_refresh_token"
	MessageInvitationEmailBusIsNotConfigured                                     = "invitation_email_bus_is_not_configured"
	MessageInvitationEmailMustMatchEmployeeEmail                                 = "invitation_email_must_match_employee_email"
	MessageInvitationHasAlreadyBeenAccepted                                      = "invitation_has_already_been_accepted"
	MessageInvitationHasBeenRevoked                                              = "invitation_has_been_revoked"
	MessageInvitationHasExpired                                                  = "invitation_has_expired"
	MessageInvitationIsAlreadyRevoked                                            = "invitation_is_already_revoked"
	MessageInvitationNotFound                                                    = "invitation_not_found"
	MessageInvoiceIsNotPayable                                                   = "invoice_is_not_payable"
	MessageInvoiceNotFound                                                       = "invoice_not_found"
	MessageOneOrMorePermissionsAreInvalidForThisClient                           = "one_or_more_permissions_are_invalid_for_this_client"
	MessageOnlyApprovedQuotationCanBeConverted                                   = "only_approved_quotation_can_be_converted"
	MessageOnlyOwnerCanInviteAdminUsers                                          = "only_owner_can_invite_admin_users"
	MessageOnlyOwnerCanManageAdminInvitations                                    = "only_owner_can_manage_admin_invitations"
	MessageOrderNotFound                                                         = "order_not_found"
	MessagePasswordIsRequired                                                    = "password_is_required"
	MessagePaymentAmountDoesNotMatchInvoiceOutstandingAmount                     = "payment_amount_does_not_match_invoice_outstanding_amount"
	MessagePaymentAttemptCannotBeRetried                                         = "payment_attempt_cannot_be_retried"
	MessagePaymentAttemptNotFound                                                = "payment_attempt_not_found"
	MessagePricePeriodOverlapsWithAnExistingActivePrice                          = "price_period_overlaps_with_an_existing_active_price"
	MessageQuotationCannotBeApprovedFromCurrentStatus                            = "quotation_cannot_be_approved_from_current_status"
	MessageQuotationNotFound                                                     = "quotation_not_found"
	MessageRevokedInvitationCannotBeResent                                       = "revoked_invitation_cannot_be_resent"
	MessageStartedAtMustUseRfc3339Format                                         = "started_at_must_use_rfc3339_format"
	MessageEndedAtMustUseRfc3339Format                                           = "ended_at_must_use_rfc3339_format"
	MessageEndedAtMustBeLaterThanStartedAt                                       = "ended_at_must_be_later_than_started_at"
	MessageTenantIsRequired                                                      = "tenant_is_required"
	MessageUnsupportedInvoiceAction                                              = "unsupported_invoice_action"
	MessageUnsupportedPaymentAttemptAction                                       = "unsupported_payment_attempt_action"
	MessageUnsupportedQuotationAction                                            = "unsupported_quotation_action"
	MessageUserActionTokenNotFound                                               = "user_action_token_not_found"
	MessageWebhookProcessed                                                      = "webhook_processed"
	MessageYouAreNotAuthorizedToInviteUsers                                      = "you_are_not_authorized_to_invite_users"
	MessageYouAreNotAuthorizedToManageInternalClientPermissions                  = "you_are_not_authorized_to_manage_internal_client_permissions"
	MessageYouAreNotAuthorizedToManageInternalProducts                           = "you_are_not_authorized_to_manage_internal_products"
	MessageYouAreNotAuthorizedToManageInternalUsers                              = "you_are_not_authorized_to_manage_internal_users"
	MessageYouAreNotAuthorizedToManageInvitations                                = "you_are_not_authorized_to_manage_invitations"
	MessageYouAreNotAuthorizedToViewInternalClients                              = "you_are_not_authorized_to_view_internal_clients"
	MessageYouAreNotAuthorizedToViewInvitations                                  = "you_are_not_authorized_to_view_invitations"
	MessageValidationDefault                                                     = "validation.default"
	MessageValidationDefaultWithParam                                            = "validation.default_with_param"
	MessageValidationRequired                                                    = "validation.required"
	MessageValidationEmail                                                       = "validation.email"
	MessageValidationEmailBlacklist                                              = "validation.email_blacklist"
	MessageValidationStrongPassword                                              = "validation.strong_password"
	MessageValidationResourceNotExist                                            = "validation.resource_not_exist"
	MessageValidationDatetime                                                    = "validation.datetime"
	MessageValidationSimpleInvalid                                               = "validation.simple_invalid"
	MessageValidationSimpleFormat                                                = "validation.simple_format"
	MessageValidationMinValue                                                    = "validation.min_value"
	MessageValidationMinChars                                                    = "validation.min_chars"
	MessageValidationMinItems                                                    = "validation.min_items"
	MessageValidationMaxValue                                                    = "validation.max_value"
	MessageValidationMaxChars                                                    = "validation.max_chars"
	MessageValidationMaxItems                                                    = "validation.max_items"
	MessageValidationCompareGt                                                   = "validation.compare_gt"
	MessageValidationCompareGte                                                  = "validation.compare_gte"
	MessageValidationCompareLt                                                   = "validation.compare_lt"
	MessageValidationCompareLte                                                  = "validation.compare_lte"
	MessageValidationNumeric                                                     = "validation.numeric"
	MessageValidationTimezone                                                    = "validation.timezone"
	MessageValidationEqfield                                                     = "validation.eqfield"
	MessageValidationOneof                                                       = "validation.oneof"
	MessageValidationUniqueInSlice                                               = "validation.unique_in_slice"
)

var AllMessageCodes = []string{
	MessageYourRequestHasFailedToProcess,
	MessageYourRequestHasBeenFailedToProcess,
	MessageYourRequestHasBeenSuccessfullyProcessed,
	MessageCompanyNotFound,
	MessageCompanyCodeAlreadyExists,
	MessageAttendanceRadiusMustBeGte0,
	MessageTimezoneIsRequired,
	MessageDateFormatIsRequired,
	MessageTimeFormatIsRequired,
	MessagePreferredLanguageIsRequired,
	MessageSupportedLanguagesIsRequired,
	MessageSupportedLanguagesContainsUnsupportedLanguage,
	MessagePreferredLanguageMustExistInSupportedLanguages,
	MessageInvalidCredentials,
	MessageFailedToGenerateToken,
	MessageEmailIsAlreadyRegistered,
	MessageTenantCodeIsAlreadyRegistered,
	MessageTenantCodeIsInvalid,
	MessageTenantCodeIsReserved,
	MessageTenantSetupConfirmationIsRequired,
	MessageInvalidGoogleToken,
	MessageGoogleRegistrationSuccessful,
	MessageGoogleLoginSuccessful,
	MessageGoogleAccountIsNotConnectedToAPresenseUser,
	MessageThisGoogleEmailExistsInMultipleTenantsSignInWithEmailAndPasswordFirst,
	MessageGoogleEmailIsAlreadyRegistered,
	MessageGoogleIdentityIsAlreadyLinked,
	MessageTenantNotFound,
	MessageRoleNotFound,
	MessageUserNotFound,
	MessageUserHasNoRoles,
	MessageInvalidOrExpiredRefreshToken,
	MessagePleaseVerifyYourEmailToContinue,
	MessageEmailIsNotVerified,
	MessageEmailIsAlreadyVerified,
	MessageEmailVerifiedSuccessfully,
	MessageVerificationEmailHasBeenResent,
	MessageIfTheEmailIsRegisteredAPasswordResetLinkHasBeenSent,
	MessagePasswordHasBeenResetSuccessfully,
	MessageInvalidOrExpiredEmailVerificationToken,
	MessageInvalidOrExpiredPasswordResetToken,
	MessageTooManyVerificationEmailRequestsPleaseWaitBeforeTryingAgain,
	MessageTooManyPasswordResetRequestsPleaseWaitBeforeTryingAgain,
	MessageEmployeeNotFound,
	MessageWorkLocationNotFound,
	MessageWorkShiftNotFound,
	MessageOrgUnitNotFound,
	MessageOrganizationUnitNotFound,
	MessageYouCannotDeleteYourOwnAccount,
	MessageFileCanNoLongerBeDeleted,
	MessageInvalidJoinDateOrJoinDateTimezone,
	MessageEmployeeNumberAlreadyExistsInThisTenant,
	MessageWorkShiftIsRequired,
	MessageFailedToGetEmployeeRole,
	MessageFailedToCreateUserAccount,
	MessageFailedToUpdateEmployeePassword,
	MessageFailedToCreateUser,
	MessageFailedToUpdateUser,
	MessageWorkLocationNameAlreadyExists,
	MessageAttendanceLogNotFound,
	MessageAttendanceRecordedSuccessfully,
	MessageCheckInAlreadyRecordedForToday,
	MessageCheckOutAlreadyRecordedForToday,
	MessageCheckInMustBeRecordedBeforeCheckOut,
	MessageCompanyPolicyNotFound,
	MessageEmployeeProfileNotFound,
	MessageExportJobMessageBusIsNotConfigured,
	MessageEmailMessageBusIsNotConfigured,
	MessageExportJobNotFound,
	MessageExportJobResourceMetadataIsRequired,
	MessageFileFilterIsRequired,
	MessageFileNotFound,
	MessageInvalidCategoryForTheSelectedParent,
	MessageInvalidLoggedAtFormat,
	MessageInvalidSelfieFile,
	MessageJobPositionNotFound,
	MessageOrganizationUnitCodeAlreadyExists,
	MessageRouteNotFound,
	MessageSelfieFileIsNoLongerAvailable,
	MessageSelfieFileNotFound,
	MessageUserRoleNotFound,
	MessageYouAreNotAllowedToAccessExportJobs,
	MessageYouAreNotAllowedToAccessThisResource,
	MessageYouAreNotAllowedToExportAttendanceReports,
	MessageYouDontHaveAccessToThisCompany,
	MessageProductNotFound,
	MessageFileIsRequired,
	MessageBodyIsRequired,
	MessageFilenameIsRequired,
	MessageFailedToGetPageOfResults,
	MessageAcceptedInvitationCannotBeResent,
	MessageAcceptedInvitationCannotBeRevoked,
	MessageActiveProductPricingNotFound,
	MessageAmountMustBePositive,
	MessageAmountsMustNotBeNegative,
	MessageAnActiveInvitationAlreadyExists,
	MessageClientOwnerRoleNotFound,
	MessageCompanyCanOnlyBeUpdatedFromCompanyDetails,
	MessageCompanyCannotBeDeleted,
	MessageCurrencyCodeFormatIsInvalid,
	MessageCurrencyCodeAlreadyExists,
	MessageCurrencyIsNotActive,
	MessageCurrencyNotFound,
	MessageProviderCurrencyIsNotActive,
	MessagePriceIsNotActive,
	MessageDefaultCurrencyCannotBeInactive,
	MessageDefaultCurrencyCannotBeDeleted,
	MessageCurrencyIsStillUsed,
	MessageDokuClientIsNotConfigured,
	MessageEmployeeIdIsNotAllowedForAdminInvitations,
	MessageEmployeeIdIsRequiredForEmployeeInvitations,
	MessageEmployeeAlreadyHasAnActiveUser,
	MessageFailedToAcceptInvitation,
	MessageFailedToCreateInvitation,
	MessageFailedToHashPassword,
	MessageFailedToResendInvitation,
	MessageFullNameIsRequired,
	MessageGoogleRegistrationSessionCreated,
	MessageGoogleRegistrationSessionIsInvalidOrExpired,
	MessageInternalProductCodeAlreadyExists,
	MessageInternalProductCodeFormatIsInvalid,
	MessageInternalProductNotFound,
	MessageInternalProductPriceNotFound,
	MessageInternalProductPricingCodeAlreadyExists,
	MessageInternalProductPricingCodeFormatIsInvalid,
	MessageInternalProductPricingNotFound,
	MessageInternalProductPricingStatusIsInvalid,
	MessageInternalProductStatusIsInvalid,
	MessageInternalUserEmailAlreadyExists,
	MessageInternalUserIsNotActive,
	MessageInternalUserNotFound,
	MessageInternalUserRoleIsInvalid,
	MessageInternalUserStatusIsInvalid,
	MessageInvalidInvitationRole,
	MessageInvalidRefreshToken,
	MessageInvitationEmailBusIsNotConfigured,
	MessageInvitationEmailMustMatchEmployeeEmail,
	MessageInvitationHasAlreadyBeenAccepted,
	MessageInvitationHasBeenRevoked,
	MessageInvitationHasExpired,
	MessageInvitationIsAlreadyRevoked,
	MessageInvitationNotFound,
	MessageInvoiceIsNotPayable,
	MessageInvoiceNotFound,
	MessageOneOrMorePermissionsAreInvalidForThisClient,
	MessageOnlyApprovedQuotationCanBeConverted,
	MessageOnlyOwnerCanInviteAdminUsers,
	MessageOnlyOwnerCanManageAdminInvitations,
	MessageOrderNotFound,
	MessagePasswordIsRequired,
	MessagePaymentAmountDoesNotMatchInvoiceOutstandingAmount,
	MessagePaymentAttemptCannotBeRetried,
	MessagePaymentAttemptNotFound,
	MessagePricePeriodOverlapsWithAnExistingActivePrice,
	MessageQuotationCannotBeApprovedFromCurrentStatus,
	MessageQuotationNotFound,
	MessageRevokedInvitationCannotBeResent,
	MessageStartedAtMustUseRfc3339Format,
	MessageEndedAtMustUseRfc3339Format,
	MessageEndedAtMustBeLaterThanStartedAt,
	MessageTenantIsRequired,
	MessageUnsupportedInvoiceAction,
	MessageUnsupportedPaymentAttemptAction,
	MessageUnsupportedQuotationAction,
	MessageUserActionTokenNotFound,
	MessageWebhookProcessed,
	MessageYouAreNotAuthorizedToInviteUsers,
	MessageYouAreNotAuthorizedToManageInternalClientPermissions,
	MessageYouAreNotAuthorizedToManageInternalProducts,
	MessageYouAreNotAuthorizedToManageInternalUsers,
	MessageYouAreNotAuthorizedToManageInvitations,
	MessageYouAreNotAuthorizedToViewInternalClients,
	MessageYouAreNotAuthorizedToViewInvitations,
	MessageValidationDefault,
	MessageValidationDefaultWithParam,
	MessageValidationRequired,
	MessageValidationEmail,
	MessageValidationEmailBlacklist,
	MessageValidationStrongPassword,
	MessageValidationResourceNotExist,
	MessageValidationDatetime,
	MessageValidationSimpleInvalid,
	MessageValidationSimpleFormat,
	MessageValidationMinValue,
	MessageValidationMinChars,
	MessageValidationMinItems,
	MessageValidationMaxValue,
	MessageValidationMaxChars,
	MessageValidationMaxItems,
	MessageValidationCompareGt,
	MessageValidationCompareGte,
	MessageValidationCompareLt,
	MessageValidationCompareLte,
	MessageValidationNumeric,
	MessageValidationTimezone,
	MessageValidationEqfield,
	MessageValidationOneof,
	MessageValidationUniqueInSlice,
}
