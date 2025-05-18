package middleware

import "net/http"

func Headers(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
		 * Prevent clickjacking.
		 *
		 * Prevents your web pages from being embedded in frames on other
		 * websites.
		 *
		 * Protects against clickjacking attacks where attackers overlay your
		 * site with invisible elements to trick users into clicking malicious
		 * elements.
		 *
		 * Example attack: An attacker could embed your internal app in an
		 * invisible iframe and trick users into clicking buttons they didn't
		 * intend to
		 */
		w.Header().Set("X-Frame-Options", "DENY")

		/*
		 * Prevent MIME type sniffing.
		 *
		 * Prevents browsers from trying to guess ("sniff") a file's MIME type.
		 *
		 * Forces browsers to strictly use the Content-Type header you specify.
		 *
		 * Example attack: Without this, if you have a file upload feature, an
		 * attacker could upload a malicious JavaScript file but name it as an
		 * image. Some browsers might "sniff" it and execute it as JavaScript
		 * instead of treating it as an image
		 */
		w.Header().Set("X-Content-Type-Options", "nosniff")

		/*
		 * XSS protection.
		 *
		 * Enables the browser's built-in XSS filter.
		 *
		 * Blocks the page if a cross-site scripting attack is detected.
		 *
		 * While somewhat dated (modern browsers use Content-Security-Policy
		 * instead), it's still useful for older browsers.
		 *
		 * Example attack: Without this, if your app accidentally reflects user
		 * input without proper escaping, an attacker could inject malicious
		 * scripts.
		 */
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		/* Content security policy (internal use).
		 *
		 * This is like a firewall for your web page.
		 *
		 * default-src 'self': Only allow resources from your own domain.
		 *
		 * img-src 'self' data:: Images can only come from your domain or be
		 * inline base64.
		 *
		 * style-src 'self' 'unsafe-inline': Styles can only come from your
		 * domain or be inline.
		 *
		 * script-src 'self': Scripts can only come from your domain.
		 *
		 * Example attack: Without CSP, if an attacker somehow injected a script
		 * tag pointing to their malicious server, the browser would load and
		 * execute it. CSP prevents this by only allowing scripts from your
		 * domain.
		 *
		 */
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; "+
				"img-src 'self' data:; "+
				"style-src 'self' 'unsafe-inline'; "+
				"script-src 'self'",
		)

		/* Referrer policy.
		 *
		 * Controls how much referrer information is included when users
		 * navigate away from your pages.
		 *
		 * Only sends referrer information to pages on your own domain.
		 *
		 * Example privacy issue: Without this, if a user clicks a link to leave
		 * your app, the destination site would see the full URL they came from,
		 * which might include sensitive information in the path or query
		 * parameters.
		 */
		w.Header().Set("Referrer-Policy", "same-origin")

		next.ServeHTTP(w, r)
	}
}
