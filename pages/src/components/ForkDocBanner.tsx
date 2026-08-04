import React from 'react';
import { DocSlug } from '../content/docs';

/** Slugs that describe upstream-only features (LLM config, viewer, npm, etc.). */
export const UPSTREAM_ONLY_SLUGS: ReadonlySet<DocSlug> = new Set([
  'quickstart',
  'installation',
  'configuration',
  'cli-reference',
  'architecture',
  'tools',
  'mcp',
  'viewer',
  'telemetry',
  'agent-skill',
  'claude-code',
  'cicd',
  'faq',
  'contributing',
]);

interface ForkDocBannerProps {
  slug: DocSlug;
  language: string;
}

type BannerCopy = {
  title: string;
  upstreamOnly: React.ReactNode;
  forkOnly: React.ReactNode;
};

const upstreamRepo = 'https://github.com/alibaba/open-code-review';
const forkReadme = 'https://github.com/haotool/open-code-review-delegate#quick-install-3-steps';

const linkStyle = { color: '#2BDE5E' };
const codeStyle = { color: '#fff' };

function bannerCopy(language: string, delegateHref: string): BannerCopy {
  switch (language) {
    case 'zh':
      return {
        title: 'Delegate Edition（haotool fork）',
        upstreamOnly: (
          <>
            本页描述上游{' '}
            <a href={upstreamRepo} style={linkStyle} target="_blank" rel="noreferrer">
              alibaba/open-code-review
            </a>{' '}
            （npm 安装、内置 LLM、viewer、遥测），不适用于本 fork。请改用{' '}
            <code style={codeStyle}>ocr-delegate</code> 与 Agent Skill — 参见{' '}
            <a href={delegateHref} style={linkStyle}>
              委托模式
            </a>{' '}
            与{' '}
            <a href={forkReadme} style={linkStyle} target="_blank" rel="noreferrer">
              README
            </a>
            。
          </>
        ),
        forkOnly: (
          <>
            本 fork 提供 <code style={codeStyle}>ocr-delegate</code>（零出站网络工程 CLI）。标准工作流：{' '}
            <a href={delegateHref} style={linkStyle}>
              委托模式
            </a>
            。
          </>
        ),
      };
    case 'ja':
      return {
        title: 'Delegate Edition（haotool fork）',
        upstreamOnly: (
          <>
            このページは上流{' '}
            <a href={upstreamRepo} style={linkStyle} target="_blank" rel="noreferrer">
              alibaba/open-code-review
            </a>{' '}
            （npm インストール、組み込み LLM、viewer、テレメトリ）を説明しており、本 fork には適用されません。代わりに{' '}
            <code style={codeStyle}>ocr-delegate</code> と Agent Skill を使用してください —{' '}
            <a href={delegateHref} style={linkStyle}>
              デリゲーションモード
            </a>{' '}
            と{' '}
            <a href={forkReadme} style={linkStyle} target="_blank" rel="noreferrer">
              README
            </a>{' '}
            を参照。
          </>
        ),
        forkOnly: (
          <>
            本 fork は <code style={codeStyle}>ocr-delegate</code>（外部ネットワーク不要のエンジニアリング CLI）を提供します。標準ワークフロー：{' '}
            <a href={delegateHref} style={linkStyle}>
              デリゲーションモード
            </a>
            。
          </>
        ),
      };
    default:
      return {
        title: 'Delegate Edition (haotool fork)',
        upstreamOnly: (
          <>
            This page documents upstream{' '}
            <a href={upstreamRepo} style={linkStyle} target="_blank" rel="noreferrer">
              alibaba/open-code-review
            </a>{' '}
            (npm install, embedded LLM, viewer, telemetry). It does not apply to this fork.
            Use <code style={codeStyle}>ocr-delegate</code> and the agent skill instead — see{' '}
            <a href={delegateHref} style={linkStyle}>
              Delegation Mode
            </a>{' '}
            and the{' '}
            <a href={forkReadme} style={linkStyle} target="_blank" rel="noreferrer">
              README
            </a>
            .
          </>
        ),
        forkOnly: (
          <>
            This fork ships <code style={codeStyle}>ocr-delegate</code> (zero-network engineering CLI).
            Canonical workflow:{' '}
            <a href={delegateHref} style={linkStyle}>
              Delegation Mode
            </a>
            .
          </>
        ),
      };
  }
}

const ForkDocBanner: React.FC<ForkDocBannerProps> = ({ slug, language }) => {
  if (slug === 'delegate') return null;

  const isUpstreamOnly = UPSTREAM_ONLY_SLUGS.has(slug);
  const delegateHref = '/docs/delegate';
  const copy = bannerCopy(language, delegateHref);

  return (
    <div
      style={{
        marginBottom: 24,
        padding: '12px 16px',
        borderRadius: 8,
        border: '1px solid rgba(43, 222, 94, 0.35)',
        background: 'rgba(43, 222, 94, 0.08)',
        fontSize: 14,
        lineHeight: '22px',
        color: 'rgba(255,255,255,0.85)',
      }}
    >
      <strong style={{ color: '#2BDE5E' }}>{copy.title}</strong>
      {' — '}
      {isUpstreamOnly ? copy.upstreamOnly : copy.forkOnly}
    </div>
  );
};

export default ForkDocBanner;
