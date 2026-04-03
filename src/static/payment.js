

// DOM helpers
const $ = (id) => document.getElementById(id);
const $$ = (sel) => document.querySelectorAll(sel);

// ---- Network Tag ----
const _CHAIN_PRESETS = {
  Ethereum: { bg: 'hsla(0,0%,55%,0.25)',      color: 'hsla(0,0%,55%,1)',      icon: 'ethereum' },
  Solana:   { bg: 'hsla(256,85%,65%,0.25)',   color: 'hsla(256,85%,65%,1)',   icon: 'solana' },
  BSC:      { bg: 'hsla(46,91%,49%,0.25)',    color: 'hsla(46,91%,49%,1)',    icon: 'bsc' },
  Arbitrum: { bg: 'hsla(202,100%,54%,0.25)',  color: 'hsla(202,100%,54%,1)',  icon: 'arbitrum' },
  TRON:     { bg: 'hsla(350,100%,47%,0.20)',  color: 'hsla(350,100%,47%,1)',  icon: 'tron' },
  Polygon:  { bg: 'hsla(263,73%,56%,0.25)',   color: 'hsla(263,73%,56%,1)',   icon: 'polygon' },
  Base:     { bg: 'hsla(240,100%,61%,0.25)',  color: 'hsla(240,100%,61%,1)',  icon: 'base' },
  OP:       { bg: 'hsla(353,99%,51%,0.25)',   color: 'hsla(353,99%,51%,1)',   icon: 'op' },
  HyperEVM: { bg: 'hsla(166,94%,79%,0.20)',  color: 'hsla(166,94%,79%,1)',   icon: 'hyperevm' },
  Plasma:   { bg: 'hsla(166,64%,32%,0.25)',   color: 'hsla(166,64%,32%,1)',   icon: 'plasma' },
  Bitcoin:  { bg: 'hsla(33,93%,54%,0.20)',    color: 'hsla(33,93%,54%,1)',    icon: 'bitcoin' },
  Binance:  { bg: 'hsla(46,91%,49%,0.24)',    color: 'hsla(46,91%,49%,1)',    icon: 'binance' },
};
const _CHAIN_ALIASES = {
  evm:'Ethereum', eth:'Ethereum', ethereum:'Ethereum',
  bsc:'BSC', bnb:'BSC',
  arbitrum:'Arbitrum', arb:'Arbitrum',
  sol:'Solana', solana:'Solana',
  tron:'TRON', trx:'TRON',
  polygon:'Polygon', matic:'Polygon', pol:'Polygon',
  base:'Base',
  op:'OP', optimism:'OP',
  hyperevm:'HyperEVM', hyperliquid:'HyperEVM', hype:'HyperEVM',
  plasma:'Plasma', xpl:'Plasma',
  bitcoin:'Bitcoin', btc:'Bitcoin',
  binance:'Binance', bnb_chain:'Binance',
};
function _resolveChain(name) {
  if (!name) return null;
  const k = name.trim().toLowerCase();
  const canonical = _CHAIN_ALIASES[k] || Object.keys(_CHAIN_PRESETS).find(c => c.toLowerCase() === k);
  return canonical ? { label: canonical, ..._CHAIN_PRESETS[canonical] } : null;
}
function networkTag(network) {
  const chain = _resolveChain(network);
  if (!chain) return `<span style="font-weight:600;color:var(--card-foreground)">${network}</span>`;
  return `<span style="display:inline-flex;align-items:center;gap:4px;padding:2px 8px 2px 4px;border-radius:99px;background:${chain.bg};color:${chain.color};font-size:12px;font-weight:600;line-height:1.6"><img src="/static/images/${chain.icon}.png" width="14" height="14" style="width:14px;height:14px;border-radius:50%;object-fit:cover;flex-shrink:0" />${chain.label}</span>`;
}

// Initialize Lucide icons
if (typeof lucide !== 'undefined') lucide.createIcons();

// Theme
let currentTheme = localStorage.getItem('theme') || 'dark';
applyTheme(currentTheme);

function applyTheme(t) {
  currentTheme = t;
  document.documentElement.setAttribute('data-theme', t);
  localStorage.setItem('theme', t);
  $('icon-moon').style.display = t === 'light' ? 'block' : 'none';
  $('icon-sun').style.display  = t === 'dark'  ? 'block' : 'none';
}

function toggleTheme(e) {
  const next = currentTheme === 'light' ? 'dark' : 'light';
  if (!document.startViewTransition) { applyTheme(next); return; }
  const x = e ? e.clientX : window.innerWidth / 2;
  const y = e ? e.clientY : window.innerHeight / 2;
  const r = Math.hypot(Math.max(x, window.innerWidth - x), Math.max(y, window.innerHeight - y));
  const transition = document.startViewTransition(() => applyTheme(next));
  transition.ready.then(() => {
    const clip = [`circle(0px at ${x}px ${y}px)`, `circle(${r}px at ${x}px ${y}px)`];
    document.documentElement.animate(
      { clipPath: clip },
      { duration: 420, easing: 'ease-in-out', pseudoElement: '::view-transition-new(root)' }
    );
  });
}

// i18n
const LANGS = {
  en: {
    scan_title:          'Scan or copy address to pay',
    amount_to_pay:       'Amount to pay',
    payment_address:     'Payment address',
    i_have_transferred:  'I have transferred',
    checking_blockchain: 'Checking blockchain…',
    expires:             'Expires',
    copied:              'Copied to clipboard',
    verifying:           'Verifying…',
    payment_success:     'Payment Successful',
    redirecting:         'Redirecting…',
    payment_expired:     'Payment Expired',
    expired_sub:         'Please initiate a new payment',
    order_id:            'Order ID',
    order_amount:        'Order amount',
    network_timeout:     'Connection Timeout',
    timeout_sub:         'Unable to connect to the payment server',
    retry:               'Retry',
    back:                'Back',
    order_not_found:     'Order Not Found',
    not_found_sub:       'The order does not exist or has already expired',
  },
  zh: {
    scan_title:          '扫码或复制地址付款',
    amount_to_pay:       '付款金额',
    payment_address:     '付款地址',
    i_have_transferred:  '我已转账',
    checking_blockchain: '链上核验中…',
    expires:             '到期时间',
    copied:              '已复制到剪贴板',
    verifying:           '核验中…',
    payment_success:     '支付成功',
    redirecting:         '正在跳转…',
    payment_expired:     '支付已过期',
    expired_sub:         '请重新发起支付',
    order_id:            '订单 ID',
    order_amount:        '订单金额',
    network_timeout:     '连接超时',
    timeout_sub:         '无法连接至支付服务器',
    retry:               '重试',
    back:                '返回',
    order_not_found:     '订单不存在',
    not_found_sub:       '待支付订单不存在或已过期',
  },
  ja: {
    scan_title:          'アドレスをスキャンまたはコピーして支払う',
    amount_to_pay:       '支払い金額',
    payment_address:     '支払いアドレス',
    i_have_transferred:  '送金しました',
    checking_blockchain: 'ブロックチェーン確認中…',
    expires:             '有効期限',
    copied:              'コピーしました',
    verifying:           '確認中…',
    payment_success:     '支払い完了',
    redirecting:         'リダイレクト中…',
    payment_expired:     '支払い期限切れ',
    expired_sub:         '新しい支払いを開始してください',
    order_id:            '注文 ID',
    order_amount:        '注文金額',
    network_timeout:     '接続タイムアウト',
    timeout_sub:         '支払いサーバーに接続できません',
    retry:               '再試行',
    back:                '戻る',
    order_not_found:     '注文が見つかりません',
    not_found_sub:       '注文が存在しないか、すでに期限切れです',
  },
  ko: {
    scan_title:          '주소를 스캔하거나 복사하여 결제',
    amount_to_pay:       '결제 금액',
    payment_address:     '결제 주소',
    i_have_transferred:  '이체 완료',
    checking_blockchain: '블록체인 확인 중…',
    expires:             '만료 시간',
    copied:              '복사됨',
    verifying:           '확인 중…',
    payment_success:     '결제 성공',
    redirecting:         '리다이렉트 중…',
    payment_expired:     '결제 만료',
    expired_sub:         '새로운 결제를 시작하세요',
    order_id:            '주문 ID',
    order_amount:        '주문 금액',
    network_timeout:     '연결 시간 초과',
    timeout_sub:         '결제 서버에 연결할 수 없습니다',
    retry:               '다시 시도',
    back:                '돌아가기',
    order_not_found:     '주문 없음',
    not_found_sub:       '주문이 존재하지 않거나 이미 만료되었습니다',
  },
  'zh-hk': {
    scan_title:          '掃碼或複製地址付款',
    amount_to_pay:       '付款金額',
    payment_address:     '付款地址',
    i_have_transferred:  '我已轉帳',
    checking_blockchain: '鏈上核驗中…',
    expires:             '到期時間',
    copied:              '已複製到剪貼簿',
    verifying:           '核驗中…',
    payment_success:     '支付成功',
    redirecting:         '正在跳轉…',
    payment_expired:     '支付已過期',
    expired_sub:         '請重新發起支付',
    order_id:            '訂單 ID',
    order_amount:        '訂單金額',
    network_timeout:     '連線逾時',
    timeout_sub:         '無法連線至支付伺服器',
    retry:               '重試',
    back:                '返回',
    order_not_found:     '訂單不存在',
    not_found_sub:       '待支付訂單不存在或已過期',
  },
  ru: {
    scan_title:          'Отсканируйте или скопируйте адрес для оплаты',
    amount_to_pay:       'Сумма к оплате',
    payment_address:     'Адрес для оплаты',
    i_have_transferred:  'Я перевёл',
    checking_blockchain: 'Проверка блокчейна…',
    expires:             'Истекает',
    copied:              'Скопировано',
    verifying:           'Проверка…',
    payment_success:     'Оплата прошла успешно',
    redirecting:         'Перенаправление…',
    payment_expired:     'Срок оплаты истёк',
    expired_sub:         'Пожалуйста, создайте новый платёж',
    order_id:            'Заказ ID',
    order_amount:        'Сумма заказа',
    network_timeout:     'Тайм-аут подключения',
    timeout_sub:         'Не удалось подключиться к серверу',
    retry:               'Повторить',
    back:                'Назад',
    order_not_found:     'Заказ не найден',
    not_found_sub:       'Заказ не существует или уже истёк',
  },
};

const LANG_LABELS = { en: 'EN', zh: '中文', 'zh-hk': '繁體', ja: '日本語', ko: '한국어', ru: 'RU' };
const SUPPORTED_LANGS = Object.keys(LANGS); // ['en','zh','ja','ko','zh-hk','ru']
let currentLang = 'en';
const t = (key) => LANGS[currentLang]?.[key] ?? key;

function detectLang() {
  const saved = localStorage.getItem('lang');
  if (saved && SUPPORTED_LANGS.includes(saved)) return saved;
  for (const nav of (navigator.languages || [navigator.language])) {
    const lc = nav.toLowerCase();
    if (lc.startsWith('zh-hk') || lc.startsWith('zh-tw') || lc.startsWith('zh-mo')) return 'zh-hk';
    if (lc.startsWith('zh')) return 'zh';
    if (lc.startsWith('ja')) return 'ja';
    if (lc.startsWith('ko')) return 'ko';
    if (lc.startsWith('ru')) return 'ru';
    if (lc.startsWith('en')) return 'en';
  }
  return 'en';
}

function setLang(lang) {
  currentLang = lang;
  localStorage.setItem('lang', lang);
  $$('[data-i18n]').forEach(el => {
    const v = t(el.dataset.i18n);
    if (v) el.textContent = v;
  });
  $('lang-label').textContent = LANG_LABELS[lang];
  $$('#dd-lang-menu .select-option').forEach(o =>
    o.classList.toggle('is-selected', o.dataset.lang === lang));
  closeAllSelects();
  if (ORDER.tradeId) { renderOrderId(); }
  if (ORDER.amount && !ORDER.amount.startsWith('{{')) {
    renderRow('display-fiat', t('order_amount'), `${ORDER.amount} ${ORDER.currency || ''}`);
  }
  requestAnimationFrame(syncStatusCardHeight);
}

// Select dropdown
function toggleSelect(id) {
  const wrap = $(id);
  const trigger = wrap.querySelector('.select-trigger');
  const menu = wrap.querySelector('.select-menu');
  const wasOpen = menu.classList.contains('is-open');
  closeAllSelects();
  if (!wasOpen) {
    menu.classList.add('is-open');
    trigger.classList.add('is-open');
  }
}

function closeAllSelects() {
  $$('.select-menu.is-open').forEach(m => {
    m.classList.remove('is-open');
    m.closest('.select-wrap').querySelector('.select-trigger').classList.remove('is-open');
  });
}

document.addEventListener('click', (e) => {
  if (!e.target.closest('.select-wrap')) closeAllSelects();
});

// Clipboard fallback (used when ClipboardJS is unavailable)
const CHECK_ICON = '<i data-lucide="check" width="15" height="15" stroke-width="2.5" color="#22c55e"></i>';

function showToast() {
  const el = $('toast');
  el.textContent = t('copied');
  el.style.opacity = '1';
  el.style.transform = 'translateX(-50%) translateY(0)';
  clearTimeout(el._tid);
  el._tid = setTimeout(() => {
    el.style.opacity = '0';
    el.style.transform = 'translateX(-50%) translateY(10px)';
  }, 2000);
}

function flashCheck(btnId) {
  const btn = $(btnId);
  if (!btn) return;
  const orig = btn.innerHTML;
  btn.innerHTML = CHECK_ICON;
  lucide.createIcons({ nodes: [btn] });
  setTimeout(() => { btn.innerHTML = orig; lucide.createIcons({ nodes: [btn] }); }, 1800);
}

// Transfer button
function handleTransfer() {
  const btn = $('btn-transferred');
  const span = btn.querySelector('span');
  span.textContent = t('verifying');
  btn.disabled = true;
  _pollErrors = 0;
  _pollStopped = false;
  clearTimeout(_pollTimer);
  checkOrderStatus();
  setTimeout(() => {
    span.textContent = t('i_have_transferred');
    btn.disabled = false;
  }, 4000);
}

// Order
function truncateAddr(addr) {
  return addr.length > 10 ? `${addr.slice(0, 4)}...${addr.slice(-6)}` : addr;
}

function formatAddr(addr) {
  if (!addr || addr.length <= 10) return addr;
  const pre = addr.slice(0, 4);
  const mid = addr.slice(4, -6);
  const suf = addr.slice(-6);
  return `<span style="font-weight:800;color:var(--primary)">${pre}</span><span style="color:var(--muted-foreground)">${mid}</span><span style="font-weight:800;color:var(--primary)">${suf}</span>`;
}

function renderRow(id, label, value) {
  const el = $(id);
  if (el) el.innerHTML = `<td style="width:1%;white-space:nowrap;padding-right:0.75em;">${label}</td><td style="color:var(--card-foreground);font-weight:500;word-break:break-all">${value}</td>`;
}

function renderOrderId() {
  if (ORDER.tradeId && !ORDER.tradeId.startsWith('{{')) {
    renderRow('display-order-id', t('order_id'), ORDER.tradeId);
  }
}

function initOrder() {
  if (!ORDER.tradeId || ORDER.tradeId.startsWith('{{')) {
    showNotFound();
    return;
  }

  $('display-amount').textContent = `${ORDER.actualAmount} ${ORDER.token}`;
  $('field-amount').textContent   = `${ORDER.actualAmount} ${ORDER.token}`;
  $('field-address').innerHTML    = formatAddr(ORDER.receiveAddress);
  if (ORDER.amount && !ORDER.amount.startsWith('{{')) {
    renderRow('display-fiat', t('order_amount'), `${ORDER.amount} ${ORDER.currency || ''}`);
  }
  $('display-network').innerHTML = networkTag(ORDER.network);
  renderOrderId();

  new QRCode($('qrcode'), {
    text:         ORDER.receiveAddress,
    width:        176,
    height:       176,
    colorDark:    '#111111',
    colorLight:   '#ffffff',
    correctLevel: QRCode.CorrectLevel.M,
  });

  initCountdown();
  checkOrderStatus();
  initClipboard();
}

// Clipboard
let _clipboard = null;

function initClipboard() {
  const targets = [
    { id: 'copy-amt-box',    text: ORDER.actualAmount, iconId: 'btn-copy-amt2'   },
    { id: 'copy-addr-box',   text: ORDER.receiveAddress, iconId: 'btn-copy-addr' },
    { id: 'btn-copy-amount', text: ORDER.actualAmount, iconId: 'btn-copy-amount' },
  ];

  if (typeof ClipboardJS !== 'undefined') {
    targets.forEach(({ id, text }) => {
      const el = $(id);
      if (el) el.setAttribute('data-clipboard-text', text);
    });
    if (_clipboard) _clipboard.destroy();
    _clipboard = new ClipboardJS(targets.map(c => '#' + c.id).join(', '));
    _clipboard.on('success', e => {
      e.clearSelection();
      const tgt = targets.find(c => c.id === e.trigger.id);
      showToast();
      if (tgt) flashCheck(tgt.iconId);
    });
  }
}
// Countdown color interpolation (green → orange → red by ratio)
function lerpColor(a, b, t) {
  return `rgb(${Math.round(a[0]+(b[0]-a[0])*t)},${Math.round(a[1]+(b[1]-a[1])*t)},${Math.round(a[2]+(b[2]-a[2])*t)})`;
}
function ratioColor(ratio) {
  if (ratio >= 0.5) return '#22c55e';
  if (ratio >= 0.2) return lerpColor([34,197,94], [249,115,22], (0.5 - ratio) / 0.3);
  return lerpColor([249,115,22], [239,68,68], (0.2 - ratio) / 0.2);
}

const CIRCUMFERENCE = 2 * Math.PI * 20; // r=20 in iOS SVG ring
let countdownInterval = null;
let totalSeconds = 0;
let _expiresAt = null;

const pad = (n) => String(n).padStart(2, '0');

function initCountdown() {
  const parseTime = (raw) => new Date(/^\d+$/.test(String(raw)) ? Number(raw) : raw);
  _expiresAt = parseTime(ORDER.expirationTime);
  const remaining = Math.max(0, Math.round((_expiresAt - Date.now()) / 1000));
  if (remaining <= 0) { showExpired(); return; }
  // 用 createdAt 和 expirationTime 推算总有效时长，否则默认 600 秒
  const created = parseTime(ORDER.createdAt);
  totalSeconds = (!isNaN(created.getTime()) && !isNaN(_expiresAt.getTime()))
    ? Math.max(1, Math.round((_expiresAt - created) / 1000))
    : 600;
  $('ring').style.strokeDasharray = CIRCUMFERENCE;
  // 让初始渲染就从 remaining-1 开始，避免视觉上停顿一秒
  _expiresAt = new Date(_expiresAt.getTime() - 1000);
  tickCountdown();
  countdownInterval = setInterval(tickCountdown, 1000);
}

function tickCountdown() {
  const remaining = Math.max(0, Math.round((_expiresAt - Date.now()) / 1000));
  const h = Math.floor(remaining / 3600);
  const m = Math.floor((remaining % 3600) / 60);
  const s = remaining % 60;

  const timeStr = h > 0
    ? `${pad(h)}:${pad(m)}:${pad(s)}`
    : `${pad(m)}:${pad(s)}`;

  $('countdown').textContent = timeStr;
  const inlineEl = $('countdown-inline');

  const color = ratioColor(totalSeconds > 0 ? remaining / totalSeconds : 0);
  if (inlineEl) { inlineEl.textContent = timeStr; inlineEl.style.color = color; }

  $('ring').style.strokeDashoffset = CIRCUMFERENCE * (1 - (totalSeconds > 0 ? remaining / totalSeconds : 0));

  $('countdown').style.color = color;
  $('ring').style.stroke = color;

  if (remaining <= 0) {
    clearInterval(countdownInterval);
    showExpired();
  }
}

// Polling
let _pollTimer = null;
let _pollStopped = false;
let _pollErrors = 0;
const MAX_POLL_ERRORS = 5;

function checkOrderStatus() {
  if (_pollStopped || !ORDER.tradeId || ORDER.tradeId.startsWith('{{')) return;

  fetch(`/pay/check-status/${encodeURIComponent(ORDER.tradeId)}`, {
    headers: { Accept: 'application/json' },
    signal: AbortSignal.timeout(8000),
  })
    .then(res => {
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      return res.json();
    })
    .then(({ data }) => {
      _pollErrors = 0;
      if (data?.status === 2) onPaymentSuccess();
      else if (data?.status === 3) showExpired();
      else scheduleNextPoll();
    })
    .catch(() => {
      if (++_pollErrors >= MAX_POLL_ERRORS) showTimeout();
      else scheduleNextPoll();
    });
}

const scheduleNextPoll = () => {
  if (!_pollStopped) _pollTimer = setTimeout(checkOrderStatus, 5000);
};

function stopPolling() {
  _pollStopped = true;
  clearTimeout(_pollTimer);
}

// State panels
const SCREENS = ['screen-success', 'screen-expired', 'screen-timeout', 'screen-not-found'];

function syncStatusCardHeight() {
  const card = $('status-card');
  const inner = card?.querySelector('.status-card-inner');
  if (!card || !inner) return;

  const activeFace = card.classList.contains('is-flipped')
    ? card.querySelector('.status-face-back')
    : card.querySelector('.status-face-front');

  const nextHeight = activeFace?.scrollHeight || 0;
  if (nextHeight > 0) {
    card.style.height = `${nextHeight}px`;
    inner.style.height = `${nextHeight}px`;
  }
}

function setActiveScreen(state) {
  SCREENS.forEach(id => $(id).classList.toggle('is-active', id === `screen-${state}`));
}

function flipToState(state) {
  setActiveScreen(state);
  $('status-card').classList.add('is-flipped');
  requestAnimationFrame(syncStatusCardHeight);
}

function showStatePanel(state) {
  $('order-info').style.display = state === 'not-found' ? 'none' : '';
  flipToState(state);
}

function goBack() {
  if (document.referrer) {
    window.location.href = document.referrer;
  } else if (ORDER.redirectUrl && !ORDER.redirectUrl.startsWith('{{')) {
    window.location.href = ORDER.redirectUrl;
  } else {
    history.back();
  }
}

function hideStatePanel() {
  $('status-card').classList.remove('is-flipped');
  $('order-info').style.display = '';
  requestAnimationFrame(syncStatusCardHeight);
}

function onPaymentSuccess() {
  stopPolling();
  clearInterval(countdownInterval);
  showStatePanel('success');
  if (ORDER.redirectUrl && !ORDER.redirectUrl.startsWith('{{')) {
    setTimeout(() => { window.location.href = ORDER.redirectUrl; }, 3000);
  }
}

function showExpired() {
  stopPolling();
  clearInterval(countdownInterval);
  $('btn-transferred').disabled = true;
  showStatePanel('expired');
}

function showTimeout() {
  stopPolling();
  showStatePanel('timeout');
}

function showNotFound() {
  stopPolling();
  clearInterval(countdownInterval);
  $('btn-transferred').disabled = true;
  showStatePanel('not-found');
}

function retryPolling() {
  _pollErrors = 0;
  _pollStopped = false;
  hideStatePanel();
  checkOrderStatus();
}

document.addEventListener('DOMContentLoaded', () => {
  setLang(detectLang());
  initOrder();
  syncStatusCardHeight();
  window.addEventListener('resize', syncStatusCardHeight);
});

