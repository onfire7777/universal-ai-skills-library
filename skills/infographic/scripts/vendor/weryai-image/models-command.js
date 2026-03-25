import { fetchModelRegistry } from './model-registry.js';

export async function execute(input, ctx) {
  const registry = await fetchModelRegistry(ctx);

  if (!registry) {
    return {
      ok: false,
      phase: 'failed',
      errorCode: 'REGISTRY_UNAVAILABLE',
      errorMessage: 'Could not fetch model list from WeryAI API. Check API key and network.',
    };
  }

  const mode = input.mode;

  if (mode && mode !== 'text_to_image' && mode !== 'image_to_image') {
    return {
      ok: false,
      phase: 'failed',
      errorCode: 'VALIDATION',
      errorMessage: 'Invalid mode. Use "text_to_image" or "image_to_image".',
    };
  }

  const result = { ok: true, phase: 'completed' };
  if (!mode || mode === 'text_to_image') result.text_to_image = Array.from(registry.text_to_image.values());
  if (!mode || mode === 'image_to_image') result.image_to_image = Array.from(registry.image_to_image.values());
  return result;
}

export default execute;
