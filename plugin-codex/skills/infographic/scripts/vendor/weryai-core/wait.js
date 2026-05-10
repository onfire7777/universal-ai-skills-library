import { log, sleep } from './client.js';
import { hasTaskOutput, normalizeBatchPhase, normalizeTaskCollection, normalizeTaskResult, toPhase } from './normalize.js';
import { isApiSuccess } from './errors.js';

export async function pollSubmittedTasks(client, options) {
  const { batchId, taskIds, allowBatch = false } = options;
  if (!Array.isArray(taskIds) || taskIds.length === 0) {
    return {
      ok: false,
      phase: 'failed',
      errorCode: 'PROTOCOL',
      errorMessage: 'API returned success but no task_ids.',
    };
  }
  if (taskIds.length > 1) {
    return batchId && allowBatch
      ? pollBatch(client, options)
      : pollMultipleTasks(client, options);
  }
  return pollSingleTask(client, { ...options, taskId: taskIds[0] });
}

export async function pollSingleTask(client, options) {
  const { taskId, batchId = null, taskIds = taskId ? [taskId] : [], ctx, outputKey, outputLabel } = options;
  const start = Date.now();

  while (true) {
    const elapsed = Date.now() - start;
    if (elapsed >= ctx.pollTimeoutMs) {
      return {
        ok: false,
        phase: 'failed',
        batchId,
        taskIds,
        taskId,
        taskStatus: 'unknown',
        [outputKey]: null,
        errorCode: 'TIMEOUT',
        errorMessage: `Poll timeout after ${Math.round(elapsed / 1000)}s.`,
      };
    }

    await sleep(ctx.pollIntervalMs);

    let res;
    try {
      res = await client.get(`/v1/generation/${taskId}/status`, { retries: 3 });
    } catch (err) {
      log(`Warning: poll query failed (${err.message}), will retry next interval.`);
      continue;
    }

    if (!isApiSuccess(res)) {
      log(`Warning: poll returned status ${res.status}, will retry next interval.`);
      continue;
    }

    const task = normalizeTaskResult(res.data);
    const phase = toPhase(task.taskStatus);
    const missingOutputs = phase === 'completed' && !hasTaskOutput(task, outputKey);
    const elapsedSec = Math.round((Date.now() - start) / 1000);
    log(`Polling ${taskId}... status: ${task.taskStatus} (${elapsedSec}s elapsed)`);

    if (phase === 'completed' || phase === 'failed') {
      return {
        ok: phase === 'completed' && !missingOutputs,
        phase: missingOutputs ? 'failed' : phase,
        batchId,
        taskIds,
        taskId: task.taskId,
        taskStatus: task.taskStatus,
        [outputKey]: task[outputKey],
        lyrics: task.lyrics,
        coverUrl: task.coverUrl,
        balance: null,
        errorCode: phase === 'failed' || missingOutputs ? 'TASK_FAILED' : null,
        errorMessage: missingOutputs
          ? `Task reached a completed state but returned no ${outputLabel} URLs.`
          : phase === 'failed'
            ? (task.msg || 'Task failed')
            : null,
      };
    }
  }
}

export async function pollMultipleTasks(client, options) {
  const { taskIds, ctx, outputKey, outputLabel } = options;
  const start = Date.now();

  while (true) {
    const elapsed = Date.now() - start;
    if (elapsed >= ctx.pollTimeoutMs) {
      return {
        ok: false,
        phase: 'failed',
        batchId: null,
        taskIds,
        tasks: null,
        errorCode: 'TIMEOUT',
        errorMessage: `Poll timeout after ${Math.round(elapsed / 1000)}s.`,
      };
    }

    await sleep(ctx.pollIntervalMs);

    const tasks = [];
    let shouldRetry = false;

    for (const taskId of taskIds) {
      let res;
      try {
        res = await client.get(`/v1/generation/${taskId}/status`, { retries: 3 });
      } catch (err) {
        log(`Warning: poll query failed for ${taskId} (${err.message}), will retry next interval.`);
        shouldRetry = true;
        break;
      }

      if (!isApiSuccess(res)) {
        log(`Warning: poll returned status ${res.status} for ${taskId}, will retry next interval.`);
        shouldRetry = true;
        break;
      }

      tasks.push(normalizeTaskResult(res.data));
    }

    if (shouldRetry) {
      continue;
    }

    const phase = normalizeBatchPhase(tasks, { outputKey });
    const elapsedSec = Math.round((Date.now() - start) / 1000);
    const summary = tasks.map((task) => `${task.taskId}:${task.taskStatus}`).join(', ');
    log(`Polling tasks... [${summary}] (${elapsedSec}s elapsed)`);

    if (phase === 'completed' || phase === 'failed') {
      return {
        ok: phase === 'completed',
        phase,
        batchId: null,
        taskIds,
        tasks,
        errorCode: phase === 'failed' ? 'TASK_FAILED' : null,
        errorMessage: phase === 'failed'
          ? `One or more tasks failed or returned no ${outputLabel} URLs.`
          : null,
      };
    }
  }
}

export async function pollBatch(client, options) {
  const { batchId, taskIds, ctx, outputKey, outputLabel } = options;
  const start = Date.now();

  while (true) {
    const elapsed = Date.now() - start;
    if (elapsed >= ctx.pollTimeoutMs) {
      return {
        ok: false,
        phase: 'failed',
        batchId,
        taskIds,
        tasks: null,
        errorCode: 'TIMEOUT',
        errorMessage: `Poll timeout after ${Math.round(elapsed / 1000)}s.`,
      };
    }

    await sleep(ctx.pollIntervalMs);

    let res;
    try {
      res = await client.get(`/v1/generation/batch/${batchId}/status`, { retries: 3 });
    } catch (err) {
      log(`Warning: batch poll failed (${err.message}), will retry next interval.`);
      continue;
    }

    if (!isApiSuccess(res)) {
      log(`Warning: batch poll returned status ${res.status}, will retry next interval.`);
      continue;
    }

    const tasks = normalizeTaskCollection(res.data);
    const phase = normalizeBatchPhase(tasks, { outputKey });
    const elapsedSec = Math.round((Date.now() - start) / 1000);
    const summary = tasks.map((task) => task.taskStatus).join(', ');
    log(`Polling batch ${batchId}... [${summary}] (${elapsedSec}s elapsed)`);

    if (phase === 'completed' || phase === 'failed') {
      return {
        ok: phase === 'completed',
        phase,
        batchId,
        taskIds,
        tasks,
        errorCode: phase === 'failed' ? 'TASK_FAILED' : null,
        errorMessage: phase === 'failed'
          ? `One or more tasks in the batch failed or returned no ${outputLabel} URLs.`
          : null,
      };
    }
  }
}
