import Link from 'next/link'
import BrandLogo from '@/components/common/BrandLogo'
import type { LandingPageCopy } from '@/lib/i18n/pages/landing'

export default function Footer({ copy, homeHref }: { copy: LandingPageCopy['footer']; homeHref: string }) {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="landingFooter" aria-label="LearnCosmos footer">
      <div className="landingFooterInner">
        <div className="landingFooterOrbit" aria-hidden="true" />

        <div className="landingFooterBrandPanel">
          <div className="landingFooterBrand">
            <BrandLogo href={homeHref} iconSize={38} textSize="22px" />
          </div>
          <p className="landingFooterTagline">
            {copy.taglineLines.map((line) => (
              <span key={line}>{line}</span>
            ))}
          </p>
        </div>

        <nav className="landingFooterColumns" aria-label="Footer links">
          {copy.columns.map((column) => (
            <section key={column.title} className="landingFooterColumn" aria-label={column.title}>
              <h2>{column.title}</h2>
              <ul>
                {column.links.map((item) => (
                  <li key={`${column.title}-${item.label}`}>
                    {!item.href ? (
                      <span>{item.label}</span>
                    ) : item.href.startsWith('http') || item.href.startsWith('mailto:') ? (
                      <a
                        href={item.href}
                        target={item.href.startsWith('http') ? '_blank' : undefined}
                        rel={item.href.startsWith('http') ? 'noreferrer' : undefined}
                      >
                        {item.label}
                      </a>
                    ) : (
                      <Link href={item.href}>{item.label}</Link>
                    )}
                  </li>
                ))}
              </ul>
            </section>
          ))}
        </nav>

        <div className="landingFooterBottom">
          <p>© {currentYear} LearnCosmos. All rights reserved.</p>
          <p>AI based learning tool platform LearnCosmos</p>
        </div>
      </div>
    </footer>
  )
}
