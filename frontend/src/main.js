// UpscalerHub Frontend – main.js
// Wails bindings are available as window.go.main.App.*

const App = window.go?.main?.App;

// ─── Localization ────────────────────────────────────────────────────────
const translations = {
  en: {
    // Nav & Common
    scan: "⟳ Scan Games",
    addManual: "＋ Add Manually",
    search: "Search games...",
    manage: "Manage",
    openFolder: "Open Folder",
    play: "Play",
    restore: "Restore",
    detectingGpu: "Detecting GPU…",
    gpuError: "⚠️ GPU error",
    noGpu: "No GPU",
    optiStatus: "OptiScaler State",
    nativeSwap: "Native Swap",
    installedSuccess: "Installed successfully!",
    restoreSuccess: "Original native DLL restored!",
    favorite: "Favorite",
    favorites: "FAVORITES",
    otherGames: "OTHER GAMES",
    navGames: "Games",
    navSettings: "Settings",
    navLibrary: "Upscaler Library",
    navAbout: "About",
    noGamesFound: "No games found",
    scanHint: "Click <strong>Scan</strong> to detect installed games",

    // Manage View
    optiStatusTitle: "OptiScaler Status",
    installed: "Installed",
    notInstalled: "Not Installed",
    currentVer: "Current",
    versionLabel: "Version:",
    installUpdateTitle: "Install / Update",
    versionLabelBare: "Version",
    injectionLabel: "Injection",
    dxgiOpt: "dxgi.dll (DirectX 12)",
    versionOpt: "version.dll (Universal)",
    winmmOpt: "winmm.dll (Alternative)",
    wininetOpt: "wininet.dll (Alternative)",
    winhttpOpt: "winhttp.dll (Alternative)",
    installFakenvapi: "Install Fakenvapi (AMD/Intel)",
    installNukem: "Install NukemFG (Frame Gen)",
    installBtn: "✦ Install",
    uninstallBtn: "✕ Uninstall",
    detectedComponentsTitle: "Detected Components",
    nativeSwapperTitle: "Official Upscalers (Native Swap)",
    nativeSwapperDesc: "Replace native original DLLs with other official versions downloaded from your Library.",
    noNativeUpscalers: "No native upscalers detected in this game.",
    dlssPresetsTitle: "DLSS Presets (NVAPI)",
    dlssPresetsDesc: "Modify driver-level DLSS Render Presets for this game.",

    // Library View
    refreshManifest: "Refresh Manifest",
    tabDLSS: "DLSS",
    tabDLSSFG: "DLSS Frame Gen",
    tabDLSSRR: "DLSS Ray Reconstruction",
    tabFSR31DX12: "FSR 3.1 (DX12)",
    tabFSR31VK: "FSR 3.1 (Vulkan)",
    tabXeSS: "XeSS",
    thVersion: "Version",
    thSize: "Size",
    thStatus: "Status",
    thAction: "Action",

    // Settings
    cacheMgmtTitle: "Cache Management",
    cacheMgmtDesc: "Manage downloaded OptiScaler versions",
    componentsTitle: "Components",
    componentsDesc: "Installed versions of managed components",

    // About
    updateAvailable: "Update available",
    notProvided: "Not provided",
    lastBuildDate: "Last Build Date",
    checkUpdatesBtn: "Check for updates",
    githubBtn: "GitHub Repository",

    // Toasts & Dynamic System messages
    dlssPresetUpdated: "DLSS Preset updated successfully!",
    failedUpdatePreset: "Failed to update preset",
    scanFound: "Scan complete. {0} games found.",
    scanFailed: "Scan failed",
    addedGame: "Added {0}",
    failedAddGame: "Failed to add game",
    removedGame: "Removed {0}",
    confirmUninstall: "Uninstall OptiScaler from this game?",
    uninstallSuccess: "Uninstalled successfully.",
    uninstallFailed: "Uninstall failed",
    confirmDelete: "Delete OptiScaler v{0}?",
    versionDeleted: "Version deleted",
    deleteFailed: "Failed to delete",
    updateAvailToast: "Update available: v{0}",
    latestVersionToast: "You are running the latest version",
    updateCheckFailed: "Update check failed",
    manifestRefreshed: "Manifest refreshed.",
    failedManifest: "Failed to fetch Upscaler Library manifest.",
    downloadedSucc: "Successfully downloaded v{0}",
    failedDownload: "Failed to download",
    noVersionsCache: "No cached versions",
    deletePlain: "Delete",
    libDownloaded: "✓ Downloaded",
    libCloud: "☁ Cloud",
    libReady: "Ready",
    libDownloadBtn: "Download",
    btnInstallVer: "Install",
    btnDownInstallVer: "⬇ Download & Install",
    noVersManifest: "No versions available in manifest",
    installingPrefix: "Installing…",
    checkingUpdatesPre: "Checking…",
    loadingManifestPre: "Loading manifest…",
    confirmRemoveList: "Remove \"{0}\" from the list?"
  },
  pt: {
    // Nav & Common
    scan: "⟳ Escanear Jogos",
    addManual: "＋ Adicionar Manualmente",
    search: "Pesquisar jogos...",
    manage: "Gerenciar",
    openFolder: "Abrir Pasta",
    play: "Iniciar",
    restore: "Restaurar",
    detectingGpu: "Detectando GPU…",
    gpuError: "⚠️ Erro na GPU",
    noGpu: "Nenhuma GPU",
    optiStatus: "Status do OptiScaler",
    nativeSwap: "Troca de DLL Nativa",
    installedSuccess: "Instalado com sucesso!",
    restoreSuccess: "DLL nativa original restaurada!",
    favorite: "Favorito",
    favorites: "FAVORITOS",
    otherGames: "OUTROS JOGOS",
    navGames: "Jogos",
    navSettings: "Configurações",
    navLibrary: "Biblioteca",
    navAbout: "Sobre",
    noGamesFound: "Nenhum jogo encontrado",
    scanHint: "Clique em <strong>Escanear</strong> para detectar os jogos instalados",

    // Manage View
    optiStatusTitle: "Status do OptiScaler",
    installed: "Instalado",
    notInstalled: "Não Instalado",
    currentVer: "Atual",
    versionLabel: "Versão:",
    installUpdateTitle: "Instalação / Atualização",
    versionLabelBare: "Versão",
    injectionLabel: "Injeção",
    dxgiOpt: "dxgi.dll (DirectX 12)",
    versionOpt: "version.dll (Universal)",
    winmmOpt: "winmm.dll (Alternativo)",
    wininetOpt: "wininet.dll (Alternativo)",
    winhttpOpt: "winhttp.dll (Alternativo)",
    installFakenvapi: "Instalar Fakenvapi (AMD/Intel)",
    installNukem: "Instalar NukemFG (Frame Gen)",
    installBtn: "✦ Instalar",
    uninstallBtn: "✕ Desinstalar",
    detectedComponentsTitle: "Componentes Detectados",
    nativeSwapperTitle: "Upscalers Oficiais (Troca Nativa)",
    nativeSwapperDesc: "Substitua DLLs nativas originais por versões oficiais baixadas da sua Biblioteca.",
    noNativeUpscalers: "Nenhum upscaler nativo detectado neste jogo.",
    dlssPresetsTitle: "Presets DLSS (NVAPI)",
    dlssPresetsDesc: "Modifique os presets de renderização DLSS no nível do driver para este jogo.",

    // Library View
    refreshManifest: "Atualizar Manifesto",
    tabDLSS: "DLSS",
    tabDLSSFG: "DLSS Frame Gen",
    tabDLSSRR: "DLSS Ray Reconstruction",
    tabFSR31DX12: "FSR 3.1 (DX12)",
    tabFSR31VK: "FSR 3.1 (Vulkan)",
    tabXeSS: "XeSS",
    thVersion: "Versão",
    thSize: "Tamanho",
    thStatus: "Status",
    thAction: "Ação",

    // Settings
    cacheMgmtTitle: "Gerenciamento de Cache",
    cacheMgmtDesc: "Gerenciar versões do OptiScaler baixadas localmente",
    componentsTitle: "Componentes",
    componentsDesc: "Versões instaladas dos componentes gerenciados",

    // About
    updateAvailable: "Atualização disponível",
    notProvided: "Não fornecido",
    lastBuildDate: "Data da Última Build",
    checkUpdatesBtn: "Verificar atualizações",
    githubBtn: "Repositório do GitHub",

    // Toasts & Dynamic System messages
    dlssPresetUpdated: "Preset DLSS atualizado com sucesso!",
    failedUpdatePreset: "Falha ao atualizar preset",
    scanFound: "Escaneamento concluído. {0} jogos encontrados.",
    scanFailed: "Falha no escaneamento",
    addedGame: "{0} adicionado",
    failedAddGame: "Falha ao adicionar o jogo",
    removedGame: "{0} removido",
    confirmUninstall: "Desinstalar o OptiScaler deste jogo?",
    uninstallSuccess: "Desinstalado com sucesso.",
    uninstallFailed: "Falha na desinstalação",
    confirmDelete: "Excluir o OptiScaler v{0}?",
    versionDeleted: "Versão excluída",
    deleteFailed: "Falha ao excluir",
    updateAvailToast: "Atualização disponível: v{0}",
    latestVersionToast: "Você está usando a versão mais recente",
    updateCheckFailed: "Falha ao verificar as atualizações",
    manifestRefreshed: "Manifesto atualizado.",
    failedManifest: "Falha ao buscar o manifesto da Biblioteca.",
    downloadedSucc: "Versão v{0} baixada com sucesso",
    failedDownload: "Falha ao baixar",
    noVersionsCache: "Nenhuma versão em cache",
    deletePlain: "Excluir",
    libDownloaded: "✓ Baixado",
    libCloud: "☁ Na Nuvem",
    libReady: "Pronto",
    libDownloadBtn: "Baixar",
    btnInstallVer: "Instalar",
    btnDownInstallVer: "⬇ Baixar & Instalar",
    noVersManifest: "Nenhuma versão disponível no manifesto",
    installingPrefix: "Instalando…",
    checkingUpdatesPre: "Verificando…",
    loadingManifestPre: "Carregando manifesto…",
    confirmRemoveList: "Remover \"{0}\" da lista?"
  }
};

let currentLang = "en";

function updateUI() {
  const t = translations[currentLang];
  
  // Data-bound elements
  document.querySelectorAll("[data-i18n]").forEach(el => {
    const key = el.getAttribute("data-i18n");
    if (t[key]) el.textContent = t[key];
  });
  document.querySelectorAll("[data-i18n-title]").forEach(el => {
    const key = el.getAttribute("data-i18n-title");
    if (t[key]) el.title = t[key];
  });
  document.querySelectorAll("[data-i18n-html]").forEach(el => {
    const key = el.getAttribute("data-i18n-html");
    if (t[key]) el.innerHTML = t[key];
  });

  document.getElementById("btn-scan").innerHTML = `
    <span style="display: flex; align-items: center; gap: 6px;">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>
      ${t.scan.replace("⟳ ", "")}
    </span>
  `;
  document.getElementById("btn-add").innerHTML = `
    <span style="display: flex; align-items: center; gap: 6px;">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
      ${t.addManual.replace("＋ ", "")}
    </span>
  `;
  document.getElementById("search-input").placeholder = t.search;
  
  const gpuText = document.getElementById("gpu-text");
  if (gpuText.textContent === translations.en.detectingGpu || gpuText.textContent === translations.pt.detectingGpu) {
    gpuText.textContent = t.detectingGpu;
  }
  
  // Update views if active
  renderGames(document.getElementById("search-input").value);
  if (currentView === "manage" && managingIndex >= 0) {
    openManage(managingIndex);
  }
  if (currentView === "settings") loadSettings();
  if (currentView === "library") loadLibrary();
  if (currentView === "about") loadAbout();
  
  document.getElementById("btn-lang-toggle").textContent = currentLang === "en" ? "ENG" : "PT-BR";
}

document.getElementById("btn-lang-toggle").addEventListener("click", () => {
  currentLang = currentLang === "en" ? "pt" : "en";
  updateUI();
});

// ─── View Routing ───────────────────────────────────────────────────────
const views = document.querySelectorAll(".view");
const navBtns = document.querySelectorAll(".nav-btn");
let currentView = "games";

navBtns.forEach((btn) => {
  btn.addEventListener("click", () => switchView(btn.dataset.view));
});

function switchView(name) {
  currentView = name;
  if (name === "manage") {
    // Manage is special — not in nav
    views.forEach((v) => v.classList.remove("active"));
    document.getElementById("view-manage").classList.add("active");
    return;
  }
  views.forEach((v) => v.classList.remove("active"));
  document.getElementById(`view-${name}`).classList.add("active");
  navBtns.forEach((b) => b.classList.toggle("active", b.dataset.view === name));

  if (name === "settings") loadSettings();
  if (name === "about") loadAbout();
  if (name === "library") loadLibrary();
}

// ─── Title Bar Controls ─────────────────────────────────────────────────
document.getElementById("btn-minimize").addEventListener("click", () => {
  window.runtime?.WindowMinimise();
});
document.getElementById("btn-maximize").addEventListener("click", () => {
  window.runtime?.WindowToggleMaximise();
});
document.getElementById("btn-close").addEventListener("click", () => {
  window.runtime?.Quit();
});

// ─── Game List ──────────────────────────────────────────────────────────
let games = [];
let managingIndex = -1;

async function loadGames() {
  if (!App) return;
  try {
    games = (await App.GetGames()) || [];
    renderGames();
  } catch (e) {
    console.error("Failed to load games:", e);
  }
}

function renderGames(filter = "") {
  const grid = document.getElementById("game-list");
  const countEl = document.getElementById("game-count");
  let emptyEl = document.getElementById("empty-state");

  // Keep empty state in a globally safe variable directly attached to grid if it gets wiped
  if (!emptyEl && grid._savedEmptyState) {
    emptyEl = grid._savedEmptyState;
  } else if (emptyEl) {
    grid._savedEmptyState = emptyEl;
  }

  const filtered = filter
    ? games.filter((g) => g.name.toLowerCase().includes(filter.toLowerCase()))
    : games;

  countEl.textContent = filtered.length;

  if (filtered.length === 0) {
    grid.innerHTML = "";
    if (emptyEl) {
      grid.appendChild(emptyEl);
      emptyEl.style.display = "block";
    }
    return;
  }

  if (emptyEl) emptyEl.style.display = "none";
  
  const favorites = [];
  const supported = [];
  const unsupported = [];
  filtered.forEach((g) => {
    if (g.isFavorite) {
      favorites.push(g);
    } else if (g.dlssVersion || g.fsrVersion || g.xessVersion) {
      supported.push(g);
    } else {
      unsupported.push(g);
    }
  });

  const renderGameCard = (g) => {
    const t = translations[currentLang];
    const realIndex = games.indexOf(g);
    const coverUrl = g.coverImageUrl || '';
    let cover;
    if (coverUrl) {
      const fallbackUrl = coverUrl.replace('library_600x900_2x.jpg', 'header.jpg');
      cover = `<img class="game-cover" src="${coverUrl}" alt="" loading="lazy" onerror="if(this.src.includes('library_600x900')){this.src='${fallbackUrl}'}else{this.outerHTML='<div class=\\'game-cover-placeholder\\'>🎮</div>'}">`;
    } else {
      cover = `<div class="game-cover-placeholder">🎮</div>`;
    }

    const opti = g.isOptiScalerInstalled
      ? `<span class="game-opti-badge">OptiScaler ${g.optiScalerVersion || ""}</span>`
      : "";

    let techBadges = '';
    if (g.dlssVersion) techBadges += '<span class="tech-badge tech-dlss">DLSS</span>';
    if (g.fsrVersion) techBadges += '<span class="tech-badge tech-fsr">FSR</span>';
    if (g.xessVersion) techBadges += '<span class="tech-badge tech-xess">XeSS</span>';

    const removeBtn =
      g.platform === "Manual"
        ? `<button class="game-btn-remove" data-path="${encodeURIComponent(g.installPath)}" data-name="${g.name}" title="Remove">✕</button>`
        : "";

    return `
      <div class="game-card" data-index="${realIndex}">
        ${cover}
        <div class="game-info">
          <div class="game-name" title="${g.name}">${g.name}</div>
          <div class="game-meta">
            <span class="game-platform">${g.platform}</span>
            ${opti}
          </div>
          ${techBadges ? `<div class="game-techs">${techBadges}</div>` : ''}
        </div>
        <div class="game-actions">
          <button class="game-btn-manage" data-manage="${realIndex}">${t.manage}</button>
          ${removeBtn}
        </div>
      </div>
    `;
  };

  let html = "";
  const t = translations[currentLang];
  if (favorites.length > 0) {
    html += `<h2 class="category-title">★ ${t.favorites} (${favorites.length})</h2>`;
    html += favorites.map(renderGameCard).join("");
  }
  if (supported.length > 0) {
    html += `<h2 class="category-title">${t.otherGames} (${supported.length})</h2>`;
    html += supported.map(renderGameCard).join("");
  }
  if (unsupported.length > 0) {
    html += `<h2 class="category-title">${t.otherGames} (${unsupported.length})</h2>`;
    html += unsupported.map(renderGameCard).join("");
  }
  
  grid.innerHTML = html;

  // Event delegation for manage/remove buttons
  grid.querySelectorAll(".game-btn-manage").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      e.stopPropagation();
      openManage(parseInt(btn.dataset.manage));
    });
  });
  grid.querySelectorAll(".game-btn-remove").forEach((btn) => {
    btn.addEventListener("click", (e) => {
      e.stopPropagation();
      removeGame(decodeURIComponent(btn.dataset.path), btn.dataset.name);
    });
  });
  // Clicking the card itself also opens manage
  grid.querySelectorAll(".game-card").forEach((card) => {
    card.addEventListener("click", () =>
      openManage(parseInt(card.dataset.index)),
    );
  });
}

// Search
document.getElementById("search-input").addEventListener("input", (e) => {
  renderGames(e.target.value);
});

// Scan
document.getElementById("btn-scan").addEventListener("click", async () => {
  if (!App) return;
  const t = translations[currentLang];
  const btn = document.getElementById("btn-scan");
  btn.disabled = true;
  btn.innerHTML = `<span class="spinner"></span> ${t.checkingUpdatesPre}`;
  try {
    games = (await App.ScanGames()) || [];
    renderGames(document.getElementById("search-input").value);
    showToast(t.scanFound.replace("{0}", games.length), "success");
  } catch (e) {
    showToast(t.scanFailed + ": " + e, "error");
  } finally {
    btn.disabled = false;
    updateUI();
  }
});

// Add manually
document.getElementById("btn-add").addEventListener("click", async () => {
  if (!App) return;
  const t = translations[currentLang];
  try {
    const game = await App.AddGameManually();
    if (game) {
      games = (await App.GetGames()) || [];
      renderGames();
      showToast(t.addedGame.replace("{0}", game.name), "success");
    }
  } catch (e) {
    showToast(t.failedAddGame, "error");
  }
});

async function removeGame(installPath, name) {
  if (!App) return;
  const t = translations[currentLang];
  if (!confirm(t.confirmRemoveList.replace("{0}", name))) return;
  await App.RemoveGame(installPath);
  games = (await App.GetGames()) || [];
  renderGames();
  showToast(t.removedGame.replace("{0}", name), "success");
}

// ─── Manage Game View ───────────────────────────────────────────────────
async function openManage(index) {
  if (!App) return;
  managingIndex = index;
  const game = games[index];
  if (!game) return;

  document.getElementById("manage-title").textContent = game.name;
  const t = translations[currentLang];

  // Status
  const statusEl = document.getElementById("opti-status");
  const versionRow = document.getElementById("opti-version-row");
  const statusLabelEl = document.getElementById("opti-status-label");
  if (statusLabelEl) statusLabelEl.textContent = t.optiStatus + ":";

  if (game.isOptiScalerInstalled) {
    statusEl.innerHTML =
      `<span class="status-dot on"></span> <span>${t.installed}</span>`;
    versionRow.style.display = "flex";
    document.getElementById("opti-installed-ver").textContent =
      game.optiScalerVersion || "Unknown";
    document.getElementById("btn-uninstall").style.display = "inline-flex";
  } else {
    statusEl.innerHTML =
      `<span class="status-dot off"></span> <span>${t.notInstalled}</span>`;
    versionRow.style.display = "none";
    document.getElementById("btn-uninstall").style.display = "none";
  }

  // Load versions
  try {
    const status = await App.GetComponentStatus();
    const sel = document.getElementById("sel-version");
    sel.innerHTML = "";
    const versions = status.optiScalerVersions || [];
    if (versions.length === 0) {
      sel.innerHTML = "<option>No versions available</option>";
    } else {
      versions.forEach((v) => {
        const opt = document.createElement("option");
        opt.value = v;
        opt.textContent = v;
        sel.appendChild(opt);
      });
    }
  } catch (e) {
    console.error("Failed to get versions:", e);
  }

  // Detected components
  const comps = document.getElementById("detected-components");
  comps.innerHTML = "";
  const techs = [
    ["DLSS", game.dlssVersion],
    ["DLSS FG", game.dlssFrameGenVersion],
    ["FSR", game.fsrVersion],
    ["XeSS", game.xessVersion],
  ];
  techs.forEach(([name, ver]) => {
    if (ver) {
      comps.innerHTML += `<div class="component-item"><span class="label">${name}</span><span class="value">${ver}</span></div>`;
    }
  });
  if (!comps.innerHTML) {
    comps.innerHTML =
      '<div class="component-item"><span class="label">No upscaling DLLs detected</span></div>';
  }

  // Cover and Actions Update
  const coverEl = document.getElementById("manage-cover-image");
  if (game.coverImageUrl) {
    coverEl.src = game.coverImageUrl;
    coverEl.style.display = 'block';
  } else {
    coverEl.style.display = 'none';
  }

  const favIcon = document.getElementById("icon-favorite");
  if (game.isFavorite) {
    favIcon.setAttribute("fill", "currentColor");
  } else {
    favIcon.setAttribute("fill", "none");
  }

  // Native Swapper handling
  const swapperTechTypes = {
    "DLSS": "dlss",
    "DLSS FG": "dlssg",
    "FSR": "fsr31dx12",
    "XeSS": "xess"
  };

  const swapperContainer = document.getElementById("native-swapper-technologies");
  const swapperEmpty = document.getElementById("native-swapper-empty");
  swapperContainer.innerHTML = "";
  
  let hasTech = false;
  techs.forEach(([name, ver]) => {
    if (ver) {
      hasTech = true;
      const t = swapperTechTypes[name];
      const techHtml = `
      <div class="swapper-tech-row" data-tech="${t}" data-current-ver="${ver}">
         <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:8px;">
            <label style="font-weight: 600;">${name}</label>
            <div style="display:flex; gap:6px; align-items:center;">
               <span style="font-size:11px; color:var(--text-muted);">${translations[currentLang].currentVer}: v${ver}</span>
               <button class="btn btn-secondary btn-restore-native-row" data-tech="${t}" title="${translations[currentLang].restore}" style="padding:4px 8px; font-size:11px;">↺ ${translations[currentLang].restore}</button>
            </div>
         </div>
         <div class="native-version-list" data-tech="${t}" style="max-height:200px; overflow-y:auto; border:1px solid var(--border-subtle); border-radius:var(--radius-sm); background:var(--bg-input);">
            <div style="padding:12px; text-align:center; color:var(--text-muted); font-size:11px;">${translations[currentLang].loadingManifestPre}</div>
         </div>
      </div>`;
      swapperContainer.innerHTML += techHtml;
    }
  });

  if (!hasTech) {
    swapperContainer.style.display = "none";
    swapperEmpty.style.display = "block";
  } else {
    swapperContainer.style.display = "flex";
    swapperEmpty.style.display = "none";
    if (swapperManifest) {
      populateNativeVersionsForEachRow();
    } else {
      loadLibrary().then(() => populateNativeVersionsForEachRow());
    }
  }

  // Load and populate NVAPI Presets
  const presetsCard = document.getElementById("dlss-presets-card");
  const rowDlss = document.getElementById("row-dlss-preset");
  const rowDlssd = document.getElementById("row-dlssd-preset");
  const selDlss = document.getElementById("sel-preset-dlss");
  const selDlssd = document.getElementById("sel-preset-dlssd");
  
  let hasDlss = !!game.dlssVersion;
  
  if (hasDlss) {
    App.GetNVAPIPresets(game.executablePath).then(res => {
      if (res && res.foundProfile) {
        presetsCard.style.display = "block";
        if (hasDlss) {
          rowDlss.style.display = "flex";
          selDlss.value = res.dlssPreset;
          rowDlssd.style.display = "flex";
          selDlssd.value = res.dlssdPreset;
        }
      } else {
        presetsCard.style.display = "none";
      }
    }).catch(err => {
      console.error(err);
      presetsCard.style.display = "none";
    });
  } else {
    presetsCard.style.display = "none";
  }

  switchView("manage");
}

// Back button
document.getElementById("btn-back").addEventListener("click", () => {
  switchView("games");
});

// Open folder - open the game's actual game directory (where DLLs are)
document.getElementById("btn-open-folder").addEventListener("click", () => {
  if (!App || managingIndex < 0) return;
  const game = games[managingIndex];
  if (!game) return;
  // Use the path of a detected DLL to find the actual game directory
  const dllPath = game.dlssPath || game.fsrPath || game.xessPath || game.dlssFrameGenPath;
  if (dllPath) {
    const lastSlash = Math.max(dllPath.lastIndexOf('\\'), dllPath.lastIndexOf('/'));
    const dir = lastSlash > 0 ? dllPath.substring(0, lastSlash) : game.installPath;
    App.OpenFolder(dir);
  } else {
    App.OpenFolder(game.installPath);
  }
});

// Launch Game
document.getElementById("btn-launch").addEventListener("click", () => {
  if (!App || managingIndex < 0) return;
  App.LaunchGame(managingIndex).catch(console.error);
});

// Favorite Game
document.getElementById("btn-favorite").addEventListener("click", async () => {
  if (!App || managingIndex < 0) return;
  const isFav = await App.ToggleFavorite(managingIndex);
  games = await App.GetGames() || [];
  
  const favIcon = document.getElementById("icon-favorite");
  if (isFav) {
    favIcon.setAttribute("fill", "currentColor");
  } else {
    favIcon.setAttribute("fill", "none");
  }
});

// Install
document.getElementById("btn-install").addEventListener("click", async () => {
  if (!App || managingIndex < 0) return;
  const t = translations[currentLang];
  const version = document.getElementById("sel-version").value;
  const injection = document.getElementById("sel-injection").value;
  const fakenvapi = document.getElementById("chk-fakenvapi").checked;
  const nukem = document.getElementById("chk-nukem").checked;

  const btn = document.getElementById("btn-install");
  const progress = document.getElementById("install-progress");
  btn.disabled = true;
  btn.innerHTML = `<span class="spinner"></span> ${t.installingPrefix}`;
  progress.style.display = "block";

  try {
    const result = await App.InstallOptiScaler(
      managingIndex,
      version,
      injection,
      fakenvapi,
      nukem,
    );
    if (result.success) {
      showToast(t.installedSuccess, "success");
      games = (await App.GetGames()) || [];
      openManage(managingIndex);
    } else {
      showToast(result.message, "error");
    }
  } catch (e) {
    showToast("Install failed: " + e, "error");
  } finally {
    btn.disabled = false;
    btn.textContent = t.installBtn;
    progress.style.display = "none";
  }
});

// Uninstall
document.getElementById("btn-uninstall").addEventListener("click", async () => {
  if (!App || managingIndex < 0) return;
  const t = translations[currentLang];
  if (!confirm(t.confirmUninstall)) return;

  try {
    const result = await App.UninstallOptiScaler(managingIndex);
    if (result.success) {
      showToast(t.uninstallSuccess, "success");
      games = (await App.GetGames()) || [];
      openManage(managingIndex);
    } else {
      showToast(result.message, "error");
    }
  } catch (e) {
    showToast(t.uninstallFailed + ": " + e, "error");
  }
});

// Download progress events
window.runtime?.EventsOn("download-progress", (pct) => {
  const fill = document.querySelector(".progress-fill");
  if (fill) fill.style.width = pct + "%";
});

// Games updated event (cover images loaded)
window.runtime?.EventsOn("games-updated", (updatedGames) => {
  games = updatedGames || [];
  renderGames(document.getElementById("search-input").value);
});

// ─── Settings View ──────────────────────────────────────────────────────
async function loadSettings() {
  if (!App) return;
  const t = translations[currentLang];

  // Cache list
  try {
    const versions = (await App.GetDownloadedVersions()) || [];
    const cacheList = document.getElementById("cache-list");
    if (versions.length === 0) {
      cacheList.innerHTML =
        `<div style="color:var(--text-muted);font-size:12px">${t.noVersionsCache}</div>`;
    } else {
      cacheList.innerHTML = versions
        .map(
          (v) => `
        <div class="cache-item">
          <span class="cache-item-label">OptiScaler v${v}</span>
          <button class="cache-delete-btn" data-version="${v}">🗑 ${t.deletePlain}</button>
        </div>
      `,
        )
        .join("");

      cacheList.querySelectorAll(".cache-delete-btn").forEach((btn) => {
        btn.addEventListener("click", async () => {
          if (!confirm(t.confirmDelete.replace("{0}", btn.dataset.version))) return;
          try {
            await App.DeleteCachedVersion(btn.dataset.version);
            showToast(t.versionDeleted, "success");
            loadSettings();
          } catch (e) {
            showToast(t.deleteFailed + ": " + e, "error");
          }
        });
      });
    }
  } catch (e) {
    console.error("Failed to load cache:", e);
  }

  // Component versions
  try {
    const status = await App.GetComponentStatus();
    const compsEl = document.getElementById("component-versions");
    compsEl.innerHTML = [
      ["OptiScaler", status.local.OptiScalerVersion || t.notInstalled],
      ["Fakenvapi", status.local.FakenvapiVersion || t.notInstalled],
      ["NukemFG", status.local.NukemFGVersion || t.notProvided],
    ]
      .map(
        ([name, ver]) => `
      <div class="component-info-row">
        <span class="label">${name}</span>
        <span>${ver}</span>
      </div>
    `,
      )
      .join("");
  } catch (e) {
    console.error("Failed to load components:", e);
  }
}

// ─── About View ─────────────────────────────────────────────────────────
async function loadAbout() {
  if (!App) return;
  const t = translations[currentLang];
  try {
    const status = await App.GetComponentStatus();
    const local = status.local || {};
    const remote = status.remote || {};

    // OptiScaler
    document.getElementById('about-opti-ver').textContent = local.OptiScalerVersion || t.notInstalled;
    const optiUpdate = document.getElementById('about-opti-update');
    if (remote.OptiScalerVersion && remote.OptiScalerVersion !== local.OptiScalerVersion) {
      optiUpdate.classList.remove('hidden');
    } else {
      optiUpdate.classList.add('hidden');
    }

    // Fakenvapi
    document.getElementById('about-fv-ver').textContent = local.FakenvapiVersion || t.notInstalled;
    const fvUpdate = document.getElementById('about-fv-update');
    if (remote.FakenvapiVersion && remote.FakenvapiVersion !== local.FakenvapiVersion) {
      fvUpdate.classList.remove('hidden');
    } else {
      fvUpdate.classList.add('hidden');
    }

    // NukemFG
    document.getElementById('about-nk-ver').textContent = local.NukemFGVersion || t.notProvided;
    const nkUpdate = document.getElementById('about-nk-update');
    if (remote.NukemFGVersion && remote.NukemFGVersion !== local.NukemFGVersion) {
      nkUpdate.classList.remove('hidden');
    } else {
      nkUpdate.classList.add('hidden');
    }

    // Build date
    document.getElementById('about-build-date').textContent = new Date().toISOString().split('T')[0];
  } catch (e) {
    console.error('Failed to load about info:', e);
  }
}

document
  .getElementById("btn-check-update")
  .addEventListener("click", async () => {
    if (!App) return;
    const t = translations[currentLang];
    const btn = document.getElementById("btn-check-update");
    btn.disabled = true;
    btn.innerHTML = `<span class="spinner"></span> ${t.checkingUpdatesPre}`;
    try {
      const info = await App.CheckAppUpdate();
      if (info.available) {
        showToast(t.updateAvailToast.replace("{0}", info.version), "success");
      } else {
        showToast(t.latestVersionToast, "success");
      }
    } catch (e) {
      showToast(t.updateCheckFailed, "error");
    } finally {
      btn.disabled = false;
      btn.textContent = t.checkUpdatesBtn;
    }
  });

document.getElementById("btn-github").addEventListener("click", async () => {
  if (!App) return;
  try {
    const config = await App.GetConfig();
    const url = `https://github.com/${config.App.RepoOwner}/${config.App.RepoName}`;
    window.runtime?.BrowserOpenURL(url);
  } catch (e) {
    window.runtime?.BrowserOpenURL("https://github.com");
  }
});

// ─── GPU Badge ──────────────────────────────────────────────────────────
async function loadGPU() {
  if (!App) return;
  const t = translations[currentLang];
  try {
    const result = await App.GetGPU();
    let display = result.display || t.noGpu;

    // Detect emoji and replace with styled span for better rendering
    if (display.includes("🟢")) {
      display = display.replace("🟢", '<span class="gpu-dot" style="--dot-color: #22c55e"></span>');
    } else if (display.includes("🔴")) {
      display = display.replace("🔴", '<span class="gpu-dot" style="--dot-color: #ef4444"></span>');
    } else if (display.includes("🔵")) {
      display = display.replace("🔵", '<span class="gpu-dot" style="--dot-color: #3b82f6"></span>');
    } else if (display.includes("⚪")) {
      display = display.replace("⚪", '<span class="gpu-dot" style="--dot-color: #94a3b8"></span>');
    }

    const gpuText = document.getElementById("gpu-text");
    gpuText.innerHTML = display;
    gpuText.style.display = "flex";
    gpuText.style.alignItems = "center";
    gpuText.style.gap = "8px";
  } catch (e) {
    document.getElementById("gpu-text").textContent = t.gpuError;
  }
}

// ─── Native Swapper / Upscaler Library ──────────────────────────────────
let swapperManifest = null;
let currentLibTab = "dlss";

async function loadLibrary() {
  if (!App) return;
  try {
    if (!swapperManifest) {
      swapperManifest = await App.FetchSwapperManifest();
    }
    renderLibraryTab(currentLibTab);
  } catch (e) {
    console.error("Failed to load swapper manifest:", e);
    showToast(translations[currentLang].failedManifest, "error");
  }
}

document.getElementById("btn-refresh-manifest").addEventListener("click", async () => {
    swapperManifest = null;
    await loadLibrary();
    showToast(translations[currentLang].manifestRefreshed, "success");
});

document.querySelectorAll(".lib-tab").forEach(tab => {
    tab.addEventListener("click", (e) => {
        document.querySelectorAll(".lib-tab").forEach(t => t.classList.remove("active"));
        tab.classList.add("active");
        currentLibTab = tab.dataset.target;
        renderLibraryTab(currentLibTab);
    });
});

function getManifestRecordsForType(type) {
  if (!swapperManifest) return [];
  switch (type) {
    case "dlss": return swapperManifest.DLSS || [];
    case "dlssg": return swapperManifest.DLSS_G || [];
    case "dlssd": return swapperManifest.DLSS_D || [];
    case "fsr31dx12": return swapperManifest.FSR_31_DX12 || [];
    case "fsr31vk": return swapperManifest.FSR_31_VK || [];
    case "xess": return swapperManifest.XeSS || [];
    case "xessfg": return swapperManifest.XeSS_FG || [];
    default: return [];
  }
}

function renderLibraryTab(type) {
  const tbody = document.getElementById("library-table-body");
  const records = getManifestRecordsForType(type);
  const t = translations[currentLang];
  
  if (records.length === 0) {
    tbody.innerHTML = `<tr><td colspan='4'>${t.noVersManifest}</td></tr>`;
    return;
  }

  // Sort descending by version number string
  const sorted = [...records].sort((a,b) => b.version.localeCompare(a.version, undefined, {numeric: true}));

  tbody.innerHTML = sorted.map((rec, idx) => {
    const sizeMb = (rec.zip_file_size / 1024 / 1024).toFixed(2);
    let statusHTML = rec.downloaded 
      ? `<span class="status-downloaded">${t.libDownloaded}</span>`
      : `<span class="status-cloud">${t.libCloud}</span>`;
    
    let actionHTML = rec.downloaded
      ? `<button class="btn btn-secondary btn-sm" disabled>${t.libReady}</button>`
      : `<button class="btn btn-primary btn-sm btn-down-dll" data-type="${type}" data-hash="${rec.md5_hash}">${t.libDownloadBtn}</button>`;

    return `
      <tr>
        <td><strong>v${rec.version}</strong> ${rec.is_dev_file ? '<span class="update-badge">Dev</span>' : ''}</td>
        <td>${sizeMb} MB</td>
        <td>${statusHTML}</td>
        <td>${actionHTML}</td>
      </tr>
    `;
  }).join("");

  tbody.querySelectorAll(".btn-down-dll").forEach(btn => {
     btn.addEventListener("click", async (e) => {
        const hash = btn.dataset.hash;
        const targetType = btn.dataset.type;
        const record = getManifestRecordsForType(targetType).find(r => r.md5_hash === hash);
        if (!record) return;
        
        btn.disabled = true;
        btn.innerHTML = `<span class="spinner" style="width:12px;height:12px;border-width:2px;"></span>`;
        try {
           await App.DownloadOfficialDLL(record);
           showToast(t.downloadedSucc.replace("{0}", record.version), "success");
           // Reload to update status cleanly
           swapperManifest = await App.FetchSwapperManifest();
           renderLibraryTab(currentLibTab);
           if (currentView === "manage") populateNativeVersionsForEachRow();
        } catch (err) {
           showToast(`${t.failedDownload}: ${err}`, "error");
           btn.disabled = false;
           btn.innerHTML = t.libDownloadBtn;
        }
     });
  });
}

function populateNativeVersionsForEachRow() {
   const rows = document.querySelectorAll(".swapper-tech-row");
   if (!rows.length) return;
   const t = translations[currentLang];

   rows.forEach(row => {
      const type = row.dataset.tech;
      const currentVer = row.dataset.currentVer;
      const listContainer = row.querySelector(".native-version-list");
      const restoreBtn = row.querySelector(".btn-restore-native-row");

      // Show ALL versions from the manifest as a scrollable list
      const allRecords = getManifestRecordsForType(type);
      if (allRecords.length === 0) {
         listContainer.innerHTML = `<div style="padding:12px; text-align:center; color:var(--text-muted); font-size:11px;">${t.noVersManifest}</div>`;
      } else {
         const sorted = [...allRecords].sort((a,b) => b.version.localeCompare(a.version, undefined, {numeric: true}));
         listContainer.innerHTML = sorted.map(r => {
            const sizeMb = r.zip_file_size ? (r.zip_file_size / 1024 / 1024).toFixed(1) : '?';
            let actionBtn;
            if (r.downloaded) {
               actionBtn = `<button class="btn btn-primary btn-sm btn-install-ver" data-hash="${r.md5_hash}" style="padding:2px 10px; font-size:10px;">${t.btnInstallVer}</button>`;
            } else {
               actionBtn = `<button class="btn btn-secondary btn-sm btn-install-ver" data-hash="${r.md5_hash}" style="padding:2px 10px; font-size:10px;">${t.btnDownInstallVer}</button>`;
            }
            return `<div class="native-ver-item" style="display:flex; align-items:center; justify-content:space-between; padding:6px 10px; border-bottom:1px solid var(--border-subtle);">
               <div>
                  <strong style="font-size:12px;">v${r.version}</strong>
                  <span style="font-size:10px; color:var(--text-muted); margin-left:6px;">${sizeMb} MB</span>
                  ${r.downloaded ? '<span style="font-size:9px; color:var(--accent); margin-left:4px;">✓</span>' : ''}
               </div>
               ${actionBtn}
            </div>`;
         }).join("");

         // Bind install buttons
         listContainer.querySelectorAll(".btn-install-ver").forEach(btn => {
            btn.addEventListener("click", async () => {
               if (!App || managingIndex < 0) return;
               const game = games[managingIndex];
               const hash = btn.dataset.hash;
               const record = getManifestRecordsForType(type).find(r => r.md5_hash === hash);
               if (!record) return;

               btn.disabled = true;
               btn.innerHTML = `<span class="spinner" style="width:10px;height:10px;border-width:2px;"></span>`;
               try {
                  if (!record.downloaded) {
                     await App.DownloadOfficialDLL(record);
                     showToast(t.downloadedSucc.replace("{0}", record.version), "success");
                     swapperManifest = await App.FetchSwapperManifest();
                  }
                  await App.SwapOfficialDLL(record, game.installPath);
                  showToast(t.installedSuccess, "success");
                  // Re-scan to detect new DLL versions
                  games = (await App.ScanGames()) || [];
                  openManage(managingIndex);
               } catch(err) {
                  showToast(`Failed: ${err}`, "error");
                  btn.disabled = false;
                  btn.textContent = record.downloaded ? t.btnInstallVer : t.btnDownInstallVer;
               }
            });
         });
      }

      // Bind Restore Event
      restoreBtn.replaceWith(restoreBtn.cloneNode(true));
      const newRestoreBtn = row.querySelector(".btn-restore-native-row");
      newRestoreBtn.addEventListener("click", async () => {
         if (!App || managingIndex < 0) return;
         const game = games[managingIndex];
         try {
           await App.RestoreOfficialDLL(game.installPath, type);
           showToast(t.restoreSuccess, "success");
           games = (await App.ScanGames()) || [];
           openManage(managingIndex);
         } catch(err) {
           showToast(`Restore failed: ${err}`, "error");
         }
      });
   });
}

// NVAPI Preset Change Event
document.querySelectorAll(".select-preset").forEach(el => {
   el.addEventListener("change", async (e) => {
      if (!App || managingIndex < 0) return;
      const game = games[managingIndex];
      const tech = e.target.dataset.tech;
      const presetVal = parseInt(e.target.value);
      try {
         await App.SetNVAPIPreset(game.executablePath, tech, presetVal);
         showToast(`DLSS Preset updated successfully!`, "success");
      } catch (err) {
         showToast(`Failed to update preset: ${err}`, "error");
      }
   });
});

// ─── Toast ──────────────────────────────────────────────────────────────
function showToast(message, type = "success") {
  const toast = document.getElementById("toast");
  toast.textContent = message;
  toast.className = `toast ${type}`;
  setTimeout(() => toast.classList.add("hidden"), 3000);
}

// ─── Init ───────────────────────────────────────────────────────────────
document.addEventListener("DOMContentLoaded", () => {
  loadGames();
  loadGPU();
});
