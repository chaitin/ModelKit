/// <reference types="vite/client" />

import type { PaletteColorChannel } from '@mui/material';

declare module '@mui/material/styles' {
  interface TypeText {
    tertiary: string;
    auxiliary: string;
  }

  interface TypeBackground {
    paper2: string;
  }

  interface Palette {
    light: Palette['primary'] & PaletteColorChannel;
    dark: Palette['primary'] & PaletteColorChannel;
    disabled: Palette['primary'] & PaletteColorChannel;
    risk: {
      severe: string;
      critical: string;
      suggest: string;
    };
  }

  // allow configuration using `createTheme`
  interface PaletteOptions {
    light?: PaletteOptions['primary'] & Partial<PaletteColorChannel>;
    dark?: PaletteOptions['primary'] & Partial<PaletteColorChannel>;
    disabled?: PaletteOptions['primary'] & Partial<PaletteColorChannel>;
    risk?: {
      severe?: string;
      critical?: string;
      suggest?: string;
    };
    text?: Partial<TypeText>;
    background?: Partial<TypeBackground>;
  }
}

declare module '@mui/material/Button' {
  interface ButtonPropsColorOverrides {
    light: true;
    dark: true;
  }
}

import type {} from '@mui/material/themeCssVarsAugmentation';