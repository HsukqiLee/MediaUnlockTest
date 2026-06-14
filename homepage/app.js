(function(){

/* ==========================================
   State
   ========================================== */
const S = {
  os: 'linux',
  paramTab: 'ip',
};

/* ==========================================
   Data
   ========================================== */
const OS_DATA = {
  linux:  { label: 'Linux',      icon: 'linux' },
  macos:  { label: 'macOS',      icon: 'apple' },
  windows:{ label: 'Windows',    icon: 'windows' },
  android:{ label: 'Android',    icon: 'android' },
};

const COMMANDS = {
  test: {
    linux:   'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)',
    macos:   'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)',
    windows: 'irm https://unlock.icmp.ing/scripts/download_test.ps1 | iex',
    android: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh)',
  },
  monitor: {
    linux:   'bash <(curl -Ls unlock.icmp.ing/scripts/monitor.sh) -service',
    macos:   'bash <(curl -Ls unlock.icmp.ing/scripts/monitor.sh) -service',
    windows: 'irm https://unlock.icmp.ing/scripts/download_monitor.ps1 | iex',
    android: 'bash <(curl -Ls unlock.icmp.ing/scripts/monitor.sh) -service',
  },
  migrate: {
    linux:   'bash <(curl -Ls unlock.icmp.ing/scripts/migrate.sh)',
    macos:   'bash <(curl -Ls unlock.icmp.ing/scripts/migrate.sh)',
    windows: 'irm https://unlock.icmp.ing/scripts/download_migrate.ps1 | iex',
    android: 'bash <(curl -Ls unlock.icmp.ing/scripts/migrate.sh)',
  },
};

const PARAMS = {
  ip: [
    { flag: '-m 4', title: '仅 IPv4 测试', desc: '只使用 IPv4 连接进行所有项目测试，适用于只有 IPv4 环境', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 4' },
    { flag: '-m 6', title: '仅 IPv6 测试', desc: '只测试已知支持 IPv6 的项目，不会对不支持 IPv6 的服务进行测试', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 6' },
    { flag: '-f', title: '强制 IPv6 测试', desc: '强制对所有项目使用 IPv6 测试，即使某些服务可能不支持 IPv6', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -f' },
    { flag: '-I [IP/接口]', title: '绑定网络接口', desc: '使用特定 IP 或网络接口进行测试，适合多 IP 环境下指定出口 IP', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -I 192.168.1.100' },
  ],
  network: [
    { flag: '-dns-servers', title: '指定 DNS 服务器', desc: '使用指定的 DNS 服务器进行域名解析，有助于解决 DNS 污染问题', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -dns-servers "1.1.1.1:53"' },
    { flag: '-http-proxy', title: 'HTTP 代理', desc: '通过 HTTP 代理进行测试，支持带认证的代理服务器', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -http-proxy "http://127.0.0.1:1080"' },
    { flag: '-socks-proxy', title: 'SOCKS5 代理', desc: '通过 SOCKS5 代理进行测试，支持带认证的代理服务器', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -socks-proxy "socks5://127.0.0.1:1080"' },
  ],
  debug: [
    { flag: '-debug', title: '调试模式', desc: '输出详细的错误和调试信息，用于排查测试中的问题', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -debug' },
    { flag: '-conc [数值]', title: '并发请求数', desc: '调整同时发送的请求数量，提高检测速度，默认为系统自动选择', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -conc 10' },
    { flag: '-u', title: '检查更新', desc: '检查并获取最新版本的脚本，确保使用最新功能', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -u' },
    { flag: '-v', title: '显示版本', desc: '输出当前脚本的版本信息，便于判断是否需要更新', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -v' },
  ],
};

const PARAM_LABELS = { ip: 'IP 模式', network: '网络参数', debug: '调试选项' };

const SVG_COPY = '<svg viewBox="0 0 24 24"><path d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/></svg>';
const SVG_CHECK = '<svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg>';
/* ==========================================
   DOM Helpers
   ========================================== */
function h(tag, cls, content) {
  const el = document.createElement(tag);
  if (cls) el.className = cls;
  if (content !== undefined) el.textContent = content;
  return el;
}

/* ==========================================
   Ripple
   ========================================== */
function createRipple(e) {
  const el = e.currentTarget;
  const ripple = document.createElement('span');
  ripple.className = 'md-ripple';
  const rect = el.getBoundingClientRect();
  const size = Math.max(rect.width, rect.height);
  ripple.style.width = ripple.style.height = size + 'px';
  ripple.style.left = (e.clientX - rect.left - size / 2) + 'px';
  ripple.style.top = (e.clientY - rect.top - size / 2) + 'px';
  el.appendChild(ripple);
  ripple.addEventListener('animationend', () => ripple.remove());
}

/* ==========================================
   Copy to Clipboard
   ========================================== */
async function copyText(btn, text) {
  try {
    await navigator.clipboard.writeText(text);
    btn.innerHTML = SVG_CHECK;
    btn.style.color = 'var(--md-sys-color-primary)';
    showSnackbar('已复制到剪贴板');
    setTimeout(() => {
      btn.innerHTML = SVG_COPY;
      btn.style.color = '';
    }, 2000);
  } catch {
    showSnackbar('复制失败');
  }
}

/* ==========================================
   Snackbar
   ========================================== */
let snackbarTimer = null;
function showSnackbar(msg) {
  let sb = document.getElementById('snackbar');
  if (!sb) {
    sb = h('div', 'snackbar');
    sb.id = 'snackbar';
    document.body.appendChild(sb);
  }
  sb.textContent = msg;
  sb.classList.add('show');
  clearTimeout(snackbarTimer);
  snackbarTimer = setTimeout(() => sb.classList.remove('show'), 2400);
}

/* ==========================================
   OS Icon
   ========================================== */
function renderOsIcon(key) {
  return '<svg class="chip-icon" viewBox="0 0 24 24" fill="currentColor">' + (
    key === 'linux' ? '<path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>' :
    key === 'apple' ? '<path d="M17.05 20.28c-.98.95-2.05.8-3.08.35-1.09-.46-2.09-.48-3.24 0-1.44.62-2.2.44-3.06-.35C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z"/>' :
    key === 'windows' ? '<path d="M3 12V6.99l6-1.02v6.03H3zm7-6.34l11-1.79v8.13H10V5.66zm-7 6.68h6v6.01l-6-1.02v-4.99zm7 5.68l11-1.79V12H10v5.34z"/>' :
    key === 'android' ? '<path d="M17.523 16.18c-.553 0-1-.447-1-1s.447-1 1-1 1 .447 1 1-.447 1-1 1zm-3.614 0c-.553 0-1-.447-1-1s.447-1 1-1 1 .447 1 1-.447 1-1 1zM5.003 8.273l.014 9.154c0 .672.318.928.934.928h.852l.015-3.413h1.669c2.349 0 3.967.589 3.967 2.21 0 1.776-1.7 2.313-3.818 2.313-1.637 0-2.296-.417-2.434-.994v-1.504h-1.63v2.972c0 1.534 1.619 2.061 3.169 2.061 2.453 0 5.266-.612 5.266-3.606 0-2.129-1.128-3.2-2.966-3.695v-.041c1.14-.335 2.298-1.189 2.298-2.947 0-2.436-1.887-3.238-4.542-3.238H5.003zm1.63 1.732h2.189c1.474 0 2.397.383 2.397 1.515 0 1.168-.887 1.626-2.504 1.626H6.633V10.01zm0 4.666h1.86c1.782 0 2.888.282 2.888 1.602 0 1.389-1.09 1.866-2.739 1.866H6.633v-3.468zM17.451 7.808l2.747-4.757.051-.09h-1.734l-2.144 3.714A6.273 6.273 0 0012.406 6.4c-.649 0-1.276.086-1.874.246l.025-.043 1.862-3.224-.968-1.769H9.654L7.641 5.043C6.904 5.26 6.219 5.585 5.6 6.001l.008-.007H5.59c.326-.202.682-.36 1.059-.469l1.504-2.607h1.734l-1.373 2.379a6.284 6.284 0 00.43.091l.484.126L11.382 2.9h.968L10.63 5.84l.025.043c.34.082.672.183.993.303l.148.052 1.765-3.057h1.734L13.30 6.27c.898.58 1.621 1.389 2.095 2.34l.003-.003s-.001.001-.001.003c.397.83.624 1.738.624 2.693l.009 6.151h1.628l-.014-7.738c0-.452-.064-.897-.193-1.323l.003-.004c.106.038.212.07.318.107.556.191 1.04.513 1.395.954l.035.041h1.736c-1.013-1.906-3.158-3.27-5.629-3.645l-.003-.004z"/>' : ''
  ) + '</svg>';
}

/* ==========================================
   OS Selector Chip
   ========================================== */
function osChip(key, onClick) {
  const active = S.os === key;
  const c = document.createElement('button');
  c.className = 'chip state-layer' + (active ? ' active' : '');
  c.type = 'button';
  c.setAttribute('aria-pressed', String(active));
  c.innerHTML = renderOsIcon(key) + '<span>' + OS_DATA[key].label + '</span>';
  c.addEventListener('click', createRipple);
  c.addEventListener('click', () => { S.os = key; onClick(); });
  return c;
}

/* ==========================================
   Param Filter Chip
   ========================================== */
function paramChip(key, onClick) {
  const active = S.paramTab === key;
  const c = document.createElement('button');
  c.className = 'chip state-layer' + (active ? ' active' : '');
  c.type = 'button';
  c.textContent = PARAM_LABELS[key];
  c.addEventListener('click', createRipple);
  c.addEventListener('click', () => { S.paramTab = key; onClick(); });
  return c;
}

/* ==========================================
   Command Block
   ========================================== */
function commandBlock(cmd) {
  const wrap = h('div', 'command-block');
  const pre = document.createElement('pre');
  pre.textContent = cmd;
  const btn = document.createElement('button');
  btn.className = 'copy-btn state-layer';
  btn.innerHTML = SVG_COPY;
  btn.addEventListener('click', createRipple);
  btn.addEventListener('click', () => copyText(btn, cmd));
  wrap.append(pre, btn);
  return wrap;
}

/* ==========================================
   Render: OS Commands
   ========================================== */
function renderOSCommands(container, cmdKey) {
  container.innerHTML = '';
  const selector = h('div', 'os-selector chip-group');
  ['linux', 'macos', 'windows', 'android'].forEach(key => {
    selector.appendChild(osChip(key, () => renderOSCommands(container, cmdKey)));
  });
  container.appendChild(selector);

  const cmd = COMMANDS[cmdKey][S.os];
  const wrap = h('div', 'command-wrap');
  const label = h('div', 'command-label', '');
  const osIcon = renderOsIcon(S.os);
  label.innerHTML = osIcon + ' ' + OS_DATA[S.os].label;
  wrap.appendChild(label);
  wrap.appendChild(commandBlock(cmd));
  container.appendChild(wrap);
}

/* ==========================================
   Render: Parameters
   ========================================== */
function renderParams(container) {
  container.innerHTML = '';
  const selector = h('div', 'chip-group');
  selector.style.marginBottom = '16px';
  ['ip', 'network', 'debug'].forEach(key => {
    selector.appendChild(paramChip(key, () => renderParams(container)));
  });
  container.appendChild(selector);

  const params = PARAMS[S.paramTab];
  const grid = h('div', 'param-grid');
  params.forEach(p => {
    const card = h('div', 'param-card');
    const header = h('div', 'param-header');
    const name = h('span', 'param-name', p.flag);
    header.appendChild(name);
    const body = h('div', 'param-body');
    const desc = h('div', 'param-desc');
    const strong = document.createElement('strong');
    strong.textContent = p.title;
    desc.appendChild(strong);
    desc.appendChild(document.createTextNode(p.desc));
    body.appendChild(desc);
    body.appendChild(commandBlock(p.cmd));
    card.append(header, body);
    grid.appendChild(card);
  });
  container.appendChild(grid);
}

/* ==========================================
   Top Bar Scroll
   ========================================== */
(function() {
  const tb = document.querySelector('.top-bar');
  if (!tb) return;
  let ticking = false;
  window.addEventListener('scroll', function() {
    if (!ticking) {
      requestAnimationFrame(function() {
        tb.classList.toggle('scrolled', window.scrollY > 0);
        ticking = false;
      });
      ticking = true;
    }
  }, { passive: true });
})();

/* ==========================================
   Init
   ========================================== */
document.addEventListener('DOMContentLoaded', function() {
  // Initialize OS commands
  const testCmdEl = document.getElementById('test-commands');
  const monitorCmdEl = document.getElementById('monitor-commands');
  const migrateCmdEl = document.getElementById('migrate-commands');
  const paramsEl = document.getElementById('param-grid');

  const rerenderAll = () => {
    if (testCmdEl) renderOSCommands(testCmdEl, 'test');
    if (monitorCmdEl) renderOSCommands(monitorCmdEl, 'monitor');
    if (migrateCmdEl) renderOSCommands(migrateCmdEl, 'migrate');
    if (paramsEl) renderParams(paramsEl);
  };

  rerenderAll();

  // Buttons with ripple
  document.querySelectorAll('.btn, .copy-btn, .chip').forEach(el => {
    el.addEventListener('click', createRipple);
  });
});

})();
