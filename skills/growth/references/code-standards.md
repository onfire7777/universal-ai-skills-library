# Growth Code Standards

## Good Growth Code

```typescript
// Rich Snippet (JSON-LD) for Search Engines
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Product",
  "name": "Awesome Tool",
  "description": "Boost your productivity..."
}
</script>

// Clear Call-to-Action (CTA) with descriptive link
<a href="/signup" className="btn-primary" onClick={trackSignupClick}>
  Start your free trial
</a>
```

## Bad Growth Code

```typescript
// "Click here" is bad for SEO and Accessibility
<a href="/signup">Click here</a>

// Missing Open Graph tags (looks ugly on Twitter/Slack)
<head>
  <title>Home</title>
  {/* No description, no image... */}
</head>
```
