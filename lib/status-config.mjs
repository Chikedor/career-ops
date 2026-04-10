#!/usr/bin/env node

import { existsSync, readFileSync } from 'fs';
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const ROOT = dirname(dirname(fileURLToPath(import.meta.url)));
const STATES_FILE = join(ROOT, 'templates', 'states.yml');

const LEGACY_ALIAS_MAP = {
  hold: 'evaluated',
  condicional: 'evaluated',
  evaluar: 'evaluated',
  verificar: 'evaluated',
  applied: 'applied',
  sent: 'applied',
  'geo blocker': 'skip',
  duplicado: 'discarded',
  dup: 'discarded',
  repost: 'discarded',
};

function parseInlineList(value) {
  const trimmed = value.trim();
  if (!trimmed.startsWith('[') || !trimmed.endsWith(']')) {
    return [];
  }
  return trimmed
    .slice(1, -1)
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean);
}

function parseStatesFile() {
  if (!existsSync(STATES_FILE)) {
    throw new Error(`Missing states file: ${STATES_FILE}`);
  }

  const lines = readFileSync(STATES_FILE, 'utf-8').split('\n');
  const states = [];
  let current = null;

  for (const rawLine of lines) {
    const line = rawLine.trim();
    if (!line || line.startsWith('#') || line === 'states:') {
      continue;
    }

    if (line.startsWith('- id:')) {
      if (current) states.push(current);
      current = {
        id: line.slice(5).trim(),
        label: '',
        aliases: [],
        dashboardGroup: '',
        rank: 999,
        actionable: false,
        topFilter: false,
      };
      continue;
    }

    if (!current) continue;

    if (line.startsWith('label:')) {
      current.label = line.slice(6).trim();
    } else if (line.startsWith('aliases:')) {
      current.aliases = parseInlineList(line.slice(8));
    } else if (line.startsWith('dashboard_group:')) {
      current.dashboardGroup = line.slice(16).trim();
    } else if (line.startsWith('rank:')) {
      current.rank = Number.parseInt(line.slice(5).trim(), 10);
    } else if (line.startsWith('actionable:')) {
      current.actionable = line.slice(11).trim().toLowerCase() === 'true';
    } else if (line.startsWith('top_filter:')) {
      current.topFilter = line.slice(11).trim().toLowerCase() === 'true';
    }
  }

  if (current) states.push(current);

  const byId = new Map();
  const byLabel = new Map();
  const aliases = new Map();

  for (const state of states) {
    if (!state.id || !state.label) {
      continue;
    }
    byId.set(state.id, state);
    byLabel.set(state.label.toLowerCase(), state);
    aliases.set(state.id.toLowerCase(), state);
    aliases.set(state.label.toLowerCase(), state);
    for (const alias of state.aliases) {
      aliases.set(alias.toLowerCase(), state);
    }
  }

  for (const [alias, targetId] of Object.entries(LEGACY_ALIAS_MAP)) {
    const state = byId.get(targetId);
    if (state) {
      aliases.set(alias, state);
    }
  }

  return {
    states,
    byId,
    byLabel,
    aliases,
  };
}

const statusConfig = parseStatesFile();

export function stripStatusDecorations(raw) {
  return raw
    .replace(/\*\*/g, '')
    .replace(/\s+\d{4}-\d{2}-\d{2}.*$/, '')
    .trim();
}

export function getCanonicalState(raw, fallbackId = 'evaluated') {
  const clean = stripStatusDecorations(raw).toLowerCase();
  const state = statusConfig.aliases.get(clean);
  if (state) {
    return state;
  }
  return statusConfig.byId.get(fallbackId);
}

export function normalizeStatusId(raw, fallbackId = 'evaluated') {
  return getCanonicalState(raw, fallbackId)?.id ?? fallbackId;
}

export function normalizeStatusLabel(raw, fallbackId = 'evaluated') {
  return getCanonicalState(raw, fallbackId)?.label ?? statusConfig.byId.get(fallbackId)?.label ?? 'Evaluated';
}

export function getStatusById(id) {
  return statusConfig.byId.get(id);
}

export function getStatusConfig() {
  return statusConfig;
}
