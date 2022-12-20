const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    baseUrl: "http://localhost:3000",
    supportFile: false,
    videoUploadOnPasses: false,
    setupNodeEvents(on, config) {},
  },
});
