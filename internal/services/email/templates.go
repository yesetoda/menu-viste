package email

import "fmt"

const (
	// Email templates
	adminNotificationNewOwnerSubject = "New Restaurant Owner Registration - Action Required"
	ownerApprovalSubject             = "Your MenuVista Account Has Been Approved!"
	ownerRejectionSubject            = "MenuVista Account Registration Update"
	restaurantApprovalSubject        = "Your Restaurant Has Been Approved!"
	restaurantRejectionSubject       = "Restaurant Submission Update"

	// Payment templates
	paymentSuccessSubject = "Payment Confirmation - MenuVista"
	paymentFailedSubject  = "Payment Failed - Action Required"
	paymentPendingSubject = "Payment Pending - MenuVista"
)

// AdminNotificationNewOwnerTemplate generates email for admins when a new owner registers
func AdminNotificationNewOwnerTemplate(ownerEmail, firstName, lastName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #1f2937; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">New Owner Registration 👤</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">A new restaurant owner has registered and is awaiting approval.</p>
                            
                            <div style="background: #f3f4f6; padding: 20px; border-radius: 8px; margin: 24px 0;">
                                <h3 style="color: #374151; margin: 0 0 12px 0; font-size: 18px;">Owner Details:</h3>
                                <p style="color: #4b5563; margin: 0 0 8px 0;"><strong>Name:</strong> %s %s</p>
                                <p style="color: #4b5563; margin: 0;"><strong>Email:</strong> %s</p>
                            </div>
                            
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="http://localhost:8080/admin/pending-users" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Review in Dashboard</a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName, lastName, ownerEmail)
}

// OwnerApprovalTemplate generates approval email for restaurant owners
func OwnerApprovalTemplate(firstName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #1f2937; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Welcome to MenuVista! 🎉</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Great news! Your MenuVista account has been approved.</p>
                            
                            <div style="background: #ecfdf5; border-left: 4px solid #10b981; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #065f46; margin: 0; font-size: 16px;">✅ <strong>Account Status: Active</strong></p>
                            </div>

                            <p style="color: #4b5563; margin: 0 0 12px 0; font-size: 16px; font-weight: 600;">You can now:</p>
                            <ul style="color: #4b5563; margin: 0 0 24px 0; padding-left: 20px;">
                                <li style="margin-bottom: 8px;">Create and manage your restaurant profiles</li>
                                <li style="margin-bottom: 8px;">Design beautiful digital menus</li>
                                <li style="margin-bottom: 8px;">Generate QR codes for easy customer access</li>
                            </ul>
                            
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="http://localhost:8080/login" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Get Started</a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName)
}

// OwnerRejectionTemplate generates rejection email for restaurant owners
func OwnerRejectionTemplate(firstName, reason string) string {
	reasonText := "your application did not meet our current requirements"
	if reason != "" {
		reasonText = reason
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #dc2626; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Registration Update</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Thank you for your interest in MenuVista.</p>
                            
                            <div style="background: #fef2f2; border-left: 4px solid #dc2626; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #991b1b; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">Application Status: Not Approved</p>
                                <p style="color: #7f1d1d; margin: 0; font-size: 14px;">%s.</p>
                            </div>
                            
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">If you believe this is an error or would like to discuss this decision, please contact our support team.</p>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName, reasonText)
}

// RestaurantApprovalTemplate generates approval email for restaurant submissions
func RestaurantApprovalTemplate(firstName, restaurantName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #1f2937; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Restaurant Approved! 🎊</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Excellent news! Your restaurant "<strong>%s</strong>" has been approved and is now live on MenuVista.</p>
                            
                            <div style="background: #ecfdf5; border-left: 4px solid #10b981; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #065f46; margin: 0; font-size: 16px;">✅ <strong>Status: Live & Visible</strong></p>
                            </div>

                            <p style="color: #4b5563; margin: 0 0 12px 0; font-size: 16px; font-weight: 600;">Your restaurant is now:</p>
                            <ul style="color: #4b5563; margin: 0 0 24px 0; padding-left: 20px;">
                                <li style="margin-bottom: 8px;">Visible to customers</li>
                                <li style="margin-bottom: 8px;">Ready to accept menu updates</li>
                                <li style="margin-bottom: 8px;">Available for QR code generation</li>
                            </ul>
                            
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="http://localhost:8080/my-restaurants" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Manage Restaurant</a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName, restaurantName)
}

// RestaurantRejectionTemplate generates rejection email for restaurant submissions
func RestaurantRejectionTemplate(firstName, restaurantName, reason string) string {
	reasonText := "it did not meet our quality standards"
	if reason != "" {
		reasonText = reason
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #dc2626; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Restaurant Submission Update</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Thank you for submitting "<strong>%s</strong>" to MenuVista.</p>
                            
                            <div style="background: #fef2f2; border-left: 4px solid #dc2626; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #991b1b; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">Status: Not Approved</p>
                                <p style="color: #7f1d1d; margin: 0; font-size: 14px;">Reason: %s.</p>
                            </div>
                            
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">You can review and update your restaurant information and resubmit for approval.</p>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName, restaurantName, reasonText)
}

// PaymentSuccessTemplate generates email for successful payment
func PaymentSuccessTemplate(firstName, invoiceNumber string, amount float64, currency string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #1f2937; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Payment Successful! 🎉</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Thank you for your payment. Your subscription has been activated.</p>
                            
                            <div style="background: #f3f4f6; padding: 20px; border-radius: 8px; margin: 24px 0;">
                                <h3 style="color: #374151; margin: 0 0 12px 0; font-size: 18px;">Payment Details:</h3>
                                <p style="color: #4b5563; margin: 0 0 8px 0;"><strong>Invoice:</strong> %s</p>
                                <p style="color: #4b5563; margin: 0 0 8px 0;"><strong>Amount:</strong> %.2f %s</p>
                                <p style="color: #059669; margin: 0; font-weight: 600;"><strong>Status:</strong> Paid</p>
                            </div>
                            
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="http://localhost:8080/billing" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">View Invoice</a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName, invoiceNumber, amount, currency)
}

// PaymentFailedTemplate generates email for failed payment
func PaymentFailedTemplate(firstName, updatePaymentURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #dc2626; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Payment Failed ⚠️</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">We were unable to process your payment for your MenuVista subscription.</p>
                            
                            <div style="background: #fef2f2; border-left: 4px solid #dc2626; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #991b1b; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">Action Required</p>
                                <p style="color: #7f1d1d; margin: 0; font-size: 14px;">Please update your payment method to avoid service interruption.</p>
                            </div>
                            
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="%s" style="display: inline-block; background: #dc2626; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Update Payment Method</a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="color: #6b7280; margin: 24px 0 0 0; font-size: 14px; text-align: center;">We will retry the payment automatically in 24 hours.</p>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName, updatePaymentURL)
}

// PaymentPendingTemplate generates email for pending payment
func PaymentPendingTemplate(firstName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #f59e0b; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Payment Pending ⏳</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">We have received your payment request and it is currently being processed.</p>
                            
                            <div style="background: #fef3c7; border-left: 4px solid #f59e0b; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #92400e; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">Processing</p>
                                <p style="color: #78350f; margin: 0; font-size: 14px;">We will notify you once the payment is confirmed. No action is required from you at this time.</p>
                            </div>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName)
}

// StaffWelcomeTemplate generates welcome email for new staff
func StaffWelcomeTemplate(firstName, email, password string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f3f4f6;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background-color: #f3f4f6;">
        <tr>
            <td style="padding: 40px 20px;">
                <table role="presentation" width="600" cellspacing="0" cellpadding="0" border="0" style="margin: 0 auto; max-width: 600px;">
                    <tr>
                        <td style="text-align: center; padding-bottom: 32px;">
                            <h1 style="color: #667eea; font-size: 32px; margin: 0; font-weight: 700;">🍽️ MenuVista</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="background: #ffffff; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                            <h2 style="color: #1f2937; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Welcome to the Team! 👋</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">You have been added as a staff member to MenuVista.</p>
                            
                            <div style="background: #f3f4f6; padding: 20px; border-radius: 8px; margin: 24px 0;">
                                <h3 style="color: #374151; margin: 0 0 12px 0; font-size: 18px;">Your Login Credentials:</h3>
                                <p style="color: #4b5563; margin: 0 0 8px 0;"><strong>Email:</strong> %s</p>
                                <p style="color: #4b5563; margin: 0;"><strong>Password:</strong> <code style="background: #e5e7eb; padding: 2px 6px; border-radius: 4px; font-family: monospace;">%s</code></p>
                            </div>
                            
                            <p style="color: #ef4444; margin: 0 0 24px 0; font-size: 14px;">⚠️ Please log in and change your password immediately.</p>
                            
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="http://localhost:8080/login" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Log In Now</a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding-top: 32px; text-align: center;">
                            <p style="color: #9ca3af; font-size: 12px; margin: 0;">© 2026 MenuVista. All rights reserved.</p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, firstName, email, password)
}
