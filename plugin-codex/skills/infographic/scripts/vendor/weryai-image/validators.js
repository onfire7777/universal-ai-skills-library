export function validateSubmitText(input) {
  const errors = [];

  if (!input.prompt || typeof input.prompt !== 'string' || input.prompt.trim().length === 0) {
    errors.push('prompt is required and must be a non-empty string.');
  }

  return errors;
}

export function validateSubmitImage(input) {
  const errors = validateSubmitText(input);

  if (!input.images || !Array.isArray(input.images) || input.images.length === 0) {
    errors.push('images is required for image-to-image and must be a non-empty array of URLs.');
  } else {
    for (const url of input.images) {
      if (typeof url !== 'string' || !url.startsWith('https://')) {
        errors.push(`Invalid image URL: "${url}". Only https:// URLs are supported.`);
      }
    }
  }

  return errors;
}
