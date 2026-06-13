// Inline SVG language marks for the header, status bar and home chip.
// Hand-authored simplified-but-recognizable brand marks, tuned for 16–20px
// legibility, each in its conventional brand hue. Sized via the .langico
// class (~1em, vertical-align tuned to sit on the text baseline).

const ICO = `class="langico" viewBox="0 0 24 24" aria-hidden="true" focusable="false"`;

const ICONS: Record<string, string> = {
  // The two-snake interlocked mark — blue over yellow, eyes cut out (evenodd).
  python: `<svg ${ICO}>
    <path fill="#3776AB" fill-rule="evenodd" d="M11.9 1.5c-3 0-4.6 1-4.6 3v2.2h4.8v.9H5.3c-2.2 0-3.3 1.9-3.3 4.4 0 2.5 1.1 4.4 3.3 4.4h2v-2.8c0-1.8 1.5-3.3 3.3-3.3h4.4c1.5 0 2.7-1.2 2.7-2.7v-3c0-2.1-2.7-3.1-5.8-3.1zm-2.2 1.6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>
    <path fill="#FFD43B" fill-rule="evenodd" d="M12.1 22.5c3 0 4.6-1 4.6-3v-2.2h-4.8v-.9h6.8c2.2 0 3.3-1.9 3.3-4.4 0-2.5-1.1-4.4-3.3-4.4h-2v2.8c0 1.8-1.5 3.3-3.3 3.3H9c-1.5 0-2.7 1.2-2.7 2.7v3c0 2.1 2.7 3.1 5.8 3.1zm2.2-1.6a1 1 0 1 1 0-2 1 1 0 0 1 0 2z"/>
  </svg>`,
  // The gopher doesn't survive 16px — the "GO" wordmark in Go cyan does.
  go: `<svg ${ICO}>
    <text x="12" y="17" text-anchor="middle" font-family="'Arial Black',Arial,sans-serif" font-weight="800" font-size="12.5" letter-spacing="-0.6" fill="#00ADD8">GO</text>
  </svg>`,
  // Gear ring (toothed stroke) with the serif R hint.
  rust: `<svg ${ICO}>
    <circle cx="12" cy="12" r="9.4" fill="none" stroke="#B7410E" stroke-width="2.6" stroke-dasharray="2.1 2.02"/>
    <circle cx="12" cy="12" r="7.8" fill="#B7410E"/>
    <text x="12" y="16.6" text-anchor="middle" font-family="Georgia,'Times New Roman',serif" font-weight="700" font-size="12.5" fill="#DEA584">R</text>
  </svg>`,
  typescript: `<svg ${ICO}>
    <rect x="1.5" y="1.5" width="21" height="21" rx="4.5" fill="#3178C6"/>
    <text x="12" y="16.4" text-anchor="middle" font-family="Arial,sans-serif" font-weight="700" font-size="10.5" fill="#FFFFFF">TS</text>
  </svg>`,
  javascript: `<svg ${ICO}>
    <rect x="1.5" y="1.5" width="21" height="21" rx="4.5" fill="#F7DF1E"/>
    <text x="12" y="16.4" text-anchor="middle" font-family="Arial,sans-serif" font-weight="700" font-size="10.5" fill="#111111">JS</text>
  </svg>`,
  // Coffee cup with steam — orange cup, red steam (the Java duo-tone).
  java: `<svg ${ICO}>
    <path fill="#E76F00" d="M5.4 10h11.2v5.3a4.2 4.2 0 0 1-4.2 4.2H9.6a4.2 4.2 0 0 1-4.2-4.2z"/>
    <path fill="none" stroke="#E76F00" stroke-width="1.8" d="M16.6 11.3h1.6a2.4 2.4 0 0 1 0 4.8h-1.6"/>
    <path fill="none" stroke="#EA2D2E" stroke-width="1.7" stroke-linecap="round" d="M9.1 7.8C7.7 6.4 10.5 5.4 9.1 4"/>
    <path fill="none" stroke="#EA2D2E" stroke-width="1.7" stroke-linecap="round" d="M12.8 8.1c-1.2-1.2 1.1-2 0-3.3"/>
    <path fill="none" stroke="#E76F00" stroke-width="1.7" stroke-linecap="round" d="M4.6 21.6h12.8"/>
  </svg>`,
  csharp: `<svg ${ICO}>
    <polygon points="12,1 2.5,6.5 2.5,17.5 12,23 21.5,17.5 21.5,6.5" fill="#68217A"/>
    <text x="12" y="15.9" text-anchor="middle" font-family="Arial,sans-serif" font-weight="700" font-size="10" fill="#FFFFFF">C#</text>
  </svg>`,
  cpp: `<svg ${ICO}>
    <polygon points="12,1 2.5,6.5 2.5,17.5 12,23 21.5,17.5 21.5,6.5" fill="#00599C"/>
    <text x="12" y="15.5" text-anchor="middle" font-family="Arial,sans-serif" font-weight="700" font-size="9" letter-spacing="-0.5" fill="#FFFFFF">C++</text>
  </svg>`,
};

// langIcon returns the inline SVG mark for lang (a neutral ◇ for unknown).
export function langIcon(lang: string): string {
  return ICONS[lang] ?? '<span class="langico langico-fb" aria-hidden="true">&#9671;</span>';
}
