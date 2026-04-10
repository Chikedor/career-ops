import { existsSync } from 'fs';
import { join } from 'path';

export function getTrackerFilePaths(baseDir) {
  const primary = join(baseDir, 'data/applications.md');
  const legacy = join(baseDir, 'applications.md');
  const readPath = existsSync(primary) ? primary : (existsSync(legacy) ? legacy : primary);
  return { primary, legacy, readPath };
}
