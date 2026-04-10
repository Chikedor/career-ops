#!/usr/bin/env node

import { readFileSync } from 'fs';
import { spawn } from 'child_process';
import { dirname } from 'path';
import { fileURLToPath } from 'url';

const BATCH_DIR = dirname(fileURLToPath(import.meta.url));
const PROJECT_DIR = dirname(BATCH_DIR);

const promptFile = process.argv[2];
const workerCommand = process.argv[3] || process.env.CAREER_OPS_BATCH_WORKER;
const inputMode = (process.env.CAREER_OPS_BATCH_WORKER_INPUT || 'stdin').toLowerCase();

if (!promptFile) {
  console.error('Missing prompt file argument');
  process.exit(1);
}
if (!workerCommand) {
  console.error('Missing worker command. Set CAREER_OPS_BATCH_WORKER or pass it explicitly.');
  process.exit(1);
}

const prompt = readFileSync(promptFile, 'utf-8');

let command = workerCommand;
if (inputMode === 'arg') {
  const escaped = JSON.stringify(prompt);
  command = `${workerCommand} ${escaped}`;
}

const child = spawn(command, {
  cwd: PROJECT_DIR,
  shell: true,
  stdio: ['pipe', 'pipe', 'pipe'],
});

child.stdout.on('data', (chunk) => process.stdout.write(chunk));
child.stderr.on('data', (chunk) => process.stderr.write(chunk));
child.on('error', (error) => {
  console.error(error.message);
  process.exit(1);
});
child.on('close', (code) => {
  process.exit(code ?? 1);
});

if (inputMode !== 'arg') {
  child.stdin.write(prompt);
}
child.stdin.end();
