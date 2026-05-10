const STATUS_MAP = {
  waiting: 'waiting',
  WAITING: 'waiting',
  pending: 'waiting',
  PENDING: 'waiting',
  processing: 'processing',
  PROCESSING: 'processing',
  running: 'processing',
  RUNNING: 'processing',
  succeed: 'completed',
  SUCCEED: 'completed',
  success: 'completed',
  SUCCESS: 'completed',
  failed: 'failed',
  FAILED: 'failed',
};

export function normalizeStatus(raw) {
  const mapped = STATUS_MAP[raw];
  if (!mapped) {
    process.stderr.write(`[weryai] Warning: unknown task_status "${raw}", treating as "unknown"\n`);
    return 'unknown';
  }
  return mapped;
}

export function toPhase(normalizedStatus) {
  if (normalizedStatus === 'completed') return 'completed';
  if (normalizedStatus === 'failed') return 'failed';
  return 'running';
}

export function normalizeTaskResult(data) {
  if (Array.isArray(data)) {
    return pickSingleTask(data.map(normalizeSingleTask));
  }
  return normalizeSingleTask(data);
}

export function normalizeTaskCollection(data) {
  if (Array.isArray(data)) {
    return data.map(normalizeSingleTask);
  }
  return [normalizeSingleTask(data)];
}

export function hasTaskOutput(task, outputKey) {
  const value = task?.[outputKey];
  return Array.isArray(value) ? value.length > 0 : Boolean(value);
}

export function normalizeBatchPhase(tasks, { outputKey } = {}) {
  const valid = tasks.filter(Boolean);
  if (valid.length === 0) return 'unknown';
  const allDone = valid.every((task) => task.taskStatus === 'completed' || task.taskStatus === 'failed');
  if (!allDone) return 'running';
  const anyFailed = valid.some((task) => task.taskStatus === 'failed');
  const anyMissingOutput = outputKey
    ? valid.some((task) => task.taskStatus === 'completed' && !hasTaskOutput(task, outputKey))
    : false;
  if (anyFailed || anyMissingOutput) return 'failed';
  return 'completed';
}

function normalizeSingleTask(task) {
  if (!task) {
    return {
      taskId: null,
      rawStatus: 'unknown',
      taskStatus: 'unknown',
      images: null,
      videos: null,
      audios: null,
      lyrics: null,
      coverUrl: null,
      msg: null,
    };
  }

  const rawStatus = task.task_status ?? task.taskStatus ?? 'unknown';
  const status = normalizeStatus(rawStatus);
  const images = normalizeArray(task.images ?? task.output?.images);
  const videos = normalizeArray(task.videos ?? task.output?.videos);
  const audios = normalizeArray(task.audios ?? task.output?.audios);
  const lyrics = normalizeString(task.lyrics ?? task.output?.lyrics);
  const coverUrl = normalizeString(task.cover_url ?? task.coverUrl ?? task.output?.cover_url);

  return {
    taskId: task.task_id ?? task.taskId ?? null,
    rawStatus,
    taskStatus: status,
    images,
    videos,
    audios,
    lyrics,
    coverUrl,
    msg: task.msg ?? null,
  };
}

function normalizeArray(value) {
  return Array.isArray(value) && value.length > 0 ? value : null;
}

function normalizeString(value) {
  return typeof value === 'string' && value.trim() ? value : null;
}

function pickSingleTask(tasks) {
  if (!Array.isArray(tasks) || tasks.length === 0) {
    return normalizeSingleTask(null);
  }
  return tasks[0];
}
