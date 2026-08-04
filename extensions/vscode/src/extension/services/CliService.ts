import { spawn } from 'child_process';
import { CliResult, CliRunOptions, EnvCheckResult, LogLine } from '../../shared/types';
import { extractCliError, parseLogLine } from './cliParse';
import { getShellEnv, resolveBin } from './shellEnv';

export class CliService {
  private current: ReturnType<typeof spawn> | null = null;
  private envCache: { env: EnvCheckResult; at: number } | null = null;
  private static readonly ENV_CACHE_TTL_MS = 5 * 60 * 1000;

  constructor(private cliPath: string = 'ocr-delegate') {}

  invalidateEnvironmentCache(): void {
    this.envCache = null;
  }

  getCachedEnvironment(): EnvCheckResult | null {
    if (!this.envCache) return null;
    if (Date.now() - this.envCache.at > CliService.ENV_CACHE_TTL_MS) {
      this.envCache = null;
      return null;
    }
    return this.envCache.env;
  }

  async isAvailable(): Promise<boolean> {
    const env = await this.checkEnvironment();
    return env.ocr.ok;
  }

  private probeCommand(bin: string): Promise<{ ok: boolean; version?: string }> {
    return new Promise((resolve) => {
      // shell: true is safe here because args are hardcoded ['--version'] — no user input.
      const proc = spawn(resolveBin(bin), ['--version'], { env: getShellEnv(), shell: process.platform === 'win32' });
      let stdout = '';
      let errored = false;
      proc.stdout?.on('data', (d) => { stdout += d.toString(); });
      proc.on('error', () => { errored = true; resolve({ ok: false }); });
      proc.on('close', (code) => {
        if (errored || code !== 0) {
          resolve({ ok: false });
          return;
        }
        const version = stdout.trim().split('\n')[0]?.trim();
        resolve({ ok: true, version: version || undefined });
      });
    });
  }

  async checkEnvironment(force = false): Promise<EnvCheckResult> {
    if (!force) {
      const cached = this.getCachedEnvironment();
      if (cached) return cached;
    }
    const node = await this.probeCommand('node');
    const npm = node.ok ? await this.probeCommand('npm') : { ok: false };
    // ocr-delegate is probed independently — Delegate Edition does not require npm install.
    const ocr = await this.probeCommand(this.cliPath);
    const env = { node, npm, ocr };
    this.envCache = { env, at: Date.now() };
    return env;
  }

  /** npm グローバルインストールは Delegate Edition では提供されません。 */
  install(onLog: (l: LogLine) => void): Promise<boolean> {
    onLog({
      text: 'npm install is not supported in Delegate Edition — build ocr-delegate from source (make build)',
      level: 'error',
    });
    return Promise.resolve(false);
  }

  /** 运行任意参数，流式回调日志，结束返回 stdout 全文。退出码非 0 时 reject，并带上 CLI 报错文本。 */
  runRaw(
    args: string[],
    cwd: string,
    onLog: (l: LogLine) => void,
    envExtra?: Record<string, string>,
  ): Promise<string> {
    return new Promise((resolve, reject) => {
      const proc = spawn(resolveBin(this.cliPath), args, {
        cwd,
        env: envExtra ? { ...getShellEnv(), ...envExtra } : getShellEnv(),
      });
      this.current = proc;
      let stdout = '';
      let stderr = '';
      proc.stdout.on('data', (d) => { stdout += d.toString(); });
      proc.stderr.on('data', (d) => {
        const text = d.toString();
        stderr += text;
        for (const line of text.split('\n')) {
          const parsed = parseLogLine(line);
          if (parsed) onLog(parsed);
        }
      });
      proc.on('error', (err) => { this.current = null; reject(err); });
      proc.on('close', (code) => {
        this.current = null;
        if (code === 0) { resolve(stdout); return; }
        reject(new Error(extractCliError(stderr) || `CLI exited with code ${code}`));
      });
    });
  }

  async review(opts: CliRunOptions, cwd: string, onLog: (l: LogLine) => void): Promise<CliResult> {
    void opts;
    void cwd;
    const msg =
      'Delegate Edition: full LLM review runs via the host agent skill (ocr-delegate preview/rule). ' +
      'This extension provides file preview only — use your agent to complete the review.';
    onLog({ text: msg, level: 'error' });
    throw new Error(msg);
  }

  async testConnection(options?: { configPath?: string; home?: string }): Promise<{ ok: boolean; message?: string }> {
    const envExtra: Record<string, string> = {};
    if (options?.home) {
      envExtra.HOME = options.home;
      if (process.platform === 'win32') envExtra.USERPROFILE = options.home;
    }
    if (options?.configPath) envExtra.OCR_CONFIG_PATH = options.configPath;
    const env = Object.keys(envExtra).length > 0 ? envExtra : undefined;
    try {
      await this.runRaw(['llm', 'test'], process.cwd(), () => {}, env);
      return { ok: true };
    } catch (e) {
      return { ok: false, message: e instanceof Error ? e.message : String(e) };
    }
  }

  cancel(): void {
    if (this.current && this.current.pid) {
      this.current.kill('SIGTERM');
      const proc = this.current;
      const forceKillTimer = setTimeout(() => {
        if (proc.exitCode === null && proc.signalCode === null) proc.kill('SIGKILL');
      }, 3000);
      proc.once('close', () => clearTimeout(forceKillTimer));
    }
  }
}
