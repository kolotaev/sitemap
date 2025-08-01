package sitemap_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/refurbed/sitemap"
)

func Example() {
	sm := sitemap.New()
	t := time.Unix(0, 0).UTC()
	sm.Add(&sitemap.URL{
		Loc:        "http://example.com/",
		LastMod:    &t,
		ChangeFreq: sitemap.Daily,
	})
	sm.WriteTo(os.Stdout)
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xhtml="http://www.w3.org/1999/xhtml">
	//   <url>
	//     <loc>http://example.com/</loc>
	//     <lastmod>1970-01-01T00:00:00Z</lastmod>
	//     <changefreq>daily</changefreq>
	//   </url>
	// </urlset>
}

func TestSitemap_WriteTo(t *testing.T) {
	sm := sitemap.New()
	ts := time.Date(2025, 12, 30, 15, 55, 55, 0, time.UTC)
	sm.Add(&sitemap.URL{
		Loc:        "http://example.com/",
		LastMod:    &ts,
		ChangeFreq: sitemap.Daily,
		Priority:   2.5,
		Alternates: []*sitemap.AltURL{
			{
				Rel:      "alternate",
				HrefLang: "de-at",
				Loc:      "http://example.com/de-at/",
			},
		},
	})

	expectOut := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xhtml="http://www.w3.org/1999/xhtml">
  <url>
    <loc>http://example.com/</loc>
    <lastmod>2025-12-30T15:55:55Z</lastmod>
    <changefreq>daily</changefreq>
    <priority>2.5</priority>
    <xhtml:link rel="alternate" hreflang="de-at" href="http://example.com/de-at/"></xhtml:link>
  </url>
</urlset>
`

	out := bytes.NewBuffer([]byte{})
	_, err := sm.WriteTo(out)

	if err != nil {
		t.Fatalf("sitemap WriteTo got unexpected error: %q", err)
	}
	if expectOut != out.String() {
		t.Errorf("sitemap WriteTo error.\nWant: %s\nGot: %s", expectOut, out.String())
	}
}

func TestSitemap_ReadFrom(t *testing.T) {
	contents := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xhtml="http://www.w3.org/1999/xhtml">
  <url>
    <loc>http://example.com/</loc>
    <lastmod>2025-12-30T15:55:55Z</lastmod>
    <changefreq>daily</changefreq>
    <priority>2.5</priority>
    <xhtml:link rel="alternate" hreflang="de-at" href="http://example.com/de-at/"></xhtml:link>
    <xhtml:link rel="alternate" hreflang="fr-fr" href="http://example.com/fr-fr/"></xhtml:link>
  </url>
</urlset>
`
	actualSm := sitemap.New()
	_, err := actualSm.ReadFrom(bytes.NewBufferString(contents))
	if err != nil {
		t.Fatalf("sitemap WriteTo got unexpected error: %q", err)
	}

	actual := bytes.NewBuffer([]byte{})
	actualSm.WriteTo(actual)

	if contents != actual.String() {
		t.Errorf("sitemap ReadFrom error.\nWant: %s\nGot: %s", contents, actual.String())
	}
}

func TestSitemapEncodeDecodeRoundtrip(t *testing.T) {
	expectedSitemap := sitemap.New()
	ts := time.Date(2025, 12, 30, 15, 55, 55, 0, time.UTC)
	expectedSitemap.Add(&sitemap.URL{
		Loc:        "http://example.com/",
		LastMod:    &ts,
		ChangeFreq: sitemap.Daily,
		Priority:   2.5,
		Alternates: []*sitemap.AltURL{
			{
				Rel:      "alternate",
				HrefLang: "de-at",
				Loc:      "http://example.com/de-at/",
			},
			{
				Rel:      "alternate",
				HrefLang: "fr-fr",
				Loc:      "http://example.com/fr-fr/",
			},
		},
	})
	expected := bytes.NewBuffer([]byte{})
	_, err := expectedSitemap.WriteTo(expected)
	if err != nil {
		t.Fatalf("expectedSitemap WriteTo got unexpected error: %q", err)
	}

	actualSitemap := sitemap.New()
	_, err = actualSitemap.ReadFrom(bytes.NewBufferString(expected.String()))
	if err != nil {
		t.Fatalf("actualSitemap WriteTo got unexpected error: %q", err)
	}

	actual := bytes.NewBuffer([]byte{})
	actualSitemap.WriteTo(actual)
	if expected.String() != actual.String() {
		t.Errorf("sitemap encode-decode roundtrip error.\nWant: %s\nGot: %s", expected.String(), actual.String())
	}
}
