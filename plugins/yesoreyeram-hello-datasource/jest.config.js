const standard = require('@grafana/toolkit/src/config/jest.plugin.config');

module.exports = {
  ...{
    ...standard.jestConfig(),
    setupFilesAfterEnv: ['<rootDir>/src/tests.ts'],
    modulePathIgnorePatterns: [],
  },
};
