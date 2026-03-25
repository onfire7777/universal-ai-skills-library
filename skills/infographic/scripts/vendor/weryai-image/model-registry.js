import { createClient, log } from '../weryai-core/client.js';
import { MODELS_API_PATH } from './models.js';

export async function fetchModelRegistry(ctx) {
  try {
    const client = createClient({ ...ctx, baseUrl: ctx.modelsBaseUrl });
    const res = await client.get(MODELS_API_PATH, { retries: 1 });

    if (res.httpStatus < 200 || res.httpStatus >= 300 || !res.data) {
      log('Warning: models API returned non-success. Using permissive mode.');
      return null;
    }

    const data = res.data;
    return {
      text_to_image: indexByKey(data.text_to_image),
      image_to_image: indexByKey(data.image_to_image),
    };
  } catch (err) {
    log(`Warning: could not fetch model registry (${err.message}). Using permissive mode.`);
    return null;
  }
}

function indexByKey(arr) {
  const map = new Map();
  if (!Array.isArray(arr)) return map;
  for (const item of arr) {
    if (item.model_key) map.set(item.model_key, item);
  }
  return map;
}

export function lookupModel(registry, modelKey, mode) {
  if (!registry) return null;
  const modeMap = registry[mode];
  if (!modeMap) return null;
  return modeMap.get(modelKey) ?? null;
}

export function validateWithModel(meta, input, mode) {
  const errors = [];
  const warnings = [];

  if (input.prompt && meta.prompt_length_limit && input.prompt.length > meta.prompt_length_limit) {
    errors.push(`prompt exceeds model limit (${input.prompt.length}/${meta.prompt_length_limit} chars).`);
  }

  const aspectRatio = input.aspect_ratio || input.aspectRatio;
  const allowedSizes = meta.image_sizes;
  if (aspectRatio && Array.isArray(allowedSizes) && allowedSizes.length > 0 && !allowedSizes.includes(aspectRatio)) {
    errors.push(
      `Model "${meta.model_key}" does not support aspect_ratio "${aspectRatio}". Allowed: [${allowedSizes.join(', ')}].`,
    );
  }

  const imageNumber = input.image_number ?? input.imageNumber;
  const allowedNums = meta.num_images;
  if (imageNumber != null && Array.isArray(allowedNums) && allowedNums.length > 0) {
    const value = Number(imageNumber);
    if (Number.isNaN(value)) {
      errors.push(`image_number must be a number, got "${imageNumber}".`);
    } else if (!allowedNums.includes(value)) {
      errors.push(
        `Model "${meta.model_key}" does not support image_number=${value}. Allowed: [${allowedNums.join(', ')}].`,
      );
    }
  }

  const resolution = input.resolution;
  const allowedRes = meta.resolutions;
  if (resolution && Array.isArray(allowedRes) && allowedRes.length > 0 && !allowedRes.includes(resolution)) {
    errors.push(
      `Model "${meta.model_key}" does not support resolution "${resolution}". Allowed: [${allowedRes.join(', ')}].`,
    );
  }

  if (mode === 'image_to_image') {
    const images = input.images;
    const limit = meta.upload_image_limit;
    if (Array.isArray(images) && limit != null && images.length > limit) {
      warnings.push(
        `Model "${meta.model_key}" allows max ${limit} reference image(s). Extra images will be trimmed.`,
      );
    }
  }

  return { errors, warnings };
}

export function buildBody(meta, input, mode) {
  const body = {
    prompt: input.prompt,
    model: meta.model_key,
  };

  const resolvedAspect = resolveFromAllowed(input.aspect_ratio || input.aspectRatio, meta.image_sizes, '1:1');
  if (resolvedAspect != null) body.aspect_ratio = resolvedAspect;

  const resolvedNum = resolveNumeric(input.image_number ?? input.imageNumber, meta.num_images, meta.num_images?.[0] ?? 1);
  if (resolvedNum != null) body.image_number = resolvedNum;

  const resolvedRes = resolveFromAllowed(input.resolution, meta.resolutions, '1k');
  if (resolvedRes != null) body.resolution = resolvedRes;

  if (input.negative_prompt || input.negativePrompt) {
    body.negative_prompt = input.negative_prompt || input.negativePrompt;
  }

  if (mode === 'image_to_image' && Array.isArray(input.images)) {
    const limit = meta.upload_image_limit ?? 1;
    body.images = input.images.slice(0, limit);
  }

  return body;
}

function resolveFromAllowed(userValue, allowed, preferred) {
  if (!Array.isArray(allowed) || allowed.length === 0) return null;
  if (userValue && allowed.includes(userValue)) return userValue;
  if (userValue && !allowed.includes(userValue)) {
    log(`Warning: "${userValue}" not in allowed list [${allowed.join(', ')}]. Using default.`);
    return allowed.includes(preferred) ? preferred : allowed[0];
  }
  return allowed.includes(preferred) ? preferred : allowed[0];
}

function resolveNumeric(userValue, allowed, fallback) {
  if (!Array.isArray(allowed) || allowed.length === 0) return fallback;
  if (userValue != null) {
    const value = Number(userValue);
    if (allowed.includes(value)) return value;
  }
  return allowed[0];
}
