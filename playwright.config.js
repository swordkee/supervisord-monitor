import { defineConfig } from '@playwright/test';

export default defineConfig({
  use: {
    browserName: 'chromium',
    launchOptions: {
      executablePath: '/home/swordkee/.cache/ms-playwright/chromium-1200/chrome-linux64/chrome',
    },
  },
});
