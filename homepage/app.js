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
    { flag: '-region [编号/名称]', title: '选择检测区域', desc: '非交互式选择检测区域，支持菜单编号或地区名称，多个值用逗号分隔，适合第三方脚本调用', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -region 0,11' },
    { flag: '-m 4', title: '仅 IPv4 测试', desc: '只使用 IPv4 连接进行所有项目测试，适用于只有 IPv4 环境', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 4' },
    { flag: '-m 6', title: '仅 IPv6 测试', desc: '只测试已知支持 IPv6 的项目，不会对不支持 IPv6 的服务进行测试', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -m 6' },
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
    { flag: '-f', title: '强制更新', desc: '配合 -u 使用，强制执行更新（即使已是最新版本）', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -u -f' },
    { flag: '-v', title: '显示版本', desc: '输出当前脚本的版本信息，便于判断是否需要更新', cmd: 'bash <(curl -Ls unlock.icmp.ing/scripts/test.sh) -v' },
    { flag: '-test [名称]', title: '单个测试', desc: '运行单个测试项，支持使用显示名称（如 "Disney+"）或函数名（如 "DisneyPlus"）精确匹配', cmd: './unlock-test -test Disney+' },
  ],
};

const PARAM_LABELS = { ip: 'IP 模式', network: '网络参数', debug: '调试选项' };

const SVG_COPY = '<span class="material-icons-round" style="font-size: 16px;">content_copy</span>';
const SVG_CHECK = '<span class="material-icons-round" style="font-size: 16px; color: var(--color-primary);">check</span>';
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
    key === 'linux' ? '<path d="M12.504 0c-.155 0-.315.008-.48.021-4.226.333-3.105 4.807-3.17 6.298-.076 1.092-.3 1.953-1.05 3.02-.885 1.051-2.127 2.75-2.716 4.521-.278.832-.41 1.684-.287 2.489a.424.424 0 00-.11.135c-.26.268-.45.6-.663.839-.199.199-.485.267-.797.4-.313.136-.658.269-.864.68-.09.189-.136.394-.132.602 0 .199.027.4.055.536.058.399.116.728.04.97-.249.68-.28 1.145-.106 1.484.174.334.535.47.94.601.81.2 1.91.135 2.774.6.926.466 1.866.67 2.616.47.526-.116.97-.464 1.208-.946.587-.003 1.23-.269 2.26-.334.699-.058 1.574.267 2.577.2.025.134.063.198.114.333l.003.003c.391.778 1.113 1.132 1.884 1.071.771-.06 1.592-.536 2.257-1.306.631-.765 1.683-1.084 2.378-1.503.348-.199.629-.469.649-.853.023-.4-.2-.811-.714-1.376v-.097l-.003-.003c-.17-.2-.25-.535-.338-.926-.085-.401-.182-.786-.492-1.046h-.003c-.059-.054-.123-.067-.188-.135a.357.357 0 00-.19-.064c.431-1.278.264-2.55-.173-3.694-.533-1.41-1.465-2.638-2.175-3.483-.796-1.005-1.576-1.957-1.56-3.368.026-2.152.236-6.133-3.544-6.139zm.529 3.405h.013c.213 0 .396.062.584.198.19.135.33.332.438.533.105.259.158.459.166.724 0-.02.006-.04.006-.06v.105a.086.086 0 01-.004-.021l-.004-.024a1.807 1.807 0 01-.15.706.953.953 0 01-.213.335.71.71 0 00-.088-.042c-.104-.045-.198-.064-.284-.133a1.312 1.312 0 00-.22-.066c.05-.06.146-.133.183-.198.053-.128.082-.264.088-.402v-.02a1.21 1.21 0 00-.061-.4c-.045-.134-.101-.2-.183-.333-.084-.066-.167-.132-.267-.132h-.016c-.093 0-.176.03-.262.132a.8.8 0 00-.205.334 1.18 1.18 0 00-.09.4v.019c.002.089.008.179.02.267-.193-.067-.438-.135-.607-.202a1.635 1.635 0 01-.018-.2v-.02a1.772 1.772 0 01.15-.768c.082-.22.232-.406.43-.533a.985.985 0 01.594-.2zm-2.962.059h.036c.142 0 .27.048.399.135.146.129.264.288.344.465.09.199.14.4.153.667v.004c.007.134.006.2-.002.266v.08c-.03.007-.056.018-.083.024-.152.055-.274.135-.393.2.012-.09.013-.18.003-.267v-.015c-.012-.133-.04-.2-.082-.333a.613.613 0 00-.166-.267.248.248 0 00-.183-.064h-.021c-.071.006-.13.04-.186.132a.552.552 0 00-.12.27.944.944 0 00-.023.33v.015c.012.135.037.2.08.334.046.134.098.2.166.268.01.009.02.018.034.024-.07.057-.117.07-.176.136a.304.304 0 01-.131.068 2.62 2.62 0 01-.275-.402 1.772 1.772 0 01-.155-.667 1.759 1.759 0 01.08-.668 1.43 1.43 0 01.283-.535c.128-.133.26-.2.418-.2zm1.37 1.706c.332 0 .733.065 1.216.399.293.2.523.269 1.052.468h.003c.255.136.405.266.478.399v-.131a.571.571 0 01.016.47c-.123.31-.516.643-1.063.842v.002c-.268.135-.501.333-.775.465-.276.135-.588.292-1.012.267a1.139 1.139 0 01-.448-.067 3.566 3.566 0 01-.322-.198c-.195-.135-.363-.332-.612-.465v-.005h-.005c-.4-.246-.616-.512-.686-.71-.07-.268-.005-.47.193-.6.224-.135.38-.271.483-.336.104-.074.143-.102.176-.131h.002v-.003c.169-.202.436-.47.839-.601.139-.036.294-.065.466-.065zm2.8 2.142c.358 1.417 1.196 3.475 1.735 4.473.286.534.855 1.659 1.102 3.024.156-.005.33.018.513.064.646-1.671-.546-3.467-1.089-3.966-.22-.2-.232-.335-.123-.335.59.534 1.365 1.572 1.646 2.757.13.535.16 1.104.021 1.67.067.028.135.06.205.067 1.032.534 1.413.938 1.23 1.537v-.043c-.06-.003-.12 0-.18 0h-.016c.151-.467-.182-.825-1.065-1.224-.915-.4-1.646-.336-1.77.465-.008.043-.013.066-.018.135-.068.023-.139.053-.209.064-.43.268-.662.669-.793 1.187-.13.533-.17 1.156-.205 1.869v.003c-.02.334-.17.838-.319 1.35-1.5 1.072-3.58 1.538-5.348.334a2.645 2.645 0 00-.402-.533 1.45 1.45 0 00-.275-.333c.182 0 .338-.03.465-.067a.615.615 0 00.314-.334c.108-.267 0-.697-.345-1.163-.345-.467-.931-.995-1.788-1.521-.63-.4-.986-.87-1.15-1.396-.165-.534-.143-1.085-.015-1.645.245-1.07.873-2.11 1.274-2.763.107-.065.037.135-.408.974-.396.751-1.14 2.497-.122 3.854a8.123 8.123 0 01.647-2.876c.564-1.278 1.743-3.504 1.836-5.268.048.036.217.135.289.202.218.133.38.333.59.465.21.201.477.335.876.335.039.003.075.006.11.006.412 0 .73-.134.997-.268.29-.134.52-.334.74-.4h.005c.467-.135.835-.402 1.044-.7zm2.185 8.958c.037.6.343 1.245.882 1.377.588.134 1.434-.333 1.791-.765l.211-.01c.315-.007.577.01.847.268l.003.003c.208.199.305.53.391.876.085.4.154.78.409 1.066.486.527.645.906.636 1.14l.003-.007v.018l-.003-.012c-.015.262-.185.396-.498.595-.63.401-1.746.712-2.457 1.57-.618.737-1.37 1.14-2.036 1.191-.664.053-1.237-.2-1.574-.898l-.005-.003c-.21-.4-.12-1.025.056-1.69.176-.668.428-1.344.463-1.897.037-.714.076-1.335.195-1.814.12-.465.308-.797.641-.984l.045-.022zm-10.814.049h.01c.053 0 .105.005.157.014.376.055.706.333 1.023.752l.91 1.664.003.003c.243.533.754 1.064 1.189 1.637.434.598.77 1.131.729 1.57v.006c-.057.744-.48 1.148-1.125 1.294-.645.135-1.52.002-2.395-.464-.968-.536-2.118-.469-2.857-.602-.369-.066-.61-.2-.723-.4-.11-.2-.113-.602.123-1.23v-.004l.002-.003c.117-.334.03-.752-.027-1.118-.055-.401-.083-.71.043-.94.16-.334.396-.4.69-.533.294-.135.64-.202.915-.47h.002v-.002c.256-.268.445-.601.668-.838.19-.201.38-.336.663-.336zm7.159-9.074c-.435.201-.945.535-1.488.535-.542 0-.97-.267-1.28-.466-.154-.134-.28-.268-.373-.335-.164-.134-.144-.333-.074-.333.109.016.129.134.199.2.096.066.215.2.36.333.292.2.68.467 1.167.467.485 0 1.053-.267 1.398-.466.195-.135.445-.334.648-.467.156-.136.149-.267.279-.267.128.016.034.134-.147.332a8.097 8.097 0 01-.69.468zm-1.082-1.583V5.64c-.006-.02.013-.042.029-.05.074-.043.18-.027.26.004.063 0 .16.067.15.135-.006.049-.085.066-.135.066-.055 0-.092-.043-.141-.068-.052-.018-.146-.008-.163-.065zm-.551 0c-.02.058-.113.049-.166.066-.047.025-.086.068-.14.068-.05 0-.13-.02-.136-.068-.01-.066.088-.133.15-.133.08-.031.184-.047.259-.005.019.009.036.03.03.05v.02h.003z"/>' :
    key === 'macos' ? '<path d="M12.152 6.896c-.948 0-2.415-1.078-3.96-1.04-2.04.027-3.91 1.183-4.961 3.014-2.117 3.675-.546 9.103 1.519 12.09 1.013 1.454 2.208 3.09 3.792 3.039 1.52-.065 2.09-.987 3.935-.987 1.831 0 2.35.987 3.96.948 1.637-.026 2.676-1.48 3.676-2.948 1.156-1.688 1.636-3.325 1.662-3.415-.039-.013-3.182-1.221-3.22-4.857-.026-3.04 2.48-4.494 2.597-4.559-1.429-2.09-3.623-2.324-4.39-2.376-2-.156-3.675 1.09-4.61 1.09zM15.53 3.83c.843-1.012 1.4-2.427 1.245-3.83-1.207.052-2.662.805-3.532 1.818-.78.896-1.454 2.338-1.273 3.714 1.338.104 2.715-.688 3.559-1.701"/>' :
    key === 'windows' ? '<path d="M0 3.449L9.75 2.1v9.451H0m10.949-9.602L24 0v11.4H10.949M0 12.6h9.75v9.451L0 20.699M10.949 12.6H24V24l-12.9-1.801"/>' :
    key === 'android' ? '<path d="M18.4395 5.5586c-.675 1.1664-1.352 2.3318-2.0274 3.498-.0366-.0155-.0742-.0286-.1113-.043-1.8249-.6957-3.484-.8-4.42-.787-1.8551.0185-3.3544.4643-4.2597.8203-.084-.1494-1.7526-3.021-2.0215-3.4864a1.1451 1.1451 0 0 0-.1406-.1914c-.3312-.364-.9054-.4859-1.379-.203-.475.282-.7136.9361-.3886 1.5019 1.9466 3.3696-.0966-.2158 1.9473 3.3593.0172.031-.4946.2642-1.3926 1.0177C2.8987 12.176.452 14.772 0 18.9902h24c-.119-1.1108-.3686-2.099-.7461-3.0683-.7438-1.9118-1.8435-3.2928-2.7402-4.1836a12.1048 12.1048 0 0 0-2.1309-1.6875c.6594-1.122 1.312-2.2559 1.9649-3.3848.2077-.3615.1886-.7956-.0079-1.1191a1.1001 1.1001 0 0 0-.8515-.5332c-.5225-.0536-.9392.3128-1.0488.5449zm-.0391 8.461c.3944.5926.324 1.3306-.1563 1.6503-.4799.3197-1.188.0985-1.582-.4941-.3944-.5927-.324-1.3307.1563-1.6504.4727-.315 1.1812-.1086 1.582.4941zM7.207 13.5273c.4803.3197.5506 1.0577.1563 1.6504-.394.5926-1.1038.8138-1.584.4941-.48-.3197-.5503-1.0577-.1563-1.6504.4008-.6021 1.1087-.8106 1.584-.4941z"/>' : ''
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
