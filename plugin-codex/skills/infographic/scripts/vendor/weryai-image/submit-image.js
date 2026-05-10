import { createClient, log } from '../weryai-core/client.js';
import { formatApiError, formatNetworkError, isApiSuccess } from '../weryai-core/errors.js';
import { buildBody, fetchModelRegistry, lookupModel, validateWithModel } from './model-registry.js';
import { DEFAULT_MODEL, FALLBACK_DEFAULTS } from './models.js';
import { validateSubmitImage } from './validators.js';

export async function execute(input, ctx) {
  const structErrors = validateSubmitImage(input);
  if (structErrors.length > 0) {
    return { ok: false, phase: 'failed', errorCode: 'VALIDATION', errorMessage: structErrors.join(' ') };
  }

  const registry = await fetchModelRegistry(ctx);
  const model = input.model || DEFAULT_MODEL;
  const meta = lookupModel(registry, model, 'image_to_image');

  if (meta) {
    const { errors, warnings } = validateWithModel(meta, input, 'image_to_image');
    warnings.forEach((warning) => log(`Warning: ${warning}`));
    if (errors.length > 0) {
      return { ok: false, phase: 'failed', errorCode: 'VALIDATION', errorMessage: errors.join(' ') };
    }
  } else if (registry) {
    log(`Warning: model "${model}" not found in image_to_image registry. Proceeding anyway.`);
  }

  const body = meta ? buildBody(meta, input, 'image_to_image') : buildFallbackBody(input, model);

  if (ctx.dryRun) {
    return {
      ok: true,
      phase: 'dry-run',
      dryRun: true,
      requestBody: body,
      requestUrl: `${ctx.baseUrl}/v1/generation/image-to-image`,
    };
  }

  const client = createClient(ctx);
  let res;
  try {
    res = await client.post('/v1/generation/image-to-image', body);
  } catch (err) {
    return formatNetworkError(err);
  }

  if (!isApiSuccess(res)) {
    return formatApiError(res);
  }

  const data = res.data || {};
  const taskIds = data.task_ids ?? (data.task_id ? [data.task_id] : []);
  return {
    ok: true,
    phase: 'submitted',
    batchId: data.batch_id ?? null,
    taskIds,
    taskId: taskIds[0] ?? null,
    taskStatus: null,
    images: null,
    balance: null,
    errorCode: null,
    errorMessage: null,
  };
}

function buildFallbackBody(input, model) {
  const body = {
    prompt: input.prompt,
    model,
    images: input.images.slice(0, 1),
    aspect_ratio: input.aspect_ratio || input.aspectRatio || FALLBACK_DEFAULTS.aspect_ratio,
    image_number: Number(input.image_number || input.imageNumber) || FALLBACK_DEFAULTS.image_number,
  };
  if (input.negative_prompt || input.negativePrompt) {
    body.negative_prompt = input.negative_prompt || input.negativePrompt;
  }
  if (input.resolution) {
    body.resolution = input.resolution;
  }
  return body;
}

export default execute;
