package email

import "fmt"

// VerificationEmailTemplate generates email verification email
func VerificationEmailTemplate(name, verificationURL, expiresIn string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email - MenuVista</title>
</head>
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
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Thank you for joining MenuVista! We're excited to help you create beautiful digital menus for your restaurants.</p>
                            <p style="color: #4b5563; margin: 0 0 24px 0; font-size: 16px; line-height: 1.6;">To get started, please verify your email address by clicking the button below:</p>
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 0 0 24px 0;">
                                        <a href="%s" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Verify Email Address</a>
                                    </td>
                                </tr>
                            </table>
                            <p style="color: #6b7280; margin: 0 0 8px 0; font-size: 14px;">Or copy and paste this link:</p>
                            <p style="color: #667eea; margin: 0 0 24px 0; font-size: 14px; word-break: break-all;">%s</p>
                            <div style="background: #fef3c7; border-left: 4px solid #f59e0b; padding: 16px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #92400e; margin: 0; font-size: 14px;">⏰ <strong>This link expires in %s</strong></p>
                            </div>
                            <hr style="border: none; border-top: 1px solid #e5e7eb; margin: 32px 0;">
                            <p style="color: #4b5563; margin: 0 0 12px 0; font-size: 15px; font-weight: 600;">What happens next?</p>
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr><td style="padding: 8px 0;"><span style="color: #10b981; font-size: 16px;">✅</span> <span style="color: #6b7280; font-size: 14px; margin-left: 8px;">Verify your email address</span></td></tr>
                                <tr><td style="padding: 8px 0;"><span style="color: #10b981; font-size: 16px;">🎁</span> <span style="color: #6b7280; font-size: 14px; margin-left: 8px;">Get 14 days of <strong>free trial</strong></span></td></tr>
                                <tr><td style="padding: 8px 0;"><span style="color: #10b981; font-size: 16px;">🍽️</span> <span style="color: #6b7280; font-size: 14px; margin-left: 8px;">Create your first restaurant</span></td></tr>
                                <tr><td style="padding: 8px 0;"><span style="color: #10b981; font-size: 16px;">📱</span> <span style="color: #6b7280; font-size: 14px; margin-left: 8px;">Build beautiful digital menus</span></td></tr>
                            </table>
                            <p style="color: #9ca3af; margin: 32px 0 0 0; font-size: 12px;">If you didn't create this account, you can safely ignore this email.</p>
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
`, name, verificationURL, verificationURL, expiresIn)
}

// WelcomeEmailTemplate generates welcome email after verification
func WelcomeEmailTemplate(name, trialEndDate, dashboardURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
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
                            <h2 style="color: #1f2937; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Welcome aboard, %s! 🚀</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Your email has been verified and your account is now active!</p>
                            <p style="color: #4b5563; margin: 0 0 24px 0; font-size: 16px; line-height: 1.6;">You now have access to all MenuVista features for the next <strong>14 days</strong> completely free.</p>
                            <div style="background: #f0fdf4; border-left: 4px solid #10b981; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #047857; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">🎁 Your Free Trial is Active</p>
                                <p style="color: #065f46; margin: 0; font-size: 14px;">Trial ends on: <strong>%s</strong></p>
                            </div>
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 32px 0 24px 0;">
                                        <a href="%s" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Go to Dashboard</a>
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
`, name, trialEndDate, dashboardURL)
}

// TrialWarningEmailTemplate generates trial warning emails
func TrialWarningEmailTemplate(name, trialEndDate string, daysLeft int, upgradeBronzeURL, upgradeSilverURL, upgradeGoldURL string) string {
	var title, urgencyColor, urgencyBg string

	switch daysLeft {
	case 4:
		title = "Your Free Trial Ends in 4 Days ⏰"
		urgencyColor = "#f59e0b"
		urgencyBg = "#fef3c7"
	case 3:
		title = "Your Free Trial Ends in 3 Days ⏰"
		urgencyColor = "#f59e0b"
		urgencyBg = "#fef3c7"
	case 2:
		title = "Your Free Trial Ends in 2 Days ⏰"
		urgencyColor = "#f59e0b"
		urgencyBg = "#fef3c7"
	case 1:
		title = "⚠️ Your Free Trial Expires Tomorrow"
		urgencyColor = "#ef4444"
		urgencyBg = "#fef2f2"
	default:
		title = "Your Free Trial Ends Soon ⏰"
		urgencyColor = "#f59e0b"
		urgencyBg = "#fef3c7"
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
                            <h2 style="color: #1f2937; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">%s</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 24px 0; font-size: 16px; line-height: 1.6;">Your 14-day free trial ends in <strong>%d days</strong> on <strong>%s</strong>.</p>
                            <div style="background: %s; border-left: 4px solid %s; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #92400e; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">⚠️ Action Required</p>
                                <p style="color: #78350f; margin: 0; font-size: 14px;">To continue using MenuVista, please upgrade to a paid plan.</p>
                            </div>
                            <p style="color: #4b5563; margin: 24px 0 16px 0; font-size: 15px; font-weight: 600;">Choose Your Plan:</p>
                            <div style="border: 2px solid #e5e7eb; border-radius: 8px; margin: 16px 0; overflow: hidden;">
                                <div style="background: linear-gradient(135deg, #cd7f32 0%%, #b87333 100%%); padding: 16px; text-align: center;">
                                    <h3 style="color: #ffffff; margin: 0; font-size: 18px;">🥉 Bronze Plan - $19.99/month</h3>
                                </div>
                                <div style="padding: 20px;">
                                    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                        <tr>
                                            <td align="center">
                                                <a href="%s" style="display: inline-block; background: #cd7f32; color: #ffffff; padding: 12px 32px; border-radius: 6px; text-decoration: none; font-weight: 600; font-size: 14px;">Choose Bronze</a>
                                            </td>
                                        </tr>
                                    </table>
                                </div>
                            </div>
                            <div style="border: 3px solid #667eea; border-radius: 8px; margin: 16px 0; overflow: hidden;">
                                <div style="background: linear-gradient(135deg, #c0c0c0 0%%, #a8a8a8 100%%); padding: 16px; text-align: center;">
                                    <h3 style="color: #1f2937; margin: 0; font-size: 18px;">🥈 Silver Plan - $49.99/month (POPULAR)</h3>
                                </div>
                                <div style="padding: 20px; background: #f9fafb;">
                                    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                        <tr>
                                            <td align="center">
                                                <a href="%s" style="display: inline-block; background: #667eea; color: #ffffff; padding: 12px 32px; border-radius: 6px; text-decoration: none; font-weight: 600; font-size: 14px;">Choose Silver</a>
                                            </td>
                                        </tr>
                                    </table>
                                </div>
                            </div>
                            <div style="border: 2px solid #e5e7eb; border-radius: 8px; margin: 16px 0; overflow: hidden;">
                                <div style="background: linear-gradient(135deg, #ffd700 0%%, #ffed4e 100%%); padding: 16px; text-align: center;">
                                    <h3 style="color: #1f2937; margin: 0; font-size: 18px;">🥇 Gold Plan - $99.99/month</h3>
                                </div>
                                <div style="padding: 20px;">
                                    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                        <tr>
                                            <td align="center">
                                                <a href="%s" style="display: inline-block; background: #ffd700; color: #1f2937; padding: 12px 32px; border-radius: 6px; text-decoration: none; font-weight: 600; font-size: 14px;">Choose Gold</a>
                                            </td>
                                        </tr>
                                    </table>
                                </div>
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
`, title, name, daysLeft, trialEndDate, urgencyBg, urgencyColor, upgradeBronzeURL, upgradeSilverURL, upgradeGoldURL)
}

// TrialExpiredEmailTemplate generates trial expired email
func TrialExpiredEmailTemplate(name, upgradeURL string) string {
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
                            <h2 style="color: #dc2626; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">Your Free Trial Has Expired 😔</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 24px 0; font-size: 16px; line-height: 1.6;">Your 14-day free trial of MenuVista has expired. To continue managing your digital menus, please upgrade to a paid plan.</p>
                            <div style="background: #fef2f2; border-left: 4px solid #dc2626; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #991b1b; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">⛔ Trial Ended</p>
                                <p style="color: #7f1d1d; margin: 0; font-size: 14px;">Your account is now on the free tier with limited features. Upgrade to restore full access.</p>
                            </div>
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="%s" style="display: inline-block; background: #dc2626; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Upgrade Now</a>
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
`, name, upgradeURL)
}

// SubscriptionExpiringEmailTemplate generates subscription expiring warning for paid plans
func SubscriptionExpiringEmailTemplate(name, planName, expiryDate string, daysLeft int, renewURL string) string {
	var urgencyLevel string
	if daysLeft <= 2 {
		urgencyLevel = "URGENT"
	} else {
		urgencyLevel = "REMINDER"
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
                            <h2 style="color: #f59e0b; margin: 0 0 16px 0; font-size: 24px; font-weight: 600;">%s: Subscription Renewal Due</h2>
                            <p style="color: #4b5563; margin: 0 0 16px 0; font-size: 16px; line-height: 1.6;">Hi <strong>%s</strong>,</p>
                            <p style="color: #4b5563; margin: 0 0 24px 0; font-size: 16px; line-height: 1.6;">Your <strong>%s</strong> subscription will expire in <strong>%d days</strong> on <strong>%s</strong>.</p>
                            <div style="background: #fef3c7; border-left: 4px solid #f59e0b; padding: 20px; border-radius: 4px; margin: 24px 0;">
                                <p style="color: #92400e; margin: 0 0 8px 0; font-size: 16px; font-weight: 600;">⚠️ Action Required</p>
                                <p style="color: #78350f; margin: 0; font-size: 14px;">Please ensure your payment method is up to date to avoid service interruption.</p>
                            </div>
                            <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
                                <tr>
                                    <td align="center" style="padding: 24px 0;">
                                        <a href="%s" style="display: inline-block; background: #667eea; color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Update Payment Method</a>
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
`, urgencyLevel, name, planName, daysLeft, expiryDate, renewURL)
}
