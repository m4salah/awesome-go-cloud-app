package views

import (
	"canvas/k"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func FrontPage() g.Node {
	return Page(
		"Canvas",
		k.RootPath,
		H1(g.Text(`Solutions to problems.`)),
		P(g.Raw(`Do you have problems? We also had problems.`)),
		P(g.Raw(`Then we created the <em>canvas</em> app, and now we don't! 😬`)),
		P(g.Raw(`Hello 👋`)),
	)
}
