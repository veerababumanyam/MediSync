---
name: seo-optimization
description: Guidelines for high-performance SEO, GEO, and AI search presence for the MediSync platform.
---

# SEO & AI Search Optimization Skill

This skill provides comprehensive guidelines and technical requirements to ensure the MediSync platform is optimized for traditional search engines (Google, Bing), AI search engines (Perplexity, ChatGPT Search, Gemini), and modern generative engines (GEO).

## Core Pillars

### 1. Traditional SEO (Search Engine Optimization)
- **Technical Foundational**: Ensure fast LCP, low CLS, and excellent Core Web Vitals.
- **Semantic HTML**: Use proper heading hierarchy (`h1`-`h6`), `article`, `section`, and `nav` tags.
- **Meta Management**: Dynamic `<title>` and `<meta name="description">` for every route.
- **Robotic Signals**: Maintain `robots.txt` and a clean `sitemap.xml`.

### 2. GEO (Generative Engine Optimization)
- **Answer Eligibility**: Structure content as direct answers to common user questions (FAQ style).
- **Claim-Based Architecture**: Make verifiable claims followed by evidence or data citations.
- **Source Trustworthiness**: Link to authoritative healthcare sources and maintain clear E-E-A-T signals.

### 3. AI Search & Citation Optimization
- **Entity Focus**: Use clear terminology for "MediSync", "HIMS AI", "LIMS AI Integration", "Tally ERP", and "Agentic AI for Legacy Healthcare".
- **LLM Scannability**: Use short paragraphs, bullet points, and explicit summaries at the top of long pages.
- **Citation Hooks**: Include unique stats, original healthcare insights, and expert quotes that AI models are likely to cite.

### 4. Robotic & Crawler Optimization
- **Machine Readability**: Ensure core content is visible without JavaScript (SSR/Pre-rendering).
- **Internal Linking**: Maintain a shallow crawl depth (all pages within 3 clicks).
- **Clean Code**: Avoid excessive DOM nesting and ensure all interactive elements have `aria-label` or descriptive text.

## Technical Implementation

### Structured Data (JSON-LD)
Always include Schema.org markup. Use the `Organization`, `SoftwareApplication`, and `WebPage` schemas.
Example:
```json
{
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  "name": "MediSync",
  "operatingSystem": "Web",
  "applicationCategory": "BusinessApplication",
  "offers": {
    "@type": "Offer",
    "price": "0",
    "priceCurrency": "USD"
  }
}
```

### Route Optimization
- Use descriptive slugs: `/dashboard/billing-analytics` instead of `/d/b1`.
- Canonical tags: Ensure every page has a `<link rel="canonical" href="...">`.

## Usage
Apply these standards when:
1. Creating new landing pages or marketing content.
2. Modifying the `App.tsx` or routing logic.
3. Implementing data visualization summaries (to make them "citeable").
4. Running technical audits.
