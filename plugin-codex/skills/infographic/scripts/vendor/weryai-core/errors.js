const ERROR_MESSAGES = {
  400: 'Bad request — check your request parameters.',
  401: 'Authentication failed — verify your API key is valid.',
  403: 'Authentication failed or IP access denied — verify your API key and account policy.',
  404: 'The requested WeryAI resource was not found.',
  429: 'Rate limited by WeryAI API. Please wait and try again.',
  500: 'WeryAI server error. Please try again later.',
  1001: 'Parameter error — check required fields and parameter format.',
  1002: 'Authentication failed — verify your API key is valid.',
  1003: 'Resource not found — the task_id or batch_id may not exist or has expired.',
  1004: 'Insufficient credits — recharge at weryai.com before retrying.',
  1006: 'Model not supported — check the model key or try a different model.',
  1007: 'Queue full — the service is busy, please try again later.',
  1009: 'Permission denied.',
  1010: 'Data not found.',
  1011: 'Insufficient credits — recharge at weryai.com before retrying.',
  1014: 'Upload image exceeds max size limit.',
  1015: 'Upload image exceeds daily limit.',
  1017: 'VIP permission required — this feature requires a subscription.',
  2001: 'Content flagged by the safety system.',
  2003: 'Content flagged by the safety system — revise the prompt or input.',
  2004: 'Image format not supported — use jpg, png, or webp.',
  6001: 'Workflow system error — please try again later.',
  6002: 'Workflow rate limit exceeded — please try again later.',
  6003: 'Credit deduction failed — please try again later.',
  6004: 'Generation failed — please try again later.',
  6010: 'Concurrent task limit reached — wait for existing tasks to finish.',
  6101: 'Daily request limit exceeded for this API type.',
};

export function formatApiError(response) {
  const { status, httpStatus } = response;
  const messageText = response.msg || response.message || null;

  if (httpStatus === 401 || httpStatus === 403) {
    return failure(String(httpStatus), appendDetail(ERROR_MESSAGES[httpStatus], messageText));
  }
  if (httpStatus === 404) {
    return failure('404', appendDetail(ERROR_MESSAGES[404], messageText));
  }
  if (httpStatus === 429) {
    return failure('RATE_LIMIT', ERROR_MESSAGES[429]);
  }
  if (httpStatus >= 500) {
    return failure('500', `${ERROR_MESSAGES[500]} (HTTP ${httpStatus})`);
  }
  if (httpStatus === 400) {
    return failure('400', appendDetail(ERROR_MESSAGES[400], messageText));
  }

  const code = status ?? null;
  const friendly = ERROR_MESSAGES[code];
  const message = friendly
    ? appendDetail(friendly, messageText)
    : messageText || `Unknown API error (status: ${code}, HTTP ${httpStatus})`;

  return failure(code ? String(code) : null, message);
}

export function formatNetworkError(err) {
  if (err.code === 'NO_API_KEY') {
    return failure('NO_API_KEY', err.message);
  }
  if (err.code === 'TIMEOUT') {
    return failure('TIMEOUT', err.message);
  }
  return failure('NETWORK_ERROR', `Network error: ${err.message}`);
}

export function isApiSuccess(response) {
  const ok = response.httpStatus >= 200 && response.httpStatus < 300;
  const bodyOk = response.status === 0 || response.status === 200;
  return ok && bodyOk;
}

function appendDetail(base, detail) {
  return detail ? `${base} (${detail})` : base;
}

function failure(errorCode, errorMessage) {
  return { ok: false, phase: 'failed', errorCode, errorMessage };
}
