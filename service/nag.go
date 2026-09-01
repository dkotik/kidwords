package service

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

/*
<button hx-get="/target-endpoint"
        hx-target="#result"
        hx-on:htmx:config-request="if (!document.cookie.includes('your_cookie_name=expected_value')) event.preventDefault()">
    Click Me
</button>

function setCookie(name, value, daysToLive) {
    // Encode value to handle special characters or spaces safely
    let cookieString = encodeURIComponent(name) + "=" + encodeURIComponent(value);

    if (daysToLive) {
        // Calculate max-age in seconds
        cookieString += "; max-age=" + (daysToLive * 24 * 60 * 60);
    }

    // Add essential fallback attributes
    cookieString += "; path=/; secure; samesite=lax";

    // Write to the browser
    document.cookie = cookieString;
}


<button hx-on:click="setCookie("theme", "dark mode", 30);">Click Me</button>
*/

type nagView struct {
	lc      *i18n.Localizer
	User    User
	HasKeys bool
}

func (s *Service) nagView(ctx context.Context) (view nagView, err error) {
	view.User, view.lc, err = s.unpackContext(ctx)
	if err != nil {
		return view, err
	}

	keys, err := s.repository.List(ctx, view.User.GetID())
	if err != nil {
		return view, err
	}
	if len(keys) > 0 {
		view.HasKeys = true
	}

	return view, nil
}
